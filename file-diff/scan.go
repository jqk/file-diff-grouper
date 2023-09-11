package filediff

import (
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"

	"github.com/jqk/futool4go/fileutils"
	"github.com/jqk/futool4go/timeutils"
)

/*
ScanTargetDir scans the directory specified in the config parameter and save the result.

Parameters:
  - config: the configuration object. not nil.

Returns:
  - [ScanResult] object.
  - an error if any.
*/
func ScanBaseDir(config *Config) (*ScanResult, error) {
	dirConfig, _ := getDirConfig(config)
	return ScanDir(dirConfig)
}

/*
ScanTargetDir scans the directory specified in the config parameter and save the result.

Parameters:
  - config: the configuration object. not nil.

Returns:
  - [ScanResult] object.
  - an error if any.
*/
func ScanTargetDir(config *Config) (*ScanResult, error) {
	_, dirConfig := getDirConfig(config)
	return ScanDir(dirConfig)
}

/*
ScanDir scans the directory specified in the config parameter and save the result.

Parameters:
  - config: the configuration object. not nil.

Returns:
  - [ScanResult] object.
  - an error if any.
*/
func ScanDir(config *DirConfig) (*ScanResult, error) {
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
func getFileIdentity(
	filename string,
	headerSize int,
	buffer []byte,
	needFullChecksum bool,
) (*FileIdentity, error) {

	// 使用 crc32 作为校验和算法。
	crc := crc32.NewIEEE()
	// 仅对必要字段初始化，其它字段后面会更新。
	identity := &FileIdentity{
		HasFullChecksum: needFullChecksum,
		Filename:        filename,
	}

	// 开始定义 fileutils.GetFileChecksum() 需要的 3 个回调函数。
	// 1. 函数定义参见 fileutils.GetFileChecksum() 的第 4 个参数。
	calculator := func(data []byte) (int, error) {
		return crc.Write(data) // 只是计算校验和。
	}

	// 2. 函数定义参见 fileutils.GetFileChecksum() 的第 5 个参数。
	headerReadyHandler := func(info os.FileInfo, fullIsReady bool) error {
		identity.HeaderChecksum = crc.Sum32() // 保存文件头的校验和。
		identity.FileSize = info.Size()
		identity.ModifiedTime = info.ModTime()

		if fullIsReady {
			// fullIsReady 为 true 表示文件长度小于等于 headerSize，所以整体校验和就是文件头的校验和。
			identity.HasFullChecksum = true
			identity.FullChecksum = identity.HeaderChecksum
		}
		return nil
	}

	// 3. 函数定义参见 fileutils.GetFileChecksum() 的第 6 个参数。
	fullReadyHandler := func(info os.FileInfo) error {
		// 此时 identity.HasFullChecksum 必然为 true
		identity.FullChecksum = crc.Sum32() // 保存整体校验和。
		return nil
	}

	// 是否使用最后一个回调函数由 needFullChecksum 参数决定。
	// 下面将使用 fullHander 变量而不是直接使用 fullReadyHandler 调用 GetFileChecksum()。
	fullHandler := fullReadyHandler
	if !needFullChecksum {
		// 通过给定的方法指针是否为 nil 来判断是否需要计算完整校验和。
		fullHandler = nil
	}

	// 前面都是准备过程，此处调用 GetFileChecksum() 进行实际读取及计算。
	if err := fileutils.GetFileChecksum(filename, headerSize, buffer,
		calculator, headerReadyHandler, fullHandler); err != nil {
		return nil, err
	}

	return identity, nil
}
