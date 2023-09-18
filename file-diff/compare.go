package filediff

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"

	"github.com/jqk/futool4go/fileutils"
	t "github.com/jqk/futool4go/timeutils"
)

/*
CompareDirs compares the scan results of two directories specified by the Config struct.
Before comparing, it will scan the directories if needed.

Parameters:
  - config: Pointer to the Config struct containing the configuration options for the comparison.

Returns:
  - CompareResult pointers containing the result of files in target dir but not in base dir.
  - CompareResult pointers containing the result of files in both target dir and base dir.
  - If an error occurs during the comparison process.
*/
func CompareDirs(config *Config, handler FileScanedFunc) (resultMore *CompareResult, resultSame *CompareResult, err error) {
	baseConfig, targetConfig := getDirConfig(config)
	var baseScanResult, targetScanResult *ScanResult
	var more, same *FileGroup

	// 获取两个目录的扫描结果。根据配置定义，有可能是装载以前的扫描结果，也有可能是新的扫描结果。
	if baseScanResult, err = getScanResult(baseConfig, handler); err != nil {
		return nil, nil, err
	}
	if targetScanResult, err = getScanResult(targetConfig, handler); err != nil {
		return nil, nil, err
	}

	// 比较两个目录的扫描结果。将得到两个目录都有的文件列表 same，以及 target 比 base 多出的文件列表 more。
	elapsedTime, _ := (&t.Stopwatch{}).Elapsing(func() error {
		more, same, err = compareScanResults(
			baseScanResult,
			targetScanResult,
			config.HeaderSize,
			config.BufferSize,
			config.CompareTarget.CompareFullChecksum,
		)

		return err
	})

	if err != nil {
		return nil, nil, err
	}

	// 根据得到的两个文件列表，创建比较结果。
	resultMore = createCompareResult(config, baseScanResult, targetScanResult, elapsedTime, CompareResultMore, more)
	resultSame = createCompareResult(config, baseScanResult, targetScanResult, elapsedTime, CompareResultSame, same)

	// 保存 target 比 base 多的文件列表。
	filename := filepath.Join(config.CompareTarget.BackupDir, "target-more-than-base.txt")
	if err = saveComareResult(resultMore, filename); err != nil {
		return nil, nil, err
	}

	// 保存 target 和 base 相同的文件列表。
	filename = filepath.Join(config.CompareTarget.BackupDir, "target-same-with-base.txt")
	if err = saveComareResult(resultSame, filename); err != nil {
		return nil, nil, err
	}

	// 如果在比较期间更新了某些文件信息的 FullChecksum，必须重新保存扫描结果。
	if baseScanResult.FullChecksumChanged {
		// 保存文件的路径如果是空字符串，将不会执行任何操作。
		if err = SaveScanResult(baseScanResult, baseConfig.ScanResultFile); err != nil {
			return nil, nil, err
		}
	}
	if targetScanResult.FullChecksumChanged {
		if err = SaveScanResult(targetScanResult, targetConfig.ScanResultFile); err != nil {
			return nil, nil, err
		}
	}

	return resultMore, resultSame, nil
}

// getScanResult 获取扫描结果。根据配置定义，有可能是装载以前的扫描结果，也有可能是新的扫描结果。
func getScanResult(config *DirConfig, handler FileScanedFunc) (scanResult *ScanResult, err error) {
	if config.LoadScanResult {
		// 按指定读取以前的扫描结果文件。
		scanResult, err = LoadScanResult(config.ScanResultFile)
		if err != nil && err != os.ErrNotExist {
			// 文件不存在，则 scanResult 为 nil，可以执行扫描。
			// 否则，返回错误。
			var pathErr *os.PathError
			if !errors.As(err, &pathErr) {
				return nil, err
			}
		}

		if scanResult != nil && scanResult.HeaderSize != config.HeaderSize {
			// 文件头长度不同时比较没有意义。
			return nil, fmt.Errorf("HeaderSize mismatch: ScanResult(%d) != Config(%d)",
				scanResult.HeaderSize, config.HeaderSize)
		}
	}

	if scanResult == nil {
		// 本条件成立，说明要么 config.LoadScanResult 为 false，
		// 要么以前的扫描结果文件不存在，需重新执行扫描。
		if scanResult, err = ScanDir(config, handler); err != nil {
			return nil, err
		}
	}

	return scanResult, nil
}

/*
ensureFullChecksumReady ensures that a full checksum is ready for a given file identity.

Parameters:
  - r: a pointer to a ScanResult struct
  - f: a pointer to a FileIdentity struct
  - headerSize: an integer representing the size of the file header
  - bufferSize: the buffer
  - provider: a pointer to a CommonFileChecksumProvider

Returns:
  - an error if there was an issue getting the file identity.
*/
func ensureFullChecksumReady(
	r *ScanResult,
	f *FileIdentity,
	headerSize int,
	buffer []byte,
	provider *fileutils.CommonFileChecksumProvider[ChecksumType],
) error {
	if !f.HasFullChecksum {
		// 没有完整校验和，通过设置最后一个参数为 true，计算整体校验和。
		if c, err := getFileIdentity(f.Filename, headerSize, buffer, true, provider); err != nil {
			return err
		} else {
			f.HasFullChecksum = true
			f.FullChecksum = c.FullChecksum

			r.FullChecksumChanged = true
			r.FullChecksumCount++
		}
	}

	return nil
}

/*
compareScanResults compares two ScanResult structs. The headerSize and bufferSize are used
to calculate the full checksums of the files if needed.
*/
func compareScanResults(
	base,
	target *ScanResult,
	headerSize,
	bufferSize int,
	compareFullChecksum bool,
) (more, same *FileGroup, err error) {

	provider := createChecksumProvider()
	buffer := make([]byte, bufferSize)
	more = &FileGroup{}
	same = &FileGroup{}

	for targetHeaderChecksum, targetFiles := range target.Files {
		// targetHeaderChecksum 就是 target 目录中文件的 HeaderChecksum。
		if baseFiles, ok := base.Files[targetHeaderChecksum]; ok {
			// 在 base 中有和 targetHeaderChecksum 相同的文件数组，
			// 继续在 baseFiles 中查找与 targetFile 元素相同的元素。
			for _, targetFile := range targetFiles {
				foundSame := false

				for _, baseFile := range baseFiles {
					if baseFile.FileSize == targetFile.FileSize {
						if !compareFullChecksum { // HeaderChecksum 相同且文件长度相同，粗略认为两者相同。
							foundSame = true
							break
						}

						// compareFullChecksum 为 true，则要继续对比 FullChecksum。首先确保 FullChecksum 有效。
						if err = ensureFullChecksumReady(base, baseFile, headerSize, buffer, provider); err != nil {
							return // nil, nil, err
						}
						if err = ensureFullChecksumReady(target, targetFile, headerSize, buffer, provider); err != nil {
							return // nil, nil, err
						}
						if reflect.DeepEqual(baseFile.FullChecksum, targetFile.FullChecksum) {
							foundSame = true
							break
						}
					}
				}

				if foundSame {
					same.Files = append(same.Files, targetFile.Filename)
					same.Size += targetFile.FileSize
				} else {
					more.Files = append(more.Files, targetFile.Filename)
					more.Size += targetFile.FileSize
				}
			}
		} else {
			// headerChecksum 都对不上，肯定是多出的文件。
			for _, targetFile := range targetFiles {
				more.Files = append(more.Files, targetFile.Filename)
				more.Size += targetFile.FileSize
			}
		}
	}

	// 还是排一下序比较好。
	sort.Strings(more.Files)
	sort.Strings(same.Files)

	return // more, same, nil
}
