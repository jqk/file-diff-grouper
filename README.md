# File Diff Grouper

[中文](readme-CN.md)

## What is it

FileDiffGrouper is a command line tool for comparing file differences between two directories. It compares all files in the two directories byte-by-byte based on content, rather than filenames. Although only tested on Windows 10/11, it does not use any OS-specific features, so theoretically can run on Linux and MacOS as well.

FileDiffGrouper compares two directories - one called the Base directory, the other called the Target directory. FileDiffGrouper will output two result sets, saved to result files:

- MORE set: Files that exist in Target but not in Base, referred to as "extra files" in this document.
- SAME set: Files that exist in both Target and Base and have identical content, referred to as "duplicate files" in this document.

FileDiffGrouper can optionally move the files in the MORE and SAME sets to specified backup directories, based on options. For safety, it does not provide functionality to automatically delete these files.

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

FileDiffGrouper is a tool to quickly find identical and different files between two directories.

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

## Usage

### Command line

Executing the `fdg` command is simple - just provide the full path to the configuration file, e.g. `fdg c:\test\config.yaml`. There are no specific requirements for the configuration filename itself, but it must have the correct file extension, see `Configuration` section for details.

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

Since there are many parameters to specify, command line arguments are not used. Instead, a configuration file is used. The configuration file can be in formats supported by [viper](http://github.com/spf13/viper%7Cgithub.com/spf13/viper) such as `.json`, `.xml,` `.yml` and `.toml`. The example below uses `yaml`. Please refer to the comments.

You don't need to read through in detail at first. Each parameter in the configuration file will be explained in more detail in the `How It Works and Parameter Descriptions` section.

```yaml {.line-numbers}
# 快速扫描时读取各个文件头用于计算摘要的字节长度。
headerSize: 1024
# 文件读取缓冲区字节长度，必须大于等于 HeaderSize。
bufferSize: 10240

# 待执行的操作，大小写不敏感，只能是以下三种之一：
# - Compare：比较 compareBase 和 compareTarget。
# - ScanBase：仅扫描 compareBase。
# - ScanTarget：仅扫描 compareTarget。
action: "Compare"

# 待比较的基准路径信息。
compareBase:
  # 待比较的文件所在路径。可以是只读的。
  dir: "z:/compare_base"
  # 保存扫描结果的文件名，必须是可写的。
  # 可以指定相对或绝对路径文件名。也可以通过 ${dir} 引用 dir 定义的路径。
  scanResultFile: "z:/result/base.scan.json"
  # 是否直接装载以前的扫描结果的文件，以提高速度，避免每次重复扫描。
  # 在上一次扫描到本次扫描之间，待比较路径中的文件可能有增、删等变化。
  # 这要由用户自行决定是否装载以前的扫描结果以提高效率。程序不会自动关注可能存在的变化。
  # 如果为 true，但文件不存在，则会进行扫描。
  loadScanResult: true
  # 是否对待比较的文件进行整个文件摘要的计算。为 false 计算文件头的摘要值。 
  needFullChecksum: false

# 待比较的目标路径信息。
compareTarget:
  # 比较目标路径，如果不做下面设定的 moveMore 和 moveSame 操作，可以是只读的。否则，必须是可写的。
  dir: "z:/compare_target"
  # 保存扫描结果的文件名，必须是可写的。
  # 可以指定相对或绝对路径文件名。也可以通过 ${dir} 引用 dir 定义的路径。
  # 如以下路径等同于 "z:/compare_target/target.scan.json"。
  scanResultFile: "${dir}/target.scan.json"
  # 是否直接装载以前的扫描结果的文件，以提高速度，避免每次重复扫描。
  # 在上一次扫描到本次扫描之间，待比较路径中的文件可能有增、删等变化。
  # 这要由用户自行决定是否装载以前的扫描结果以提高效率。程序不会自动关注可能存在的变化。
  # 如果为 true，但文件不存在，则会进行扫描。
  loadScanResult: true
  # 是否对待比较的文件进行整个文件摘要的计算。为 false 计算文件头的摘要值。 
  needFullChecksum: false

  # 保存比较结果的路径，必须是可写的。而且，必需是无需复制，能直接从 dir 移动文件的路径。
  # 可以指定相对或绝对路径。也可以通过 ${dir} 引用 dir 定义的路径。
  backupDir: "z:/result/group"
  # 是否将 target 比 base 多的文件移动到比较结果目录中。为 false 只生成结果文件列表。
  moveMore: false
  # 是否将 target 和 base 相同的文件移动到比较结果目录中。为 false 只生成结果文件列表。
  moveSame: false

# 选择待比较文件的过滤条件。
filter:
  # 扩展名是否大小写敏感。
  caseSensitive: false
  # 包含的文件扩展名。必须提供至少一个有效的字符串的条件。空字符串表示没有扩展名的文件。
  # 本例中的 include 是手机、相机中主要的图片、视频扩展名。
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

  # 排除的文件扩展名，可以不提供。
  exclude:
    - "*.log"

  # 文件字节最小长度小于等于 0 表示不限制，但至少会从 1 字节开始，0 字节不处理。
  minFileSize: 1024

  # 文件字节最大长度小于等于 0 表示不限制。
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
