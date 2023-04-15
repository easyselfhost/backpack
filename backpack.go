package backpack

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Workflow interface {
	Run() error
}

type RetryingWorkflow struct {
	workflow Workflow
	retries  uint
}

func NewRetryingWorkflow(workflow Workflow, retries uint) Workflow {
	return &RetryingWorkflow{
		workflow: workflow,
		retries:  retries,
	}
}

func (wf *RetryingWorkflow) Run() error {
	err := wf.workflow.Run()

	if err == nil {
		return nil
	}

	for i := uint(0); i < wf.retries; i++ {
		err = wf.workflow.Run()
		if err == nil {
			return nil
		}
	}

	return err
}

type BackpackFlow struct {
	uploadRule BackupRule
	dirRules   map[string]DirRule
}

func NewBackpackFlow(uploadRule BackupRule) Workflow {
	dirRules := map[string]DirRule{}

	for _, dirRule := range uploadRule.Directories {
		dirRules[dirRule.SrcDir] = dirRule
	}

	return &BackpackFlow{
		uploadRule: uploadRule,
		dirRules:   dirRules,
	}
}

func (bf *BackpackFlow) Run() (err error) {
	tmpDir, err := ioutil.TempDir("", "backpack-tmp-dir-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}

	defer func() {
		rerr := os.RemoveAll(tmpDir)
		if rerr != nil {
			err = fmt.Errorf("failed to remove temp directory %w", rerr)
		}
	}()

	// backing up all directories
	for _, dirRule := range bf.dirRules {
		dirName := filepath.Base(dirRule.SrcDir)
		destDir := filepath.Join(tmpDir, dirName)
		err = os.Mkdir(destDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		err = SnapshotDir(dirRule, destDir)
		if err != nil {
			return fmt.Errorf("failed to backup %s: %w", dirRule.SrcDir, err)
		}
	}

	err = UploadDir(bf.uploadRule, tmpDir)
	if err != nil {
		return fmt.Errorf("failed to uploade directories: %w", err)
	}

	return nil
}
