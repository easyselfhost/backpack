package testing

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const FileContent = "this is a file."

const CreateTableQuery = `
	CREATE TABLE my_table (
		id INTEGER PRIMARY KEY,
		name TEXT,
		age INTEGER
	);
`
const InsertQuery = `
	INSERT INTO my_table (name, age)
	VALUES ('John Smith', 35);
`
const SelectQuery = `
	SELECT * FROM my_table WHERE name = 'John Smith';
`

var TestFiles = []string{
	"dir1/a.txt",
	"dir1/db/db.sqlite3",
	"dir2/b.conf",
	"dir3/.c",
	".d",
}

func VerifyTextFile(t *testing.T, path string, content string) {
	contentRead, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to open file %s: %v", path, err)
	}

	if string(contentRead) != content {
		t.Errorf("file %s does not have expected content", path)
	}
}

func VerifyIgnoredFile(t *testing.T, path string) {
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Error("ignored file is not ignored")
	}
}

func VerifyDbFile(t *testing.T, path string) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(SelectQuery)
	if err != nil {
		t.Fatal(err)
	}

	var id, age int
	var name string
	err = stmt.QueryRow().Scan(&id, &name, &age)
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()

	if id != 1 || name != "John Smith" || age != 35 {
		t.Error("sqlite file doesn't have the written row")
	}
}

func GenTestDirs() (srcDir string, destDir string, err error) {
	srcDir, err = ioutil.TempDir("", "test-src-*")
	if err != nil {
		return "", "", err
	}

	destDir, err = ioutil.TempDir("", "test-dst-*")
	if err != nil {
		return "", "", err
	}

	return
}

func GenDbFile(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(CreateTableQuery + InsertQuery)
	return err
}

func GenTextFiles(dir string, fileContent map[string]string) error {
	for name, content := range fileContent {
		path := filepath.Join(dir, name)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}

		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.WriteString(content)
		if err != nil {
			return err
		}
	}

	return nil
}
