package filediff

import (
	"errors"
	"strings"

	"github.com/jqk/futool4go/fileutils"
)

// Config holds the configuration for a file comparison operation.
type Config struct {
	HeaderSize    int               `mapstructure:"headerSize"`    // HeaderSize is the size in bytes to read from the header of each file.
	BufferSize    int               `mapstructure:"bufferSize"`    // BufferSize is the size in bytes of the buffer used when reading files.
	Action        string            `mapstructure:"action"`        // Action is the action to perform - "Compare", "ScanBase", etc.
	CompareBase   *CompareBase      `mapstructure:"compareBase"`   // CompareBase is the base directory information for comparison.
	CompareTarget *CompareTarget    `mapstructure:"compareTarget"` // CompareTarget is the target directory information for comparison.
	Filter        *fileutils.Filter `mapstructure:"filter"`        // Filter defines filters to apply for file selection.
}

// CompareBase defines the information for the base directory in a comparison.
type CompareBase struct {
	Dir              string `mapstructure:"dir"`              // Dir is the base directory path to compare.
	ScanResultFile   string `mapstructure:"scanResultFile"`   // ScanResultFile is the file path for a the scan result. Empty string indicates no need to save the scan result.
	LoadScanResult   bool   `mapstructure:"loadScanResult"`   // LoadScanResult indicates if the scan result should be loaded. Defaults to false.
	NeedFullChecksum bool   `mapstructure:"needFullChecksum"` // NeedFullChecksum indicates if full checksums are needed. Defaults to false.
}

// CompareTarget defines the information for the target directory in a comparison.
type CompareTarget struct {
	Dir              string `mapstructure:"dir"`              // Dir is the base directory path to compare.
	ScanResultFile   string `mapstructure:"scanResultFile"`   // ScanResultFile is the file path for a the scan result. Empty string indicates no need to save the scan result.
	LoadScanResult   bool   `mapstructure:"loadscanResult"`   // LoadScanResult indicates if the scan result should be loaded. Defaults to false.
	NeedFullChecksum bool   `mapstructure:"needFullChecksum"` // NeedFullChecksum indicates if full checksums are needed. Defaults to false.
	BackupDir        string `mapstructure:"backupDir"`        // BackupDir is the backup directory path for moving reulst files.
	MoveMore         bool   `mapstructure:"moveMore"`         // MoveMore indicates if more files should be moved. Defaults to false. 'more' refers to files that are present in the target dir but not in the base dir.
	MoveSame         bool   `mapstructure:"moveSame"`         // MoveSame indicates if same files should be moved. Defaults to false. 'same' refers to files that are present in both the target dir and base dir.
}

// DirConfig holds the configuration for a directory scan.
type DirConfig struct {
	HeaderSize       int               // HeaderSize is the size in bytes to read from the header of each file.
	BufferSize       int               // BufferSize is the size in bytes of the buffer used when reading files.
	Dir              string            // Dir is the path to be scanned.
	ScanResultFile   string            // ScanResultFile is the file path for a the scan result. Empty string indicates no need to save the scan result.
	LoadScanResult   bool              // LoadScanResult indicates if the scan result should be loaded. Defaults to false.
	NeedFullChecksum bool              // NeedFullChecksum indicates if full checksums are needed. Defaults to false.
	Filter           *fileutils.Filter // Filter defines filters to apply for file selection.
}

const (
	ActionCompare    = "Compare"    // Compare Base directory and Target directory.
	ActionScanBase   = "ScanBase"   // Only scan Base directory.
	ActionScanTarget = "ScanTarget" // Only scan Target directory.
)

type Validator interface {
	/**
	Validate is responsible for performing some kind of validation and returning an error if the validation fails.

	Returns:
		An error object if the validation fails, or None if the validation succeeds.
	*/
	Validate() error
}

type Differ interface {
	/**
	Diff returns a string representing the first difference found between the
	receiver and the other Differ instance. If both instances are equal, an empty
	string is returned.

	Parameters:
		- other (Differ): The other instance to compare against.

	Returns:
		A string indicating the first difference found, or an empty string if
		no differences were found.
	*/
	Diff(other Differ) string
}

// getDirConfig parses the Config information and returns the directory information for CompareBase and CompareTarget.
func getDirConfig(Config *Config) (*DirConfig, *DirConfig) {
	return &DirConfig{
			HeaderSize:       Config.HeaderSize,
			BufferSize:       Config.BufferSize,
			Dir:              Config.CompareBase.Dir,
			ScanResultFile:   Config.CompareBase.ScanResultFile,
			LoadScanResult:   Config.CompareBase.LoadScanResult,
			NeedFullChecksum: Config.CompareBase.NeedFullChecksum,
			Filter:           Config.Filter,
		}, &DirConfig{
			HeaderSize:       Config.HeaderSize,
			BufferSize:       Config.BufferSize,
			Dir:              Config.CompareTarget.Dir,
			ScanResultFile:   Config.CompareTarget.ScanResultFile,
			LoadScanResult:   Config.CompareTarget.LoadScanResult,
			NeedFullChecksum: Config.CompareTarget.NeedFullChecksum,
			Filter:           Config.Filter,
		}
}

// Diff compares two [CompareBase] objects.
//
// See [Differ] for more information.
func (c *CompareBase) Diff(other Differ) string {
	o := other.(*CompareBase)
	if c == o {
		return ""
	}

	if c.Dir != o.Dir {
		return "CompareBase.Dir"
	}
	if c.ScanResultFile != o.ScanResultFile {
		return "CompareBase.ScanResultFile"
	}
	if c.LoadScanResult != o.LoadScanResult {
		return "CompareBase.LoadScanResult"
	}
	if c.NeedFullChecksum != o.NeedFullChecksum {
		return "CompareBase.NeedFullChecksum"
	}

	return ""
}

// Validate validates [CompareBase] object.
//
// See [Validator] for more information.
func (c *CompareBase) Validate() error {
	if c.Dir = strings.TrimSpace(c.Dir); c.Dir == "" {
		return errors.New("CompareBase.Dir must not be empty")
	}

	c.ScanResultFile = strings.TrimSpace(c.ScanResultFile)
	c.ScanResultFile = strings.Replace(c.ScanResultFile, "${dir}", c.Dir, 1)

	if c.LoadScanResult && c.ScanResultFile == "" {
		return errors.New("CompareBase.ScanResultFile must not be empty")
	}

	return nil
}

// Diff compares two [CompareTarget] objects.
//
// See [Differ] for more information.
func (c *CompareTarget) Diff(other Differ) string {
	o := other.(*CompareTarget)
	if c == o {
		return ""
	}

	if c.Dir != o.Dir {
		return "CompareTarget.Dir"
	}
	if c.ScanResultFile != o.ScanResultFile {
		return "CompareTarget.ScanResultFile"
	}
	if c.LoadScanResult != o.LoadScanResult {
		return "CompareTarget.LoadScanResult"
	}
	if c.NeedFullChecksum != o.NeedFullChecksum {
		return "CompareTarget.NeedFullChecksum"
	}
	if c.BackupDir != o.BackupDir {
		return "CompareTarget.BackupDir"
	}
	if c.MoveMore != o.MoveMore {
		return "CompareTarget.MoveMore"
	}
	if c.MoveSame != o.MoveSame {
		return "CompareTarget.MoveSame"
	}

	return ""
}

// Validate validates [CompareTarget] object.
//
// See [Validator] for more information.
func (c *CompareTarget) Validate() error {
	if c.Dir = strings.TrimSpace(c.Dir); c.Dir == "" {
		return errors.New("CompareTarget.Dir must not be empty")
	}

	c.ScanResultFile = strings.TrimSpace(c.ScanResultFile)
	c.ScanResultFile = strings.Replace(c.ScanResultFile, "${dir}", c.Dir, 1)

	if c.LoadScanResult && c.ScanResultFile == "" {
		return errors.New("CompareTarget.ScanResultFile must not be empty")
	}

	c.BackupDir = strings.TrimSpace(c.BackupDir)
	c.BackupDir = strings.Replace(c.BackupDir, "${dir}", c.Dir, 1)
	if c.BackupDir == "" {
		return errors.New("CompareTarget.BackupDir must not be empty")
	}

	return nil
}

// Diff compares two [Config] objects.
//
// See [Differ] for more information.
func (c *Config) Diff(other Differ) string {
	o := other.(*Config)
	if c == o {
		return ""
	}

	if c.HeaderSize != o.HeaderSize {
		return "Config.HeaderSize"
	} else if c.BufferSize != o.BufferSize {
		return "Config.BufferSize"
	} else if !strings.EqualFold(c.Action, o.Action) {
		return "Config.Action"
	} else if s := c.CompareBase.Diff(o.CompareBase); s != "" {
		return s
	} else if s = c.CompareTarget.Diff(o.CompareTarget); s != "" {
		return s
	} else if s = c.Filter.Diff(o.Filter); s != "" {
		return s
	}

	return ""
}

// Validate validates [Config] object.
//
// See [Validator] for more information.
func (c *Config) Validate() error {
	if strings.EqualFold(c.Action, ActionCompare) {
		c.Action = ActionCompare
	} else if strings.EqualFold(c.Action, ActionScanBase) {
		c.Action = ActionScanBase
	} else if strings.EqualFold(c.Action, ActionScanTarget) {
		c.Action = ActionScanTarget
	} else {
		return errors.New("Config.Action must be one of [compare, scanbase, scantarget]")
	}

	if c.CompareBase == nil {
		return errors.New("Config.CompareBase must not be nil")
	}
	if c.CompareTarget == nil {
		return errors.New("Config.CompareTarget must not be nil")
	}
	if c.Filter == nil {
		return errors.New("Config.Filter must not be nil")
	}

	if e := c.CompareBase.Validate(); e != nil {
		return e
	} else if e = c.CompareTarget.Validate(); e != nil {
		return e
	} else if e = c.Filter.Validate(); e != nil {
		return e
	}

	if c.HeaderSize < 1024 {
		c.HeaderSize = 1024
	}
	if c.BufferSize < c.HeaderSize {
		c.BufferSize = c.HeaderSize
	}

	return nil
}
