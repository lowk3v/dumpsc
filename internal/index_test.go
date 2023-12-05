package internal

import (
	global "github.com/lowk3v/dumpsc/config"
	"os"
	"strings"
	"testing"
)

func TestOptions_Run(t *testing.T) {
	//projectDir := os.Getenv("PROJECT_DIR")
	projectDir := "/Users/lap14962/projects/me/dumpsc"
	if global.CustomConfig(projectDir+"/config.yml") != nil {
		os.Exit(0)
	}

	defaultOption := &Options{
		Explorer: "etherscan",
		Address:  "",
		ApiKey:   "",
		Output:   projectDir + "/src",
		Url:      "url",
	}

	tests := []struct {
		name   string
		config Options
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "Fail. Incorrect url 1",
			config: Options{
				Explorer: defaultOption.Explorer,
				Address:  defaultOption.Address,
				ApiKey:   defaultOption.ApiKey,
				Output:   defaultOption.Output,
				Url:      "https://incorrect/url",
			},
			want: "the url is not valid",
		},
		{
			name: "Fail. Incorrect url 2",
			config: Options{
				Explorer: defaultOption.Explorer,
				Address:  defaultOption.Address,
				ApiKey:   defaultOption.ApiKey,
				Output:   defaultOption.Output,
				Url:      "https://etherscan.io/address/",
			},
			want: "the url is not valid",
		},
		{
			name: "Fail. no support explorer",
			config: Options{
				Explorer: defaultOption.Explorer,
				Address:  defaultOption.Address,
				ApiKey:   defaultOption.ApiKey,
				Output:   defaultOption.Output,
				Url:      "https://incorrect.io/address/0x9abf23f4e439d695a7fd341a1b25873c50cfa52e",
			},
			want: "explorer is not supported",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New(&tt.config)
			_, _, err := app.RunTest()
			if !ErrorContains(err, tt.want) {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func ErrorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}
