package backpack

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rclone/rclone/librclone/librclone"
)

func init() {
	librclone.Initialize()
}

type syncRequest struct {
	SrcFs string `json:"srcFs"`
	DstFs string `json:"dstFs"`
}

func UploadDir(rule BackupRule, dir string) error {
	var dstFs string
	if rule.RcloneConfig == "" {
		dstFs = rule.RemotePath
	} else {
		dstFs = fmt.Sprintf("%s:%s", rule.RcloneConfig, rule.RemotePath)
	}

	req := syncRequest{
		SrcFs: dir,
		DstFs: dstFs,
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to encode JSON request: %w", err)
	}

	out, status := librclone.RPC("sync/sync", string(reqJson))
	if status != http.StatusOK {
		return fmt.Errorf("failed to call rclone, output: %s", out)
	}

	return nil
}
