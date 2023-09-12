package main

import (
	filediff "file-diff"
	"fmt"
	"path/filepath"

	"github.com/jqk/futool4go/common"
)

func showVersion() {
	fmt.Println()
	fmt.Println("Copyright (c) 1999-2023 Not a dream Co., Ltd.")
	fmt.Println("file difference grouper (fdg) 0.9.1, 2023-09-13")
	fmt.Println()
}

func showHelp() {
	fmt.Println("Usage:")
	fmt.Println("  fdg [path/to/configFile]")
	fmt.Println("      Compare and group the file differences according to specified config file.")
	fmt.Println()
	fmt.Println("Otherwise: show this help.")
	fmt.Println("See <https://github.com/jqk/file-diff-grouper> for more information.")
	fmt.Println()
}

func showError(header string, err error, includingHelp bool) {
	fmt.Printf("%s: %s\n", header, err)

	if includingHelp {
		showHelp()
	}
}

func showGroupResult(result *filediff.GroupResult, config *filediff.Config) {
	if progressShowed {
		fmt.Println()
	}

	baseDir, _ := filepath.Abs(result.More.BaseDir)
	targetDir, _ := filepath.Abs(result.More.TargetDir)

	fmt.Printf("Action: %s\n\n", config.Action)

	fmt.Printf("Header size        : %d\n", config.HeaderSize)
	fmt.Printf("CompareFullChecksum: %t\n", config.CompareTarget.CompareFullChecksum)
	fmt.Printf("Base   Dir         : %s\n", baseDir)
	fmt.Printf("Target Dir         : %s\n", targetDir)
	fmt.Printf("Base file count    : %d\n", result.More.BaseFileCount)
	fmt.Printf("Target file count  : %d\n", result.More.TargetFileCount)
	fmt.Printf("More file count    : %d\n", len(result.More.FileGroup.Files))
	fmt.Printf("More file size     : %s\n", common.ToSizeString(result.More.FileGroup.Size))
	fmt.Printf("Same file count    : %d\n", len(result.Same.FileGroup.Files))
	fmt.Printf("Same file size     : %s\n", common.ToSizeString(result.Same.FileGroup.Size))
	fmt.Printf("Time elapsed       : %s\n", result.ElapsedTime)

	if config.CompareTarget.MoveMore {
		fmt.Printf("Move MORE files to : %s\n", result.MoreDir)
	} else {
		fmt.Println("Moving MORE files  : not required.")
	}

	if config.CompareTarget.MoveSame {
		fmt.Printf("Move SAME files to : %s\n", result.SameDir)
	} else {
		fmt.Println("Moving SAME files  : not required.")
	}

	fmt.Println()
}

func showScanResult(result *filediff.ScanResult, config *filediff.Config, resultFilename string) {
	if progressShowed {
		fmt.Println()
	}
	dir, _ := filepath.Abs(result.Dir)

	fmt.Printf("Action: %s\n\n", config.Action)
	fmt.Printf("Header size         : %d\n", config.HeaderSize)
	fmt.Printf("Scan dir            : %s\n", dir)
	fmt.Printf("Scan file count     : %d\n", result.FileCount)
	fmt.Printf("Scan file size      : %s\n", common.ToSizeString(result.FileSize))
	fmt.Printf("HeaderChecksum Count: %d\n", result.HeaderChecksumCount)

	resultFilename, _ = filepath.Abs(resultFilename)

	fmt.Printf("Scan result file    : %s\n", resultFilename)
	fmt.Printf("Time elapsed        : %s\n\n", result.ElapsedTime)
}

var progressShowed = false

func showScanProgress(count int) {
	progressShowed = true
	fmt.Printf("File scaned: %d\n", count)
}
