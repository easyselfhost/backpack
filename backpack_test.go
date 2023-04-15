package backpack_test

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	bp "github.com/easyselfhost/backpack"
	bt "github.com/easyselfhost/backpack/testing"
	"github.com/golang/mock/gomock"
)

var backpackFiles = map[string]string{
	"dir1/a.txt":        "file1",
	"dir1/b.conf":       "conf",
	"dir2/b.conf":       "file 2",
	"dir2/fake.sqlite3": "fake db",
	"dir2/dir3/.c":      ".c",
	"dir2/.d":           "another file",
}

const backpackDbFiles = "dir1/path/db.sqlite3"

func TestBackpackFlow(t *testing.T) {
	srcDir, destDir, err := bt.GenTestDirs()
	if err != nil {
		t.Fatal(err)
	}

	err = bt.GenTextFiles(srcDir, backpackFiles)
	if err != nil {
		t.Fatal(err)
	}

	err = bt.GenDbFile(filepath.Join(srcDir, backpackDbFiles))
	if err != nil {
		t.Fatal(err)
	}

	wf := bp.NewBackpackFlow(bp.BackupRule{
		Directories: []bp.DirRule{
			{
				SrcDir: filepath.Join(srcDir, "dir1"),
				FileRules: []bp.FileRule{
					{
						Regex:   regexp.MustCompile(".*\\.sqlite3"),
						Command: bp.Sqlite,
					},
					{
						Regex:   regexp.MustCompile(".*b\\.conf"),
						Command: bp.Ignore,
					},
				},
			},
			{
				SrcDir: filepath.Join(srcDir, "dir2"),
			},
		},
		RemotePath: destDir,
	})

	err = wf.Run()
	if err != nil {
		t.Fatal(err)
	}

	bt.VerifyTextFile(t, filepath.Join(destDir, "dir1/a.txt"), backpackFiles["dir1/a.txt"])
	bt.VerifyDbFile(t, filepath.Join(destDir, backpackDbFiles))
	bt.VerifyIgnoredFile(t, filepath.Join(destDir, "dir1/b.conf"))

	// Verify files in dir2/ are all copied
	for name, content := range backpackFiles {
		if !strings.HasPrefix(name, "dir2") {
			continue
		}

		bt.VerifyTextFile(t, filepath.Join(destDir, name), content)
	}
}

func TestRetryingWorkflow_Run(t *testing.T) {
	type fields struct {
		retries uint
		fails   uint
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "single success",
			fields: fields{
				retries: 0,
				fails:   0,
			},
			wantErr: false,
		},
		{
			name: "single success with retries",
			fields: fields{
				retries: 10,
				fails:   0,
			},
			wantErr: false,
		},
		{
			name: "errors eventually success",
			fields: fields{
				retries: 3,
				fails:   3,
			},
			wantErr: false,
		},
		{
			name: "errors",
			fields: fields{
				retries: 3,
				fails:   4,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockwf := bt.NewMockWorkflow(ctrl)

			wf := bp.NewRetryingWorkflow(mockwf, tt.fields.retries)

			if tt.fields.fails >= tt.fields.retries+1 {
				mockwf.EXPECT().Run().DoAndReturn(func() error {
					return errors.New("test error")
				}).MinTimes(int(tt.fields.retries) + 1).MaxTimes(int(tt.fields.retries) + 1)
			} else {
				fc := mockwf.EXPECT().Run().DoAndReturn(func() error {
					return errors.New("tes terror")
				}).MaxTimes(int(tt.fields.fails))
				mockwf.EXPECT().Run().DoAndReturn(func() error {
					return nil
				}).After(fc)
			}

			if err := wf.Run(); (err != nil) != tt.wantErr {
				t.Errorf("RetryingWorkflow.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
