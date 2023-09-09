package filediff

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jqk/futool4go/fileutils"

	"github.com/stretchr/testify/assert"
)

func TestGroupFileDiff(t *testing.T) {
	config, err := LoadConfigFromFile("test-data/group.test.config.yaml")
	assert.NotNil(t, config)
	assert.Nil(t, err)

	// 删除 output 中的 target 目录。
	err = os.RemoveAll(config.CompareTarget.Dir)
	assert.Nil(t, err)

	// 从 origin 中复制整个 target 目录到 output 中。
	baseParent := filepath.Dir(config.CompareBase.Dir)
	targetSource := filepath.Join(baseParent, "compare_target")
	err = fileutils.CopyDir(targetSource, config.CompareTarget.Dir)
	assert.Nil(t, err)

	// 移动文件。
	result, err := GroupFileDiff(config)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	// 具体结果与 compare_test.go 中的一致。请参见其说明。
	assert.Equal(t, 1, len(result.Same.FileGroup.Files))
	assert.Equal(t, 2, len(result.More.FileGroup.Files))

	exist, isDir, err := fileutils.FileExists(result.SameDir)
	assert.Nil(t, err)
	assert.True(t, exist)
	assert.True(t, isDir)
	// 实在懒得再查检每个该移动的文件在不在了。
}
