package explorer

import (
	"encoding/json"
	"golang.org/x/crypto/sha3"
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

// ChecksumAddress Convert address into checksum address
func ChecksumAddress(address string) string {
	hex := strings.ToLower(address)[2:]

	d := sha3.NewLegacyKeccak256()
	d.Write([]byte(hex))
	hash := d.Sum(nil)

	ret := "0x"

	for i, b := range hex {
		c := string(b)
		if b < '0' || b > '9' {
			if hash[i/2]&byte(128-i%2*120) != 0 {
				c = string(b - 32)
			}
		}
		ret += c
	}

	return ret
}
