package backpack

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func SnapshotDir(rule DirRule, dest string) error {
	return filepath.Walk(rule.SrcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relName, err := filepath.Rel(rule.SrcDir, path)
		if err != nil {
			return err
		}

		fileCommand := getFileCommand(relName, rule.FileRules)

		return executeFileCommand(path, filepath.Join(dest, relName), fileCommand)
	})
}

func executeFileCommand(src string, dest string, cmd FileCommand) error {
	switch cmd {
	case Copy:
		return copyFile(src, dest)
	case Sqlite:
		return sqliteOnlineBackup(src, dest)
	case Ignore:
		return nil
	default:
		return fmt.Errorf("unsupported file command %v", cmd)
	}
}

func getFileCommand(fileName string, rules []FileRule) FileCommand {
	for _, rule := range rules {
		if rule.Regex.MatchString(fileName) {
			return rule.Command
		}
	}
	return Copy
}

func copyFile(src string, dest string) error {
	fin, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fin.Close()

	err = os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}

	fout, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)

	return err
}

func sqliteOnlineBackup(src string, dest string) error {
	err := os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}

	path := GetExecPath("sqlite3")
	cmd := exec.Command(path, src, fmt.Sprintf(".backup '%s'", dest))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error backing up sqlite with output with error %w, output: %s",
			err, string(output))
	}

	return nil
}
