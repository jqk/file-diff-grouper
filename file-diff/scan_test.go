package filediff

import (
	"encoding/base64"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanDir(t *testing.T) {
	config, err := LoadConfigFromFile("test-data/scan.test.config.yaml")
	assert.NotNil(t, config)
	assert.Nil(t, err)

	fileIdentities = fileIdentities[:0]
	base, target := getDirConfig(config)
	result, err := ScanDir(base, fileScanedHandler)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, len(fileIdentities), result.FileCount)
	assert.Equal(t, 5, result.FileCount)
	// 有 4 个文件长度小于 headerSize，所以直接算是有了 fullChecksum。
	assert.Equal(t, 4, result.FullChecksumCount)
	// 有 3 个文件是完全一样的，所以只有 3 个 headerChecksum。
	assert.Equal(t, 3, len(result.Files))
	// 以下几个值和校验和直接相对，修改文件后校验和会发生变化，所以要特别注意。
	// 以下 3 行是使用 CRC32-IEEE 时，直接使用 uint32 作为 key 和时拿到的数据。
	// assert.Equal(t, 1, len(result.Files[3096586316]))
	// assert.Equal(t, 1, len(result.Files[3222652411]))
	// assert.Equal(t, 3, len(result.Files[4245835769]))

	// 以下 3 行是使用 CRC64-ISO 时，直接使用 uint64 作为 key 和时拿到的数据。
	// assert.Equal(t, 1, len(result.Files[1570070392997697241]))
	// assert.Equal(t, 1, len(result.Files[18030433853017517113]))
	// assert.Equal(t, 3, len(result.Files[10364516874543699336]))

	// 使用 base64 的 uint64 为 key 时，与直接使用 uint64 时，结果是一样的。
	b, _ := base64.StdEncoding.DecodeString("+jj4D1tJbDk=")
	v := binary.BigEndian.Uint64(b)
	assert.Equal(t, v, uint64(18030433853017517113))
	assert.Equal(t, 1, len(result.Files["+jj4D1tJbDk="]))
	assert.Equal(t, 1, len(result.Files["FcoCtC78Vtk="]))
	assert.Equal(t, 3, len(result.Files["j9YpLw+4FYg="]))

	// 使用 MD5。
	// assert.Equal(t, 1, len(result.Files["0Z5pNX6LtUul4nLTdITfgw=="]))
	// assert.Equal(t, 1, len(result.Files["x1Usc48X8zTtWMdpWQ9lZw=="]))
	// assert.Equal(t, 3, len(result.Files["75J9Lgmiz5krPXNX+rlXyA=="]))

	// 由于 base.ScanResultFile 不是空字符串，所以在 ScanDir() 时执行了
	// SaveScanResult() 操作。
	data, err := LoadScanResult(base.ScanResultFile)
	assert.Nil(t, err)
	assert.NotNil(t, data)

	// 对比保存之前的，和从结果文件中读取的扫描结果。
	diff := result.Diff(data)
	assert.Equal(t, "", diff)

	// 也执行一下 target。但因为 target.ScanResultFile 是空字符串，所以
	// 不会保存结果文件。
	fileIdentities = fileIdentities[:0]
	result, err = ScanDir(target, fileScanedHandler)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(fileIdentities), result.FileCount)
}

var fileIdentities = make([]*FileIdentity, 10)

func fileScanedHandler(id *FileIdentity) {
	fileIdentities = append(fileIdentities, id)
}
