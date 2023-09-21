package filediff

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/jqk/futool4go/fileutils"
)

const (
	CompareResultMore = "more" // CompareResultMore indicates the files are in Target Dir but in Base Dir.
	CompareResultSame = "same" // CompareResultSame indicates the files taht exist in both Target Dir and Base Dir.
)

/*
FileGroup is a group of file names and their total sizes.
*/
type FileGroup struct {
	Size  int64    // total size of all files.
	Files []string // the list of file names.
}

/*
CompareRsult represents the compare result of files in target dir but not in base dir.
*/
type CompareResult struct {
	Method                    string            // Algorithm name
	BaseDir                   string            // Base directory path
	BaseFileCount             int               // Number of files in base directory
	BaseHeaderChecksumCount   int               // Number of header checksums calculated for base files
	TargetDir                 string            // Target directory path
	TargetFileCount           int               // Number of files in target directory
	TargetHeaderChecksumCount int               // Number of header checksums for target files
	Filter                    *fileutils.Filter // Filter used during comparison
	CompareFullChecksum       bool              // CompareFullChecksum indicates if full checksums should be compared when headerChecksum and file size are equal.
	CompareResultType         string            // Type of comparison result - "More/Same"
	ResultFileCount           int               // Number of result files, just for saving the result, same as len(FileGroup.Files)
	ResultFileSize            int64             // Total size of result files, just for saving the result, same as FileGroup.Size
	ElapsedTime               time.Duration     // Time taken for comparison
	FileGroup                 *FileGroup        `json:"-"` // Result file group
}

func createCompareResult(
	config *Config,
	baseScanResult *ScanResult,
	targetScanResult *ScanResult,
	elapsedTime time.Duration,
	compareResultType string,
	fileGroup *FileGroup,
) *CompareResult {
	return &CompareResult{
		Method:                    baseScanResult.Method, // baseScanResult.Method 和 targetScanResult.Method 是相同的。
		BaseDir:                   config.CompareBase.Dir,
		BaseFileCount:             baseScanResult.FileCount,
		BaseHeaderChecksumCount:   baseScanResult.HeaderChecksumCount,
		TargetDir:                 config.CompareTarget.Dir,
		TargetFileCount:           targetScanResult.FileCount,
		TargetHeaderChecksumCount: targetScanResult.HeaderChecksumCount,
		Filter:                    config.Filter,
		CompareFullChecksum:       config.CompareTarget.CompareFullChecksum,
		CompareResultType:         compareResultType,
		ElapsedTime:               elapsedTime,
		FileGroup:                 fileGroup,
		// 之所以有以下这两个冗余字段，是为了在保存文件时，
		// 能比较简单输出文件数量和文件大小两个信息。
		ResultFileCount: len(fileGroup.Files),
		ResultFileSize:  fileGroup.Size,
	}
}

func saveComareResult(result *CompareResult, filename string) error {
	// 确保目录可用。
	p := filepath.Dir(filename)
	if err := os.MkdirAll(p, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	// 空字符串为首行不缩进，其他行缩进 4 个空格
	encoder.SetIndent("", "    ")

	// 先将比较结果汇总信息以 JSON 格式写入文件。
	if err = encoder.Encode(result); err != nil {
		return err
	}

	if _, err = writer.WriteString("\n----------\n\n"); err != nil {
		return err
	}

	// 逐一写入比较所得的文件名。
	for _, fn := range result.FileGroup.Files {
		if _, err = writer.WriteString("\"" + fn + "\"\n"); err != nil {
			return err
		}
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	return nil
}
