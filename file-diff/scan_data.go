package filediff

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jqk/futool4go/fileutils"
)

// FileIdentity defines the identity attributes of a file.
type FileIdentity struct {
	HeaderChecksum  []byte    // HeaderChecksum is the checksum of the file header.
	HasFullChecksum bool      // HasFullChecksum indicates if there is a full checksum.
	FullChecksum    []byte    // FullChecksum is the full checksum of the file.
	Filename        string    // Filename is the name of the file.
	FileSize        int64     // FileSize is the size of the file.
	ModifiedTime    time.Time // ModifiedTime is the last modified time of the file.
}

// FileIdentities is a map of file identity.
//
// The key is the header checksum.
//
// The value is a slice of file identity with the same header checksum.
type FileIdentities map[string][]*FileIdentity

// ScanResult is the result of a scanning a given directory.
type ScanResult struct {
	Method              string            // Algorithm name
	HeaderSize          int               // HeaderSize is the size of the header read from each file.
	Dir                 string            // Dir is the directory path that was scanned.
	Filter              *fileutils.Filter // Filter is the filter used to select files to scan.
	FullChecksumChanged bool              `json:"-"` // FullChecksumChanged indicates if any full checksums changed.
	FileCount           int               // FileCount is the number of files scanned.
	FileSize            int64             // FileSize is the total size of all files scanned.
	HeaderChecksumCount int               // HeaderChecksumCount is the number of header checksums calculated.
	FullChecksumCount   int               // FullChecksumCount is the number of full checksums calculated.
	DupGroupCount       int               // DupGroupCount is the number of groups of files with duplicate header checksums and file sizes.
	DupFileCount        int               // DupFileCount is the number of files with duplicate header checksums and file sizes.
	DupFileSize         int64             // DupFileSize is the total size of files with duplicate header checksums and file sizes.
	ElapsedTime         time.Duration     // ElapsedTime is the duration of the scan.
	Files               FileIdentities    // Files contains the detailed scan results for each file.
}

// Diff compares two [FileIdentity] objects.
//
// See [Differ] for more information.
func (f *FileIdentity) Diff(other Differ) string {
	o := other.(*FileIdentity)
	if f == o {
		return ""
	}

	if !bytes.Equal(f.HeaderChecksum, o.HeaderChecksum) {
		return "FileIdentity.HeaderChecksum"
	}
	if f.HasFullChecksum != o.HasFullChecksum {
		return "FileIdentity.HasFullChecksum"
	}
	if !bytes.Equal(f.FullChecksum, o.FullChecksum) {
		return "FileIdentity.FullChecksum"
	}
	if f.Filename != o.Filename {
		return "FileIdentity.Filename"
	}
	if f.FileSize != o.FileSize {
		return "FileIdentity.FileSize"
	}
	if f.ModifiedTime != o.ModifiedTime {
		return "FileIdentity.ModifiedTime"
	}

	return ""
}

// Diff compares two [ScanResult] objects.
//
// See [Differ] for more information.
func (r *ScanResult) Diff(other Differ) string {
	o := other.(*ScanResult)

	if r.HeaderSize != o.HeaderSize {
		return "ScanResult.HeaderSize"
	}
	if s := r.Filter.Diff(o.Filter); s != "" {
		return "ScanResult." + s
	}
	if r.FullChecksumChanged != o.FullChecksumChanged {
		return "ScanResult.FullChecksumChanged"
	}
	if r.FileCount != o.FileCount {
		return "ScanResult.FileCount"
	}
	if r.FileSize != o.FileSize {
		return "ScanResult.FileSize"
	}
	if r.HeaderChecksumCount != o.HeaderChecksumCount {
		return "ScanResult.HeaderChecksumCount"
	}
	if r.FullChecksumCount != o.FullChecksumCount {
		return "ScanResult.FullChecksumCount"
	}
	if r.ElapsedTime != o.ElapsedTime {
		return "ScanResult.Elapsed"
	}
	if len(r.Files) != len(o.Files) {
		return "ScanResult.FilesLength"
	}

	for k, v := range r.Files {
		// k 是文件头校验和。
		// v 是具有相同文件头校验和的文件信息数组。
		if len(v) != len(o.Files[k]) {
			// 这里包含了 k 在 *o 中不存在的情况。
			return "ScanResult.Files.HeaderChecksum: " + k
		}

		for i := range v {
			if v[i].Diff(o.Files[k][i]) != "" {
				return "ScanResult.Files.Identity: " + v[i].Filename
			}
		}
	}

	return ""
}

// sortFileIdentities sorts the given slice of FileIdentity in ascending order by FileSize.
//
// identities: a slice of pointers to FileIdentity structs to be sorted.
//
// No return value.
func sortFileIdentities(identities []*FileIdentity) {
	sort.Slice(identities, func(i, j int) bool {
		if identities[i].FileSize < identities[j].FileSize {
			return true
		} else if identities[i].FileSize == identities[j].FileSize {
			return identities[i].Filename < identities[j].Filename
		}

		return false
	})
}

/*
SaveScanResult saves the scan result data to a file with the given filename.

Parameters:
  - result: The [ScanResult] containing the ChecksumInfo to save
  - filename: The path of the file to save the data to. empty means do not save any data.

Returns:
  - An error if the file cannot be created or the data cannot be encoded
*/
func SaveScanResult(result *ScanResult, filename string) error {
	// 文件名为空表示不保存任何数据。
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return nil
	}

	// 确保目录存在。
	p := filepath.Dir(filename)
	if err := os.MkdirAll(p, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	// 空字符串为首行不缩进，其他行缩进 4 个空格
	encoder.SetIndent("", "    ")

	if err = encoder.Encode(result); err != nil {
		return err
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	return nil
}

/*
LoadScanResult loads the scan result from the last scan and returns a [ScanResult]
containing the result data or an error if the file cannot be opened or decoded.

Parameters:
  - filename: The path of the file containing the last scan result

Returns:
  - A [ScanResult] containing the loaded scan result data
  - An error if the file cannot be opened or decoded
*/
func LoadScanResult(filename string) (*ScanResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var result ScanResult

	for decoder.More() {
		if err = decoder.Decode(&result); err != nil {
			return nil, err
		}
	}

	return &result, nil
}
