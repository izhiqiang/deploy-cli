package conf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

// FileExists
// Determine whether the file exists
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err == nil {
		return info.Mode().IsRegular()
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// LoadFiles
// Retrieve one from many files
func LoadFiles(files ...string) (string, error) {
	var (
		confFile string
		err      error
	)
	for _, file := range files {
		if FileExists(file) {
			confFile = file
			break
		}
	}
	if confFile == "" {
		err = fmt.Errorf("loading %s failed", strings.Join(files, ","))
	}
	return confFile, err
}

func Unmarshal(filePath string, out any) (err error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(raw, out)
	return
}

func UnmarshalFiles(out any, files ...string) (err error) {
	configFile, err := LoadFiles(files...)
	if err != nil {
		return err
	}
	return Unmarshal(configFile, out)
}

func Write(file string, in any) error {
	out, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return os.WriteFile(file, out, os.ModePerm)
}
