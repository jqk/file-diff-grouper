package filediff

import (
	"encoding/base64"
	"fmt"
	"hash/crc64"
	"os"
	"path/filepath"
	"strconv"

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

	provider := createChecksumProvider()
	result := &ScanResult{
		HeaderSize:          config.HeaderSize,
		Dir:                 config.Dir,
		Filter:              config.Filter,
		FileCount:           0,
		FileSize:            0,
		Method:              provider.Method(),
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
		identity, e := getFileIdentity(filename, config.HeaderSize, buffer, config.NeedFullChecksum, provider)
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

		key := checksumToString(identity.HeaderChecksum)
		if _, ok := result.Files[key]; !ok {
			// 结果集中不存在对应 HeaderChecksum 的文件数组，说明是新的 HeaderChecksum，
			// 要添加一个新的文件数组(实际不添加也会自动添加)，同时计数。
			result.Files[key] = []*FileIdentity{}
			result.HeaderChecksumCount++
		}
		result.Files[key] = append(result.Files[key], identity)

		return nil
	})

	if err != nil {
		return nil, err
	}

	sortAndFindDupFiles(result)

	// 保存文件的耗时不计入工作耗时。
	stopwatch.Stop()
	result.ElapsedTime = stopwatch.ElapsedTime()

	if err := SaveScanResult(result, config.ScanResultFile); err != nil {
		return nil, err
	}

	return result, nil
}

func sortAndFindDupFiles(r *ScanResult) {
	for _, identities := range r.Files {
		if len(identities) > 1 { // 多个文件具有相同的 headerChecksum 才需排序并查重。
			sortFileIdentities(identities) // 1. 排序。

			m := FileIdentities{} // 2. 准备查重。

			for _, id := range identities {
				// 使用文件长度加整体校验和作为 key。这是判断文件是否重复的标准。即使 fullChecksum 无值，也要加上。
				key := strconv.FormatInt(id.FileSize, 10) + "_" + checksumToString(id.FullChecksum)
				m[key] = append(m[key], id)
			}

			for _, ids := range m {
				count := len(ids) - 1 // 如果 3 个文件完全相同，说明有 2 个是重复的，所以减一。
				if count > 0 {        // 某个 key 出现了多于 1 次，说明有重复。
					r.DupGroupCount++
					r.DupFileCount += count
					r.DupFileSize += ids[0].FileSize * int64(count)
				}
			}
		}
	}
}

func checksumToString(checksum []byte) string {
	return base64.StdEncoding.EncodeToString(checksum)
}

/*
getFileIdentity 通过计算文件的整体校验和来获取文件的标识信息。
*/
func getFileIdentity(
	filename string,
	headerSize int,
	buffer []byte,
	needFullChecksum bool,
	provider *fileutils.CommonFileChecksumProvider,
) (*FileIdentity, error) {

	if err := fileutils.GetFileChecksumWithProvider(
		filename, headerSize, buffer, true, needFullChecksum, provider); err != nil {
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

func createChecksumProvider() *fileutils.CommonFileChecksumProvider {
	//return fileutils.NewCommonFileChecksumProvider("CRC32-IEEE", crc32.NewIEEE())
	return fileutils.NewCommonFileChecksumProvider("CRC64-ISO", crc64.New(crc64.MakeTable(crc64.ISO)))
	//return fileutils.NewCommonFileChecksumProvider("MD5", md5.New())
}
