package main

import (
	"errors"
	filediff "file-diff"
	"os"
	"strings"
	"time"
)

const (
	ArgumentError = 1
	RuntimeError  = 2
)

func main() {
	showVersion()
	argCount := len(os.Args)

	if argCount == 1 {
		showHelp()
	} else if argCount != 2 {
		showError("Argument error", errors.New("wrong number of argument"), true)
		os.Exit(ArgumentError)
	} else {
		config, err := filediff.LoadConfigFromFile(os.Args[1])
		if err != nil {
			showError("Load config error", err, false)
			os.Exit(RuntimeError)
		}

		var groupResult *filediff.GroupResult = nil
		var scanBaseResult, scanTargetResult *filediff.ScanResult = nil, nil
		lastCount, count := 0, 0
		done := make(chan struct{}) // 用于协程同步的通道。
		fileScanedHandler := func(*filediff.FileIdentity) {
			count++
		}

		go func() { // 启动协程执行任务，主线程负责输出工作状态及结果。
			if strings.EqualFold(config.Action, filediff.ActionCompare) {
				groupResult, err = filediff.GroupFileDiff(config, fileScanedHandler)
				if err != nil {
					showError("Group files error", err, false)
				}
			} else if strings.EqualFold(config.Action, filediff.ActionScanBase) {
				scanBaseResult, err = filediff.ScanBaseDir(config, fileScanedHandler)
				if err != nil {
					showError("ScanBaseDir error", err, false)
				}
			} else if strings.EqualFold(config.Action, filediff.ActionScanTarget) {
				scanTargetResult, err = filediff.ScanTargetDir(config, fileScanedHandler)
				if err != nil {
					showError("ScanTargetDir error", err, false)
				}
			}

			close(done) // 关闭通道。通知主线程结束等待。
		}()

		// 等待扩展名扫描结束。并显示扩展名扫描进度。
		sleepTime := 800 * time.Millisecond
		for {
			select {
			case <-done: // 等待扫描协程工作结束。
				if err != nil { // 错误已经在协程结束前打印过了。
					os.Exit(RuntimeError)
				} else if groupResult != nil {
					showGroupResult(groupResult, config)
				} else if scanBaseResult != nil {
					showScanResult(scanBaseResult, config, config.CompareBase.ScanResultFile)
				} else if scanTargetResult != nil {
					showScanResult(scanTargetResult, config, config.CompareTarget.ScanResultFile)
				}
				return
			default: // 打印扫描进度。
				// 避免扫描文件时由于 IO 等待或文件很大，造成处理时间超过循环等待时间，进而重复显示同一数量。
				if lastCount != count {
					lastCount = count
					showScanProgress(count)
				}
				time.Sleep(sleepTime)
			}
		}
	}
	// 默认结束并退出，相当于 os.Exit(0)。
}
