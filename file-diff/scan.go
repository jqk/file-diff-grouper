package filediff

import (
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"

	"github.com/jqk/futool4go/fileutils"
	"github.com/jqk/futool4go/timeutils"
)

type FileScanedFunc func(*FileIdentity)

/*
ScanTargetDir scans the directory specified in the config parameter and save the result.

Parameters:
  - config: the configuration object. not nil.

Returns:
  - [ScanResult] object.
  - an error if any.
*/
func ScanBaseDir(config *Config, handler FileScanedFunc) (*ScanResult, error) {
	dirConfig, _ := getDirConfig(config)
	return ScanDir(dirConfig, handler)
}

/*
ScanTargetDir scans the directory specified in the config parameter and save the result.

Parameters:
  - config: the configuration object. not nil.

Returns:
  - [ScanResult] object.
  - an error if any.
*/
func ScanTargetDir(config *Config, handler FileScanedFunc) (*ScanResult, error) {
	_, dirConfig := getDirConfig(config)
	return ScanDir(dirConfig, handler)
}

/*
ScanDir scans the directory specified in the config parameter and save the result.

Parameters:
  - config: the configuration object. not nil.

Returns:
  - [ScanResult] object.
  - an error if any.
*/
func ScanDir(config *DirConfig, handler FileScanedFunc) (*ScanResult, error) {
	result := &ScanResult{
		HeaderSize:          config.HeaderSize,
		Dir:                 config.Dir,
		Filter:              config.Filter,
		FileCount:           0,
		FileSize:            0,
		HeaderChecksumCount: 0,
		FullChecksumCount:   0,
		ElapsedTime:         0,
		Files:               FileIdentities{},
	}

	buffer := make([]byte, config.BufferSize) // 准备可以重复使用的缓冲区。
	var err error
	stopwatch := timeutils.Stopwatch{}
	stopwatch.Start()

	err = config.Filter.GetEachFile(config.Dir, nil, func(path string, info os.FileInfo) error {
		// 因为 GetEachFile() 调用本代码给出的 path 必然是有效的，所以 Abs() 不会返回错误，也就不必处理。
		filename, _ := filepath.Abs(path)
		identity, e := getFileIdentity(filename, config.HeaderSize, buffer, config.NeedFullChecksum)
		if e != nil {
			return fmt.Errorf("%s: %s", e, filename)
		}

		if handler != nil {
			handler(identity)
		}

		result.FileSize += identity.FileSize
		result.FileCount++
		if identity.HasFullChecksum {
			// 即使 config.NeedFullChecksum 为 false，但文件可能比较小，所以也可能有整体校验和。
			result.FullChecksumCount++
		}

		key := identity.HeaderChecksum
		if _, ok := result.Files[key]; !ok {
			// 结果集中不存在对应 HeaderChecksum 的文件数组，说明是新的 HeaderChecksum，
			// 要添加一个新的文件数组，同时计数。
			result.Files[key] = []*FileIdentity{}
			result.HeaderChecksumCount++
		}
		result.Files[key] = append(result.Files[key], identity)

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 保存文件的耗时不计入工作耗时。
	stopwatch.Stop()
	result.ElapsedTime = stopwatch.ElapsedTime()

	if err := SaveScanResult(result, config.ScanResultFile); err != nil {
		return nil, err
	}

	return result, nil
}

/*
getFileIdentity 通过计算文件的整体校验和来获取文件的标识信息。
*/
func getFileIdentity(filename string, headerSize int, buffer []byte, needFullChecksum bool) (*FileIdentity, error) {
	crc := crc32.NewIEEE() // 使用 crc32 作为校验和算法。仅需改变此处的 hash 及取得最后结果的函数引用即可。
	sumFunc := crc.Sum32   // 注意函数必须是前面建立的 hash 对象的函数引用。
	provider := fileutils.NewCommonFileChecksumProvider[uint32](crc, sumFunc)

	if err := fileutils.GetFileChecksumWithProvider[uint32](
		filename, headerSize, buffer, provider, true, needFullChecksum); err != nil {
		return nil, err
	}

	identity := &FileIdentity{
		Filename:        filename,
		HasFullChecksum: provider.IsFullChecksumReady(),
		HeaderChecksum:  provider.HeaderChecksum(),
		FullChecksum:    provider.FullChecksum(),
		FileSize:        provider.FileInfo().Size(),
		ModifiedTime:    provider.FileInfo().ModTime(),
	}

	return identity, nil
}
