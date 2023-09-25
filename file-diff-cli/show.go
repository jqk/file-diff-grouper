package main

import (
	filediff "file-diff"
	"fmt"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/jqk/futool4go/common"
)

var blue color.Style = color.New(color.LightBlue)
var green color.Style = color.New(color.LightGreen)
var white color.Style = color.New(color.White)
var yellow color.Style = color.New(color.LightYellow)
var red color.Style = color.New(color.LightRed)

func showVersion() {
	white.Println("\nCopyright (c) 1999-2023 Not a dream Co., Ltd.")
	white.Print("file difference grouper (")
	blue.Print("fdg")
	white.Print(") ")
	red.Print("1.1.1")
	white.Print(", 2023-09-25\n\n")
}

func showHelp() {
	yellow.Println("Usage:")
	yellow.Println("  fdg [path/to/configFile]")
	white.Println("      Compare and group the file differences according to specified config file.")
	white.Println()
	white.Println("Otherwise: show this help.")
	white.Print("See <")
	yellow.Print("https://github.com/jqk/file-diff-grouper")
	white.Println("> for more information.")
	fmt.Println()
}

func showError(header string, err error, includingHelp bool) {
	white.Printf("%s: ", header)
	color.Errorf("%s", err)
	fmt.Println()

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
	backupDir, _ := filepath.Abs(config.CompareTarget.BackupDir)
	moreFileCount := len(result.More.FileGroup.Files)
	sameFileCount := len(result.Same.FileGroup.Files)

	green.Print("Action: ")
	yellow.Printf("%s\n\n", config.Action)

	green.Print("Base   Dir         : ")
	yellow.Printf("%s\n", baseDir)
	green.Print("Target Dir         : ")
	yellow.Printf("%s\n", targetDir)
	green.Print("Header size        : ")
	yellow.Printf("%d\n", config.HeaderSize)
	green.Print("Base file count    : ")
	yellow.Printf("%d\n", result.More.BaseFileCount)
	green.Print("Target file count  : ")
	yellow.Printf("%d\n", result.More.TargetFileCount)
	green.Print("More file count    : ")
	yellow.Printf("%d\n", moreFileCount)
	green.Print("More file size     : ")
	yellow.Printf("%s\n", common.ToSizeString(result.More.FileGroup.Size))
	green.Print("Same file count    : ")
	yellow.Printf("%d\n", sameFileCount)
	green.Print("Same file size     : ")
	yellow.Printf("%s\n", common.ToSizeString(result.Same.FileGroup.Size))
	green.Print("Time elapsed       : ")
	yellow.Printf("%s\n", result.ElapsedTime)
	green.Print("Backup Dir         : ")
	yellow.Printf("%s\n", backupDir)
	if config.CompareTarget.MoveMore {
		if moreFileCount == 0 {
			green.Println("Moving MORE files  : nothing to move.")
		} else {
			green.Print("Move MORE files to : ")
			yellow.Printf("%s\n", result.MoreDir)
		}
	} else {
		green.Println("Moving MORE files  : not required.")
	}

	if config.CompareTarget.MoveSame {
		if sameFileCount == 0 {
			green.Println("Moving SAME files  : nothing to move.")
		} else {
			green.Print("Move SAME files to : ")
			yellow.Printf("%s\n", result.SameDir)
		}
	} else {
		green.Println("Moving SAME files  : not required.")
	}

	fmt.Println()
}

func showScanResult(result *filediff.ScanResult, config *filediff.Config, resultFilename string) {
	if progressShowed {
		green.Println()
	}

	green.Print("Action ")
	yellow.Printf("%s\n\n", config.Action)

	dir, _ := filepath.Abs(result.Dir)
	resultFilename, _ = filepath.Abs(resultFilename)

	green.Print("Method               : ")
	yellow.Printf("%s\n", result.Method)
	green.Print("Header size          : ")
	yellow.Printf("%d\n", result.HeaderSize)
	green.Print("Need full checksum   : ")
	yellow.Printf("%t\n", result.NeedFullChecksum)
	green.Print("Compare full checksum: ")
	yellow.Printf("%t\n", result.CompareFullChecksum)
	green.Print("Scan dir             : ")
	yellow.Printf("%s\n", dir)
	green.Print("Scan result file     : ")
	yellow.Printf("%s\n", resultFilename)
	green.Print("Scan file count      : ")
	yellow.Printf("%d\n", result.FileCount)
	green.Print("Scan file size       : ")
	yellow.Printf("%s\n", common.ToSizeString(result.FileSize))
	green.Print("HeaderChecksum count : ")
	yellow.Printf("%d\n", result.HeaderChecksumCount)
	green.Print("FullChecksum count   : ")
	yellow.Printf("%d\n", result.FullChecksumCount)
	green.Print("Duplicate group count: ")
	yellow.Printf("%d\n", result.DupGroupCount)
	green.Print("Duplicate file count : ")
	yellow.Printf("%d\n", result.DupFileCount)
	green.Print("Duplicate file size  : ")
	yellow.Printf("%s\n", common.ToSizeString(result.DupFileSize))
	green.Print("Time elapsed         : ")
	yellow.Printf("%s\n", result.ElapsedTime)

	fmt.Println()
}

var progressShowed = false

func showScanProgress(count int) {
	progressShowed = true
	white.Print("File scaned: ")
	green.Printf("%d\n", count)
}
