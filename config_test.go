package backpack_test

import (
	"reflect"
	"regexp"
	"testing"

	bp "github.com/easyselfhost/backpack"
)

func TestParseConfigFromBytes(t *testing.T) {
	configData := []byte(`
	{
		"version": "0.1",
		"backups": [
			{
				"directories": [
					{
						"path": "/data/dir",
						"file_rules": [
							{
								"regex": ".*\\.sqlite3",
								"command": "sqlite"
							},
							{
								"regex": ".*",
								"command": "copy"
							}
						]
					},
					{
						"path": "/data/dir2",
						"file_rules": [
							{
								"regex": ".*",
								"command": "copy"
							}
						]
					}
				],
				"rclone_remote": "default",
				"remote_path": "my_bucket",
				"schedule": {
					"daily": [
						"03:00",
						"15:00"
					],
					"every": "30m"
				}
			}
		]
	}
	`)

	expected := bp.Config{
		Version: "0.1",
		BackupRules: []bp.BackupRule{
			{
				Directories: []bp.DirRule{
					{
						SrcDir: "/data/dir",
						FileRules: []bp.FileRule{
							{
								Regex:   regexp.MustCompile(".*\\.sqlite3"),
								Command: bp.Sqlite,
							},
							{
								Regex:   regexp.MustCompile(".*"),
								Command: bp.Copy,
							},
						},
					},
					{
						SrcDir: "/data/dir2",
						FileRules: []bp.FileRule{
							{
								Regex:   regexp.MustCompile(".*"),
								Command: bp.Copy,
							},
						},
					},
				},
				RcloneConfig: "default",
				RemotePath:   "my_bucket",
				Schedule: bp.Schedule{
					DailySchedule: []string{
						"03:00",
						"15:00",
					},
					EveryInterval: "30m",
				},
			},
		},
	}

	config, err := bp.ParseConfigFromBytes(configData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(config, expected) {
		t.Errorf("result not equal expected: %v, actual: %v", expected, config)
	}
}
