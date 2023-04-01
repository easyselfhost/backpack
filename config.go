package backpack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Schedule struct {
	DailySchedule []string `json:"daily"`
	EveryInterval string   `json:"every"`
}

type BackupRule struct {
	Directories  []DirRule `json:"directories"`
	RcloneConfig string    `json:"rclone_remote"`
	RemotePath   string    `json:"remote_path"`
	Schedule     Schedule  `json:"schedule"`
}

type Config struct {
	Version     string       `json:"version"`
	BackupRules []BackupRule `json:"backups"`
}

func ParseConfigFromFile(path string) (Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read configuration file: %w", err)
	}

	return ParseConfigFromBytes(data)
}

func ParseConfigFromBytes(data []byte) (Config, error) {
	var conf Config

	err := json.Unmarshal(data, &conf)
	if err != nil {
		err = fmt.Errorf("error parsing json bytes %w", err)
	}

	return conf, err
}
