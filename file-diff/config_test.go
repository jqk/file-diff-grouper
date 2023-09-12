package filediff

import (
	"testing"

	"github.com/jqk/futool4go/fileutils"
)

// Notice that tabs are not allowed in the content string.
// The test will fail if tabs are used.
// However, when inserting new rows, tabs sometimes occur accidentally.
var content string = `
HeaderSize: 1024
BufferSize: 10240
action: "compare"

compareBase:
  dir: "test-data/origin/compare_base_0"
  scanResultFile: "${dir}/~~~data/base-scan.data.json"
  loadScanResult: true
  needFullChecksum: true

compareTarget:
  dir: "test-data/origin/compare_target"
  scanResultFile: "${dir}/~~~data/target-scan.data.json"
  loadScanResult: true
  needFullChecksum: true

  compareFullChecksum: true
  backupDir: "${dir}/~~~result"
  moveMore: false
  moveSame: false

filter:
  caseSensitive: true
  include:
  - "*.txt"
  - "*.md"

  exclude:
  - "*.logg"

  minFileSize: 1024
  maxFileSize: 1048576`

var expected Config = Config{
	HeaderSize: 1024,
	BufferSize: 10240,
	Action:     "compare",
	CompareBase: &CompareBase{
		Dir:              "test-data/origin/compare_base_0",
		ScanResultFile:   "${dir}/~~~data/base-scan.data.json",
		LoadScanResult:   true,
		NeedFullChecksum: true,
	},
	CompareTarget: &CompareTarget{
		Dir:                 "test-data/origin/compare_target",
		ScanResultFile:      "${dir}/~~~data/target-scan.data.json",
		LoadScanResult:      true,
		NeedFullChecksum:    true,
		CompareFullChecksum: true,
		BackupDir:           "${dir}/~~~result",
		MoveMore:            false,
		MoveSame:            false,
	},
	Filter: &fileutils.Filter{
		CaseSensitive: true,
		Include: []string{
			"*.md",
			"*.txt",
		},
		Exclude: []string{
			"*.logg",
		},
		MinFileSize: 1024,
		MaxFileSize: 1048576,
	},
}

func TestLoadConfigFromString(t *testing.T) {
	config, err := LoadConfigFromString(content, "yaml")
	if err != nil {
		t.Errorf("LoadConfigFromString() returned error: %v", err)
	}

	expected.Validate()
	if s := config.Diff(&expected); s != "" {
		t.Errorf("Diff found: %s\nLoadConfigFromString() returned:\n%+v\nExpected:\n%+v", s, *config, expected)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	config, err := LoadConfigFromFile("test-data/config.yaml")
	if err != nil {
		t.Errorf("LoadConfigFromFile() returned error: %v", err)
	}

	expected.Validate()
	if s := config.Diff(&expected); s != "" {
		t.Errorf("Diff found: %s\nLoadConfigFromFile() returned:\n%+v\nExpected:\n%+v", s, *config, expected)
	}
}
