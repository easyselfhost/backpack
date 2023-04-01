package backpack_test

import (
	"path/filepath"
	"testing"

	bp "github.com/easyselfhost/backpack"
	bt "github.com/easyselfhost/backpack/testing"
	_ "github.com/rclone/rclone/backend/local"
)

var uploadFiles = map[string]string{
	"dir1/a.txt":  "file1",
	"dir2/b.conf": "file 2",
	"dir3/.c":     ".c",
	".d":          "another file",
}

func TestUpload(t *testing.T) {
	srcDir, destDir, err := bt.GenTestDirs()
	if err != nil {
		t.Fatal(err)
	}

	err = bt.GenTextFiles(srcDir, uploadFiles)
	if err != nil {
		t.Fatal(err)
	}

	backupRule := bp.BackupRule{
		Directories: []bp.DirRule{},
		RemotePath:  destDir,
	}

	err = bp.UploadDir(backupRule, srcDir)
	if err != nil {
		t.Error(err)
	}

	for name, content := range uploadFiles {
		bt.VerifyTextFile(t, filepath.Join(destDir, name), content)
	}
}
