package explorer

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"golang.org/x/crypto/sha3"
	"reflect"
	"strings"
)

func containsSetting(code string) bool {
	return strings.HasPrefix(code, "{{")
}

func notContainsSetting(code string) bool {
	return strings.HasPrefix(code, "{")
}

func parseSourceCodeToStruct(code string) (SourceCodeInfo, error) {
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

func parseSourceCodeRaw(sourceCodeRaw string) string {
	trimSourceCode := ""
	if containsSetting(sourceCodeRaw) {
		trimSourceCode = strings.ReplaceAll(strings.ReplaceAll(sourceCodeRaw, "{{", "{"), "}}", "}")
	} else if notContainsSetting(sourceCodeRaw) {
		trimSourceCode = sourceCodeRaw
	}
	return trimSourceCode
}

func isDataError(apiResponse *resty.Response, response *ApiResponse) error {
	if apiResponse.StatusCode() != 200 {
		return errors.New(apiResponse.Status())
	}
	if response.Status != "1" {
		return errors.New(apiResponse.String())
	}
	if response.Results == nil {
		return errors.New("no result")
	}

	if reflect.TypeOf(response.Results).Kind() == reflect.String &&
		response.Results.(string) == "Contract source code not verified" {
		return errors.New("the contract source code is not verified")
	}
	return nil
}
