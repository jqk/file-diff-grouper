# File Diff Grouper

[中文](readme-CN.md)

## What is it

`FileDiffGrouper` is a command line tool for comparing file differences between two directories. It compares all files in the two directories byte-by-byte based on content, rather than filenames. Although only tested on Windows 10/11, it does not use any OS-specific features, so theoretically can run on Linux and MacOS as well.

`FileDiffGrouper` compares two directories - one called the Base directory, the other called the Target directory. `FileDiffGrouper` will output two result sets, saved to result files:

- MORE set: Files that exist in Target but not in Base, referred to as "extra files" in this document.
- SAME set: Files that exist in both Target and Base and have identical content, referred to as "duplicate files" in this document.

`FileDiffGrouper` can optionally move the files in the MORE and SAME sets to specified backup directories, based on options. For safety, it does not provide functionality to automatically delete these files.

## What to solve

There are many software tools that provide functionality to compare file differences between two directories and find duplicate files, such as:

- [alldup](https://alldup.info), free.
- [AntiDupl](https://github.com/ermig1979/AntiDupl), open source.
- [beyondcompare](https://www.scootersoftware.com)
- [czkawka](https://github.com/qarmin/czkawka), open source, mainly for finding duplicate pictures, also has a GUI version.
- [dupeguru](https://github.com/arsenetar/dupeguru), open source.
- [duplicatecleaner](https://www.duplicatecleaner.com)

The above tools are all very useful, but when the number of files is large, using the GUI to select and process duplicate files can be very inconvenient.

With the development of electronics technology, we generate a lot of files in work and life, especially photos and videos, as well as various electronic documents for work. Backing up this personal data is very important.

Backing up data is a relatively professional job. For most people, it is difficult to manage backups precisely. Backups are often made irregularly over different periods, with no rules for timing, quantity, or destinations. As a result, there are many backup copies, with no clear relationships between different backups, leading to a large number of duplicate and orphaned backups. Using the above tools to find and categorize duplicate files among hundreds of thousands of files is very time consuming and laborious.

`FileDiffGrouper` is a tool to quickly find identical and different files between two directories.

> FileDiffGrouper only finds duplicate and extra files in Target relative to Base. It does not deduplicate Base itself. Please use the tools listed above to deduplicate Base first. Of course, it is recommended to also deduplicate Target before using this tool, to improve efficiency.

## Install

3 way to install file difference grouper：

- Compile the source code:
  - `git clone https://github.com/jqk/file-diff-grouper.git`
  - `cd file-diff-grouper/file-diff-cli`
  - `go build`
- Download the package from <https://github.com/jqk/file-diff-grouper/releases>, unzip it and run.
- When [scoop](https://github.com/ScoopInstaller/Scoop) is installed, run:
  - `scoop bucket add ajqk https://github/jqk/scoopbucket`
  - `scoop install file-diff-grouper`

On Windows, the executable filename is `fdg.exe`, around 6MB in size.

In addition to the executable file, a sample config file is also provided: config.yaml.

## Usage

### Command line

Executing the `fdg` command is simple - just provide the full path to the configuration file.

```text {.line-numbers}
fdg c:\test\config.yaml
```

There are no specific requirements for the configuration filename itself, but it must have the correct file extension, see `Configuration` section for details.

```text {.line-numbers}
$ fdg

Copyright (c) 1999-2023 Not a dream Co., Ltd.
file difference grouper (fdg) 1.0.0, 2023-09-07

Usage:
  fdg [path/to/the/taskConfigFile]
      Compare and group the file differences according to specified config file.

  otherwise: show this help.
```

### Configuration

Since there are many parameters to specify, command line arguments are not used. Instead, a configuration file is used. The configuration file can be in formats supported by [viper](http://github.com/spf13/viper%7Cgithub.com/spf13/viper) such as `.json`, `.xml,` `.yml` and `.toml`. The example below uses `.yaml`. Please refer to the comments.

You don't need to read through in detail at first. Each parameter in the configuration file will be explained in more detail in the `How It Works and Parameter Descriptions` section.

```yaml {.line-numbers}
# The number of bytes read from each file header for calculating checksum during quick scan. 
headerSize: 1024
# The buffer size in bytes for reading files. It must be greater than or equal to HeaderSize. 
bufferSize: 10240

# Action to perform, case insensitive, must be one of: 
# - Compare: Compare compareBase and compareTarget.
# - ScanBase: Only scan compareBase.
# - ScanTarget: Only scan compareTarget. 
action: "Compare"

# Base path information for comparison.
compareBase:
  # Compare base path. Can be read-only. 
  dir: "z:/compare_base"
  # Filename to save the scan result.
  # Can be a relative or absolute path. Can also reference the path defined above using ${dir}. 
  # It equals to 'z:/compare_base/base.scan.json' below.
  scanResultFile: "${dir}/base.scan.json"
  # Whether to load previous scan result file to improve speed, avoiding rescanning every time.
  # Files in the compare paths may have changed (added, deleted etc) since last scan.
  # It's up to the user to decide whether to load previous result to improve efficiency.
  # The program does not automatically handle potential changes.
  # If set to true but the file does not exist, scanning will be performed. 
  loadScanResult: true
  # Whether to compute full file digests for the files being compared. 
  # False means computing digest of file header only.
  needFullChecksum: false

# 待比较的目标路径信息。
compareTarget:
  # Compare target path. Can be read-only if not doing moveMore and moveSame operations below. 
  # Otherwise must be writable.
  dir: "z:/compare_target"
  # Same as attributes defined in compareBase.
  scanResultFile: "z:/result/target.scan.json"
  loadScanResult: true
  needFullChecksum: false

  # Path to save the comparison results, must be writable. 
  # Can be a relative or absolute path. Can also reference the path defined in above using ${dir}.
  backupDir: "z:/result/group"
  # Whether to move files in target but not in base to the compare result dir. 
  # False means only generating a result file list.
  moveMore: false
  # Whether to move files that exist in both target and base to the compare result directory.
  # False means only generating a result file list.
  moveSame: false

# Filter criteria for selecting files to compare.
filter:
  # Whether to be case sensitive for file extensions.
  caseSensitive: false
  # File extensions to include. Must provide condition(s) with at least one valid string. 
  # Empty string means files without extension.
  # The include in this example covers major image and video extensions for phones and cameras.
  include:
    - "*.3gp" 
    - "*.amr" 
    - "*.avi" 
    - "*.bmp" 
    - "*.gif" 
    - "*.jpeg"
    - "*.jpg" 
    - "*.mov" 
    - "*.mp4" 
    - "*.mpg" 
    - "*.png" 
    - "*.webp"
    - "*.wmv" 

  # File extensions to exclude, can be left empty.
  exclude:
    - "*.log"

  # Minimum file size in bytes, 0 or less means no limit, 
  # but will process starting from 1 byte, 0 byte files are skipped.
  minFileSize: 1024

  # Maximum file size in bytes, 0 or less means no limit.
  maxFileSize: 0
```

### How it works and parameter descriptions

#### How it works

`fdg` traverses the directories specified in `compareBase.dir` and `compareTarget.dir`, and finds identical files (duplicate files) between them, as well as files that exist in Target but not in Base (extra files).

`fdg` does not compare filenames of the two files, but rather compares the file sizes and contents:

- If file sizes differ, the files are considered different.
- If file checksums differ, the files are considered different.

`fdg` first scans all files under `compareBase.dir` and `compareTarget.dir` including subdirectories, to get two scan result sets containing file sizes and checksums. It then compares the records in the two scan results based on the rules above to determine duplicate and extra files.

Currently the CRC32 algorithm is used, which should be sufficient.

#### headerSize & bufferSize

In order to calculate file checksums, the binary contents of each file needs to be read. Reading the entire contents of all files would take too much time, so `headerSize` is defined. For example, if there are 100 files of 1GB each, and `needFullChecksum` is set to true, 100GB of data will be read. If set to false and `headerSize` is 1024 bytes, only 100KB of data will be read, which is much faster.

`headerSize` should not be set too large, 1024 to 10240 is recommended. If `headerSize` is set smaller than 1024, it will be automatically adjusted to 1024 by the program.

`bufferSize` defines the buffer size for file IO, to improve speed. If `bufferSize` is smaller than headerSize, it will be automatically adjusted to the value of `headerSize`.

#### needFullChecksum

The checksum of the file header is called `headerChecksum`. If two files have the same length and `headerChecksum`, then their full file checksums `fullChecksum` need to be compared further.

If `fullChecksum` is not calculated, `fdg` will automatically calculate and save it in the scan result file. So in most cases, `needFullChecksum` can be set to false. `fdg` will automatically supplement the calculation as needed.

Setting `needFullChecksum` to true is useful in the scenario where there is a large directory that needs to be compared repeatedly with other directories. To avoid rescanning the directory every time, its complete scan result can be obtained once, and in subsequent runs `loadScanResult` can be set to true to save scanning time.

> For example, I have a USB drive with about 50,000 files totaling 300GB. After scanning it with `needFullChecksum` set to true and getting the result file `result.json`, I can then compare files on the USB drive with others using only `result.json` without connecting the USB drive.

#### loadScanResult & scanResultFile

Each comparison is based on the scan results of the two directories. The scan results are saved to a file defined by `scanResultFile`. If this value is an empty string, no scan result file will be output.

If `loadScanResult` is true and the file defined by `scanResultFile` exists, the scan results in that file will be loaded to save scanning time. Otherwise, scanning will be performed.

The scan results are saved in `JSON` format, with content like below:

```json {.line-numbers}
{
    "HeaderSize": 2000,
    "Dir": "test-data/origin/compare_base",
    "Filter": {
        "CaseSensitive": false,
        "Include": [
            "*.md",
            "*.txt"
        ],
        "Exclude": [
            "*.log"
        ],
        "MinFileSize": 1024,
        "MaxFileSize": 3072
    },
    "FileCount": 5,
    "FileSize": 9668,
    "HeaderChecksumCount": 3,
    "FullChecksumCount": 4,
    "ElapsedTime": 1050700,
    "Files": {
        "3096586316": [
            {
                "HeaderChecksum": 3096586316,
                "HasFullChecksum": true,
                "FullChecksum": 3096586316,
                "Filename": "e:\\github\\jqk\\file-diff-grouper\\file-diff\\test-data\\origin\\compare_base\\004.txt",
                "FileSize": 1588,
                "ModifiedTime": "2023-06-30T12:57:32.2270053+08:00"
            }
        ],
        "3222652411": [
            {
                "HeaderChecksum": 3222652411,
                "HasFullChecksum": false,
                "FullChecksum": 0,
                "Filename": "e:\\github\\jqk\\file-diff-grouper\\file-diff\\test-data\\origin\\compare_base\\001.md",
                "FileSize": 2278,
                "ModifiedTime": "2023-06-30T12:57:32.2260055+08:00"
            }
        ],
        "4245835769": [
            {
                "HeaderChecksum": 4245835769,
                "HasFullChecksum": true,
                "FullChecksum": 4245835769,
                "Filename": "e:\\github\\jqk\\file-diff-grouper\\file-diff\\test-data\\origin\\compare_base\\dir_0\\002.txt",
                "FileSize": 1934,
                "ModifiedTime": "2023-06-30T12:57:32.2270053+08:00"
            },
            {
                "HeaderChecksum": 4245835769,
                "HasFullChecksum": true,
                "FullChecksum": 4245835769,
                "Filename": "e:\\github\\jqk\\file-diff-grouper\\file-diff\\test-data\\origin\\compare_base\\dir_0\\dir_1\\003-same-as-002.md",
                "FileSize": 1934,
                "ModifiedTime": "2023-06-30T12:57:32.2270053+08:00"
            },
            {
                "HeaderChecksum": 4245835769,
                "HasFullChecksum": true,
                "FullChecksum": 4245835769,
                "Filename": "e:\\github\\jqk\\file-diff-grouper\\file-diff\\test-data\\origin\\compare_base\\dir_0\\dir_1\\copy-of-003.md",
                "FileSize": 1934,
                "ModifiedTime": "2023-06-30T12:57:32.2270053+08:00"
            }
        ]
    }
}
```

#### backupDir

Since the program is designed for cases with extremely large numbers of files, automatic deletion of duplicate files is not provided to avoid hard-to-recover mistakes. Instead, duplicate and extra files are moved to the specified directory for manual confirmation and deletion by the user.

`backupDir` specifies where to move the duplicate and extra files found. This value must be a valid writable path that has ***movable*** relationship with `compareTarget.dir`. After comparison, two result files will be kept in this directory:

- target-more-than-base.txt
- target-same-with-base.txt

It's obvious that `backupDir` must be writable. The requirement that it must be ***movable*** needs emphasis. Here "movable" means moving can be done without copying the file content. For example on Windows, moving `c:\doc\a.txt` to `c:\backup\a.txt` is extremely fast, without actually reading/writing the file content itself - it's similar to renaming the file. But moving it to `d:\doc\a.txt` would require first reading all content from `c:\doc\a.txt`, writing it to `d:\doc\a.txt`, and finally deleting `c:\doc\a.txt`. Considering there may be a huge number of large files, this would involve massive IO and waste time. Therefore, ***backupDir must have this kind of movable relationship with compareTarget.dir***.

The two result files have the same structure, for example:

```text {.line-numbers}
{
    "BaseDir": "test-data/origin/compare_base",
    "BaseFileCount": 6,
    "BaseHeaderChecksumCount": 4,
    "TargetDir": "test-data/origin/compare_target",
    "TargetFileCount": 3,
    "TargetHeaderChecksumCount": 3,
    "Filter": {
        "CaseSensitive": false,
        "Include": [
            "*.md",
            "*.txt"
        ],
        "Exclude": [
            "*.log"
        ],
        "MinFileSize": 1024,
        "MaxFileSize": 1048576
    },
    "CompareResultType": "more",
    "ResultFileCount": 2,
    "ResultFileSize": 35879,
    "ElapsedTime": 0
}

----------

"e:\github\jqk\file-diff-grouper\file-diff\test-data\origin\compare_target\013.md"
"e:\github\jqk\file-diff-grouper\file-diff\test-data\origin\compare_target\dir\011.txt"
```

Before the separator line `----------`, comparison result information is saved in JSON format. The meaning can be understood from the field names.

After the separator line, each line is the absolute path filename.

#### moveMore & moveSame

`moveMore` and `moveSame` specify whether to move the corresponding files to `backupDir`.
The program will create a directory named like `YYYYMMDD_HHMMSS` under `backupDir` based on current time, and then create `more` and `same` directories under it, for extra and duplicate files respectively.

The original directory structure will be kept when moving the files. For example, if `target/a/b.txt` is a duplicate file, it will be moved to `backupDir/20230907_123456/same/a/b.txt`. This makes it convenient to manually compare and locate the original files.

Here `20230907_123456` represents the execution time, which is `2023-09-07 12:34:56.`

The timestamped directory under `backupDir` isolates the result from multiple runs. The `more` and `same` subdirectories categorize the extra and duplicate files. Keeping the original structure helps identify where the files came from. This organization of the moved files aims to facilitate manual review and cleanup.

#### filter

`filter` defines the conditions to filter files. The comments in the config file have explained it clearly. Here are some additional notes:

- If the same extension appears in both `include` and `exclude`, those files will be skipped. This is because `exclude` is processed first.
- An empty string means files without any extension.

How to find out what file extensions exist in a directory? You can use a small tool called [file-extension-counter](https://github.com/jqk/FileExtensionCounter), only 3MB in size.

**Enjoy**!
