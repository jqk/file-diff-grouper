package filediff

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jqk/futool4go/fileutils"

	"github.com/spf13/viper"
)

/*
LoadConfigFromString loads a configuration from a string in the specified format.

Parameters:
  - content: The string containing the configuration data.
  - format: The format of the configuration data, such as yaml or json.

Returns:
  - a pointer to a Config struct
  - an error if there was a problem loading the configuration.
*/
func LoadConfigFromString(content string, format string) (*Config, error) {
	v := viper.New()
	v.SetConfigType(format)
	err := v.ReadConfig(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	return unmarshalConfig(v)
}

/*
LoadConfigFromString loads a configuration from a file.

Parameters:
  - filename: The path of the configuration file.

Returns:
  - a pointer to a Config struct
  - an error if there was a problem loading the configuration.
*/
func LoadConfigFromFile(filename string) (*Config, error) {
	fileExists, isDir, err := fileutils.FileExists(filename)
	if err != nil {
		return nil, err
	}
	if !fileExists {
		return nil, fmt.Errorf("file %s does not exist", filename)
	}
	if isDir {
		return nil, fmt.Errorf("file %s is a directory", filename)
	}

	dir, name, ext := splitFilename(filename)
	format := strings.ToLower(ext)

	v := viper.New()
	v.SetConfigName(name)
	v.SetConfigType(format)
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return unmarshalConfig(v)
}

func unmarshalConfig(v *viper.Viper) (*Config, error) {
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func splitFilename(filename string) (string, string, string) {
	dir, file := filepath.Split(filename)
	lastDot := strings.LastIndex(file, ".")
	name := file[:lastDot]
	ext := file[lastDot+1:]

	if dir == "" {
		dir = "./" // viper 装载文件时，dir 为空字符会找不到文件。
	}

	return dir, name, ext
}
