package filediff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanDir(t *testing.T) {
	config, err := LoadConfigFromFile("test-data/scan.test.config.yaml")
	assert.NotNil(t, config)
	assert.Nil(t, err)

	base, target := getDirConfig(config)
	result, err := ScanDir(base)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, 5, result.FileCount)
	// 有 4 个文件长度小于 headerSize，所以直接算是有了 fullChecksum。
	assert.Equal(t, 4, result.FullChecksumCount)
	// 有 3 个文件是完全一样的，所以只有 3 个 headerChecksum。
	assert.Equal(t, 3, len(result.Files))
	// 以下几个值和校验和直接相对，修改文件后校验和会发生变化，所以要特别注意。
	assert.Equal(t, 1, len(result.Files[3096586316]))
	assert.Equal(t, 1, len(result.Files[3222652411]))
	assert.Equal(t, 3, len(result.Files[4245835769]))

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
	result, err = ScanDir(target)
	assert.Nil(t, err)
	assert.NotNil(t, result)
}
