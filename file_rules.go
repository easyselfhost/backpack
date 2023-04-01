package backpack

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type FileCommand string

const (
	Copy   FileCommand = "copy"
	Sqlite FileCommand = "sqlite"
	Ignore FileCommand = "ignore"
)

type FileRule struct {
	Regex   *regexp.Regexp
	Command FileCommand
}

type DirRule struct {
	SrcDir    string     `json:"path"`
	FileRules []FileRule `json:"file_rules"`
}

func (fr *FileRule) UnmarshalJSON(data []byte) error {
	var raw struct {
		Regex   string `json:"regex"`
		Command string `json:"command"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("error parsing json object for FileRule: %w", err)
	}

	if raw.Regex == "" {
		return fmt.Errorf("error parsing json object: empty regexp")
	}
	re, err := regexp.Compile(raw.Regex)
	if err != nil {
		return fmt.Errorf("error parsing regexp in FileRule: %w", err)
	}

	if raw.Command != string(Copy) &&
		raw.Command != string(Sqlite) &&
		raw.Command != string(Ignore) {
		return fmt.Errorf("unsupported command in FileRule: %v", raw.Command)
	}

	fr.Command = FileCommand(raw.Command)
	fr.Regex = re
	return nil
}
