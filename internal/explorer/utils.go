package explorer

import (
	"encoding/json"
	"strings"
)

func sourceCodeContainsSetting(code string) bool {
	return strings.HasPrefix(code, "{{")
}

func sourceCodeNotContainsSetting(code string) bool {
	return strings.HasPrefix(code, "{")
}

func parseSourceCodeString(code string) (SourceCodeInfo, error) {
	var sources SourceCodeInfo
	err := json.Unmarshal([]byte(code), &sources)
	if err != nil {
		return SourceCodeInfo{}, err
	}
	return sources, nil
}

func contractLanguage(code string) string {
	if strings.HasPrefix(code, "vyper:") {
		return ".vy"
	}
	return ".sol"
}
