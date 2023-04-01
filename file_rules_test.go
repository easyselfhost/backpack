package backpack_test

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	bp "github.com/easyselfhost/backpack"
)

func TestFileRule_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		errSubStr      string
		expectedResult bp.FileRule
	}{
		{
			name: "basic1",
			args: args{
				data: []byte("{\"regex\": \".*\\\\.db\", \"command\": \"sqlite\"}"),
			},
			wantErr: false,
			expectedResult: bp.FileRule{
				Regex:   regexp.MustCompile(".*\\.db"),
				Command: bp.Sqlite,
			},
		},
		{
			name: "basic2",
			args: args{
				data: []byte("{\"regex\": \".*\", \"command\": \"copy\"}"),
			},
			wantErr: false,
			expectedResult: bp.FileRule{
				Regex:   regexp.MustCompile(".*"),
				Command: bp.Copy,
			},
		},
		{
			name: "basic3",
			args: args{
				data: []byte("{\"regex\": \"abc\\\\.conf\", \"command\": \"ignore\"}"),
			},
			wantErr: false,
			expectedResult: bp.FileRule{
				Regex:   regexp.MustCompile("abc\\.conf"),
				Command: bp.Ignore,
			},
		},
		{
			name: "invalid field",
			args: args{
				data: []byte("{\"regexp\": \"*.db\", \"command\": \"sqlite\"}"),
			},
			wantErr:   true,
			errSubStr: "error parsing json object: empty regexp",
		},
		{
			name: "invalid regex",
			args: args{
				data: []byte("{\"regex\": \"*.db\", \"command\": \"sqlite\"}"),
			},
			wantErr:   true,
			errSubStr: "error parsing regexp",
		},
		{
			name: "invalid command",
			args: args{
				data: []byte("{\"regex\": \".*\", \"command\": \"cp\"}"),
			},
			wantErr:   true,
			errSubStr: "unsupported command",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := &bp.FileRule{}
			var err error
			if err = fr.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("FileRule.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), tt.errSubStr) {
				t.Errorf("FileRule.UnmarshalJson() error does not contain \"%v\", actual: %v", tt.errSubStr, err)
			}
			if err == nil && !reflect.DeepEqual(*fr, tt.expectedResult) {
				t.Errorf("FileRule.UnmarshalJSON() result not equal actual = %v, expected %v", *fr, tt.expectedResult)
			}
		})
	}
}
