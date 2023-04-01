package backpack

import (
	"fmt"
	"os"
	"strings"
)

func GetExecPath(name string) string {
	if p := os.Getenv(fmt.Sprintf("%s_PATH", strings.ToUpper(name))); p != "" {
		return p
	}
	return name
}
