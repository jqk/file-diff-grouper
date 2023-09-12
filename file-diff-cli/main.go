package main

import (
	"errors"
	filediff "file-diff"
	"os"
	"strings"
)

func main() {
	showVersion()

	argCount := len(os.Args)

	if argCount == 1 {
		showHelp()
	} else if argCount == 2 {
		runConfigedTask(os.Args[1])
	} else {
		showError("Argument error", errors.New("wrong number of argument"), true)
		os.Exit(1)
	}
}

func runConfigedTask(filename string) {
	config, err := filediff.LoadConfigFromFile(filename)
	if err != nil {
		showError("Load config error", err, false)
		os.Exit(2)
	}

	if strings.EqualFold(config.Action, filediff.ActionCompare) {
		groupResult, err := filediff.GroupFileDiff(config)
		if err != nil {
			showError("Group files error", err, false)
			os.Exit(2)
		}

		showGroupResult(groupResult, config)
	} else if strings.EqualFold(config.Action, filediff.ActionScanBase) {
		scanBaseResult, err := filediff.ScanBaseDir(config)
		if err != nil {
			showError("ScanBaseDir error", err, false)
			os.Exit(2)
		}

		showScanResult(scanBaseResult, config, config.CompareBase.ScanResultFile)
	} else if strings.EqualFold(config.Action, filediff.ActionScanTarget) {
		scanTargetResult, err := filediff.ScanTargetDir(config)
		if err != nil {
			showError("ScanTargetDir error", err, false)
			os.Exit(2)
		}

		showScanResult(scanTargetResult, config, config.CompareTarget.ScanResultFile)
	}
}
