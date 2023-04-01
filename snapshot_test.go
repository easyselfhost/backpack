package backpack_test

import (
	"path/filepath"
	"regexp"
	"testing"

	bp "github.com/easyselfhost/backpack"
	bt "github.com/easyselfhost/backpack/testing"
	_ "github.com/mattn/go-sqlite3"
)

var snapshotFiles = map[string]string{
	"dir1/a.txt":  "file1",
	"dir2/b.conf": "file 2",
	"dir3/.c":     ".c",
	".d":          "another file",
}

const snapshotDbFiles = "dir1/path/db.sqlite3"

func TestSnaphost(t *testing.T) {
	srcDir, destDir, err := bt.GenTestDirs()
	if err != nil {
		t.Fatal(err)
	}

	err = bt.GenTextFiles(srcDir, snapshotFiles)
	if err != nil {
		t.Fatal(err)
	}

	err = bt.GenDbFile(filepath.Join(srcDir, snapshotDbFiles))
	if err != nil {
		t.Fatal(err)
	}

	err = bp.SnapshotDir(bp.DirRule{
		SrcDir: srcDir,
		FileRules: []bp.FileRule{
			{
				Regex:   regexp.MustCompile(".*\\.sqlite3"),
				Command: bp.Sqlite,
			},
			{
				Regex:   regexp.MustCompile("dir3/.*"),
				Command: bp.Ignore,
			},
		},
	}, destDir)
	if err != nil {
		t.Fatal(err)
	}

	bt.VerifyTextFile(t, filepath.Join(destDir, "dir1/a.txt"), snapshotFiles["dir1/a.txt"])
	bt.VerifyTextFile(t, filepath.Join(destDir, "dir2/b.conf"), snapshotFiles["dir2/b.conf"])
	bt.VerifyTextFile(t, filepath.Join(destDir, ".d"), snapshotFiles[".d"])
	bt.VerifyIgnoredFile(t, filepath.Join(destDir, "dir3/.c"))

	bt.VerifyDbFile(t, filepath.Join(destDir, snapshotDbFiles))
}
