package main

import (
	"errors"
	filediff "file-diff"
	"os"
	"strings"
	"time"
)

func main() {
	showVersion()

	argCount := len(os.Args)

	if argCount == 1 {
		showHelp()
	} else if argCount == 2 {
		config, err := filediff.LoadConfigFromFile(os.Args[1])
		if err != nil {
			showError("Load config error", err, false)
			os.Exit(2)
		}

		var groupResult *filediff.GroupResult = nil
		var scanBaseResult *filediff.ScanResult = nil
		var scanTargetResult *filediff.ScanResult = nil
		count := 0
		done := make(chan struct{}) // 用于协程同步的通道。

		go func() { // 启动单独的协程执行任务。
			if strings.EqualFold(config.Action, filediff.ActionCompare) {
				groupResult, err = filediff.GroupFileDiff(config, func(fileIdentity *filediff.FileIdentity) {
					count++
				})
				if err != nil {
					showError("Group files error", err, false)
				}
			} else if strings.EqualFold(config.Action, filediff.ActionScanBase) {
				scanBaseResult, err = filediff.ScanBaseDir(config, func(fileIdentity *filediff.FileIdentity) {
					count++
				})
				if err != nil {
					showError("ScanBaseDir error", err, false)
				}
			} else if strings.EqualFold(config.Action, filediff.ActionScanTarget) {
				scanTargetResult, err = filediff.ScanTargetDir(config, func(fileIdentity *filediff.FileIdentity) {
					count++
				})
				if err != nil {
					showError("ScanTargetDir error", err, false)
				}
			}

			close(done)
		}()

		// 等待扩展名扫描结束。并显示扩展名扫描进度。
		sleepTime := 500 * time.Millisecond
		for {
			time.Sleep(sleepTime)

			select {
			case <-done: // 等待扫描结束。
				if err != nil {
					os.Exit(2)
				}

				if groupResult != nil {
					showGroupResult(groupResult, config)
				} else if scanBaseResult != nil {
					showScanResult(scanBaseResult, config, config.CompareBase.ScanResultFile)
				} else if scanTargetResult != nil {
					showScanResult(scanTargetResult, config, config.CompareTarget.ScanResultFile)
				}
				return
			default: // 打印扫描进度。
				showScanProgress(count)
			}
		}
	} else {
		showError("Argument error", errors.New("wrong number of argument"), true)
		os.Exit(1)
	}
}
