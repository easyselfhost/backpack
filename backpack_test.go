package backpack_test

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	bp "github.com/easyselfhost/backpack"
	bt "github.com/easyselfhost/backpack/testing"
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
