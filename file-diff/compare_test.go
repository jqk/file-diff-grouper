package filediff

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareDirs(t *testing.T) {
	config, err := LoadConfigFromFile("test-data/compare.test.config.yaml")
	assert.NotNil(t, config)
	assert.Nil(t, err)

	more, same, err := CompareDirs(config, nil)
	assert.Nil(t, err)
	assert.NotNil(t, more)
	assert.NotNil(t, same)

	// 在 compare_origin 扫描到 6 个文件：
	// - size_too_large.txt
	// - 002.txt
	// - 003-same-as-002.md
	// - copy-of-003.md
	// - 004.txt
	// - 001.md
	// in_exclude_filter.log 和 not_in_include_filter.c 被过滤掉，未选取。
	//
	// ------------------------------------------
	// 在 compare_target 扫描到 3 个文件：
	// - 011.txt
	// - 013.md
	// - 010-same-as-001.md
	// 012.log 被过滤掉了。
	//
	// ------------------------------------------
	// 有 1 个重复的，有 2 个 origin 中没有的。
	// 以下信息可以查看 compare_test_base.scan.json 和 compare_test_target.scan.json，
	// 以及 target-more-than-base.txt 和 target-same-with-base.txt。

	assert.Equal(t, 6, same.BaseFileCount)
	assert.Equal(t, 4, same.BaseHeaderChecksumCount)
	assert.Equal(t, 3, same.TargetFileCount)
	assert.Equal(t, 3, same.TargetHeaderChecksumCount)
	assert.Equal(t, "same", same.CompareResultType)
	assert.Equal(t, 1, same.ResultFileCount)
	assert.Equal(t, 1, len(same.FileGroup.Files))
	assert.True(t, strings.Index(same.FileGroup.Files[0], "010-same-as-001.md") > 0)

	assert.Equal(t, "more", more.CompareResultType)
	assert.Equal(t, 2, len(more.FileGroup.Files))
	assert.Equal(t, int64(35879), more.FileGroup.Size)
}
