package filediff

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	t "github.com/jqk/futool4go/timeutils"
)

/*
GroupResult represents the grouping result correspond to the compare result.
*/
type GroupResult struct {
	More        *CompareResult // files in target dir but not in base dir.
	Same        *CompareResult // files in both target dir and base dir.
	MoreDir     string         // backup dir the contains the More files.
	SameDir     string         // backup dir the contains the Same files.
	ElapsedTime time.Duration  // time taken for comparison.
}

/*
GroupFileDiff compares two directories and moves the compare result files into the backup dir.

Parameters:
  - config: Pointer to the Config struct containing the configuration options for the comparison.

Returns:
  - [GroupResult] object.
  - If an error occurs during the comparison process.
*/
func GroupFileDiff(config *Config) (*GroupResult, error) {
	result := &GroupResult{}
	var err error
	stopwatch := &t.Stopwatch{}
	stopwatch.Start()

	// 先比较两个目录。获得相同的，和多出的文件列表。
	result.More, result.Same, err = CompareDirs(config)
	if err != nil {
		return nil, err
	}

	// 每次都将结果文件移动到由当前时间标记的目录中。
	timeFlag := time.Now().Format("20060102_150405")
	backupRoot, err := filepath.Abs(config.CompareTarget.BackupDir)
	if err != nil {
		return nil, err
	}

	// 移动目标路径，得到：backupRoot/YYYYMMDD_HHMMSS/'more or same'
	result.MoreDir = filepath.Join(backupRoot, timeFlag, CompareResultMore)
	result.SameDir = filepath.Join(backupRoot, timeFlag, CompareResultSame)

	sourceRoot, err := filepath.Abs(config.CompareTarget.Dir)
	if err != nil {
		return nil, err
	}

	if config.CompareTarget.MoveMore {
		if err = moveCompareResult(result.More.FileGroup, sourceRoot, result.MoreDir); err != nil {
			return nil, err
		}
	}
	if config.CompareTarget.MoveSame {
		if err = moveCompareResult(result.Same.FileGroup, sourceRoot, result.SameDir); err != nil {
			return nil, err
		}
	}

	stopwatch.Stop()
	result.ElapsedTime = stopwatch.ElapsedTime()

	return result, nil
}

/*
moveCompareResult moves all files in the CompareResult from sourceRoot to
targetRoot. It creates the necessary directories if they do not exist. It
returns an error if any operation fails.

Parameters:
  - result: A pointer to a CompareResult object that contains paths of all files that need to be moved.
  - sourceRoot: The root directory of the source files.
  - targetRoot: The root directory of the target files.

Returns:
  - error: Returns an error if any operation fails.
*/
func moveCompareResult(result *FileGroup, sourceRoot string, targetRoot string) error {
	for _, source := range result.Files {
		// 获得目标路径。
		target := strings.Replace(source, sourceRoot, targetRoot, 1)
		if err := os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
			return err
		} else if err = os.Rename(source, target); err != nil {
			// 执行移动文件的操作。为避免大量读写 IO，不提供复制操作。
			// 所以，要保证 compareTarget.Dir 是可写的。
			return err
		}
	}

	return nil
}
