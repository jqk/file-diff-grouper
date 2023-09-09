# File Diff Grouper

[English](readme.md)

## 一、 功能

FileDiffGrouper 是一个比较两个目录文件差异的命令行工具。它以二进制方式比较两个目录中的所有文件的内容，而不是文件名。虽然仅在 Windows 10/11 中测试过，但未使用特定的操作系统功能，所以理论上可以在运行于 Linux 和 MacOS。

FileDiffGrouper 比较两个目录，一个称为 Base 目录，一个称为 Target 目录。FileDiffGrouper 将给出两个结果集合，保存在结果文件中：

- MORE 集合：Target 比 Base 多的文件，即 Target 中有，但 Base 中没有的文本，本文中称为`多出文件`。
- SAME 集合：Target 和 Base 中同时存在的完全相同的文件，本文中称为`重复文件`。

FileDiffGrouper 可以根据选项，直接将以上两个集合的文件，移动到指定的备份目录中。为保证安全，不提供自动删除这些文件的功能。

## 二、 解决的问题

有许多软件提供了比较两个目录中文件异同、查找重复文件的功能，如：

- [alldup](https://alldup.info)，免费。
- [AntiDupl](https://github.com/ermig1979/AntiDupl)，开源。
- [beyondcompare](https://www.scootersoftware.com)
- [czkawka](https://github.com/qarmin/czkawka)，开源，主要查找重复图片，还有 GUI 版。
- [dupeguru](https://github.com/arsenetar/dupeguru)，开源。
- [duplicatecleaner](https://www.duplicatecleaner.com)

以上各个软件都很好用，但当文件数量比较多时，使用 GUI 选择处理重复文件就很不方便。

随着电子技术的发展，我们在工作和生活会产生很多的文件，尤其以照片、视频为主，当然还有工作中的各类电子文档。备份这些个人数据是非常重要的。

备份数据是一件相对专业的工作，对于大多数人而言，很难做到精确的管理，经常会有多个不同时期做的备份，备份时间周期、数量、目标无规则。因此备份数量多，备份与备份之间没有明确的关系，会有大量重复备份及孤本备份。使用上述软件在数十万个文件中找出重复文件并分组处理是一件耗时耗力的事。

FileDiffGrouper 就是这样一款快速找出两个目录之间相同与不同文件的工具。

> FileDiffGrouper 只会找出 Target 相对于 Base 的重复文件和多出文件，不会对 Base 本身查重。请先使用前面列出的工具对 Base 查重。当然，建议在使用本工具前也对 Target 查重，以提高工作效率。

## 三、 安装

三种安装方法：

- 编译源码:
  - `git clone https://github.com/jqk/file-diff-grouper.git`
  - `cd file-diff-grouper/file-diff-cli`
  - `go build`
- 从 <https://github.com/jqk/file-diff-grouper/releases> 下载解压软件包，直接运行即可。
- 通过 [scoop](https://github.com/ScoopInstaller/Scoop) 安装。在安装好 scoop 后：
  - `scoop bucket add ajqk https://github/jqk/scoopbucket`
  - `scoop install file-diff-grouper`

在 Windows 下，可执行文件名为 `fdg.exe`，约 6MB。

## 四、 使用

### 4.1 命令行

执行 `fdg` 命令很简单，仅需提供配置文件的完整路径，如 `fdg c:\test\config.yaml`。对配置文件名本身没有特别要求，但必须有正确格式的扩展名，见`配置文件`说明。

```text {.line-numbers}
$ fdg

Copyright (c) 1999-2023 Not a dream Co., Ltd.
file difference grouper (fdg) 1.0.0, 2023-09-07

Usage:
  fdg [path/to/the/taskConfigFile]
      Compare and group the file differences according to specified config file.

  otherwise: show this help.
```

### 4.2 配置文件

由于要指定的参数很多，就不使用命令行参数的形式，而使用配置文件。配置文件可以是 `.json`、`.xml`、`.yml` 和 `.toml` 等 [viper](github.com/spf13/viper) 支持的格式。以下以 `yaml` 为例说明，请参见注释。

一开始不必细读，配置文件里各个参数会在`工作原理及参数说明`中更详细的说明。

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

### 4.3 工作原理及参数说明

#### 4.3.1 工作原理

`fdg` 遍历 `compareBase.dir` 和 `compareTarget.dir` 指定的目录，找出两者相同的文件（重复文件），以及 Target 比 Base 多的文件（多出文件）。

`fdg` 不比较两个文件的文件名，而是比较两个文件的长度和内容：

- 文件长度不同，则认为文件不同。
- 文件的校验和不同，则认为文件不同。

`fdg` 先扫描 `compareBase.dir` 和 `compareTarget.dir` 及其子目录中的所有文件，得到两个包含文件长度及校验和的扫描结果，然后根据以上规则对比扫描结果中的记录，确定重复文件多出文件。

目前使用 `CRC32` 算法，应已够用。

#### 4.3.2 headerSize 和 bufferSize

为了得到文件的校验和，必需读取文件的二进制内容进行计算。而读取所有文件的整个文件内容，将耗费太多时间，因此定义了 `headerSize`。例如有 100 个 1GB 的文件，如果 `needFullChecksum` 设置为 true，将读取 100GB 的数据；而将其设为 false，并设置 `headerSize` 为 1024 字节，则仅需读取 100KB 的数据。后者速度要远远高于前者。

`headerSize` 不易设得过大，建议为 1024 到 10240。如果 `headerSize` 设置小于 1024，程序会自动将其调整为 1024。

`bufferSize` 定义的是读文件缓冲区的长度，以提高 IO 速度。如果 `bufferSize` 小于 `headerSize`，程序会自动将其调整为 `headerSize` 的值。

#### 4.3.3 needFullChecksum

文件头的校验和命名为 `headerChecksum`。如果两个文件的长度及 `headerChecksum` 相同，则需进一步比较其整个文件的校验和 `fullChecksum`。

如果未计算整个文件的校验和 `fullChecksum`，`fdg` 将自动计算并保存的扫描结果文件中。所以一般情况下，`needFullChecksum` 设置为 false 即可。`fdg` 会根据需要自动补充计算。

将 `needFullChecksum` 设置为 true 场景是，有一个比较大的目录，需要和其它多个目录反复比较，为避免每次都扫描该目录，则可一次性得到整个目录的完整扫描结果，每次都通过设置 `loadScanResult` 为 true，节约扫描时间。

> 例如，我有个 U 盘，上面有约 5 万个文件，共约 300GB。将 `needFullChecksum` 设置为 true 扫描后得到结果文件 `result.json`。以后我可以只使用 `result.json` 而不必连接该 U 盘即能完成针对该 U 盘上文件的比较。

#### 4.3.4 loadScanResult 和 scanResultFile

每次比较，均基于两个目录的扫描结果。扫描结果将保存在由 `scanResultFile` 定义的文件中。如果该值定义为空字符串，则不输出扫描结果文件。

`loadScanResult` 为 true，且 `scanResultFile` 定义的文件存在，则装载该文件的扫描结果，从而节省扫描时间；否则执行扫描。

扫描结果以 `json` 格式保存，内容如下：

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

#### 4.3.5 backupDir

由于程序主要针对文件量极大的情况设计，因此为避免自动删除重复文件产生不易挽回的错误，所以不提供自动删除功能，而是将重复文件以及多出文件移动到指定的目录中，用户可确认后手动删除。

`backupDir` 指定重复文件及多出文件的移动位置。该值必须为有效的位置，并且是**可写**的且**可移动**。比较完成后会有两个比较结果文件保留到该目录中：

- target-more-than-base.txt
- target-same-with-base.txt

`backupDir` 必须可写不必多说，必须**可移动**需要再次强调。此处的可移动指的是无需复制文件。以 Windows 为例，将 `c:\doc\a.txt` 移动到 `c:\backup\a.txt` 是极为迅速的，没有针对文件本身的读取与写入，类似于修改文件名；而将其移动到 `d:\doc\a.txt`，要先读取 `c:\doc\a.txt` 的全部内容，再将其写入到 `d:\doc\a.txt`，最后再将 `c:\doc\a.txt` 删除。考虑到文件可能很多、很大，这将涉及大量 IO，浪费时间。所以，**backupDir 一定要与 compareTarget.dir 有这种可移动关系**。

这两个文件的结构相同，示例如下：

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

在分隔线 `----------` 之前是以 json 格式保存的部分参数及比较结果信息，根据字段名基本可以了解其意义，不再赘述。

在分隔线之后是该类文件的绝对路径文件名，每个文件占据一行。

#### 4.3.6 moveMore 和 moveSame

`moveMore` 和 `moveSame` 指定是否将对应的文件移动到 `backupDir`。程序将在 `backupDir` 下根据时间创建名如 `YYYYMMDD_HHMMSS` 的目录，然后再在其下创建 `more` 和 `same` 目录，分别用于重复文件和多出文件。

移动操作将保持原有的目录结构，例如 `target/a/b.txt` 是重复文件，则将移动到 `backupDir/20230907_123456/same/a/b.txt`。这样方便人工对比，确定原有文件的位置。其中 `20230907_123456` 代表程序执行的时间点，2023 年 9 月 7 日 12 点 34 分 56 秒。

#### 4.3.7 filter

`filter` 定义过滤文件的条件，在前面配置文件中的注释已足够清晰。需要特别说明的是：

- 如果同一扩展名同时出现在 `include` 和 `exclude` 中，该类文件将不会被选中，因为程序先处理 `exclude` 中定义的扩展名。
- 空字符串表示文件没有扩展名。

> 如何确定在目录中有哪些扩展名呢？可以使用名为 [file-extension-counter](https://github.com/jqk/FileExtensionCounter) 的小工具，仅 3MB。

**Enjoy**!
