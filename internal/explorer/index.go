package explorer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"reflect"
	"strings"
)

type Explorer struct {
	ApiGetSourceCode string `yaml:"apiGetSourceCode"`
	ApiKey           string `yaml:"apiKey"`
}

func New(e *Explorer) *Explorer {
	return e
}

func (e *Explorer) GetSourceCode(address string, proxyDepth int) ([]ContractFile, error) {
	fullUrl := fmt.Sprintf("%s?module=contract&action=getsourcecode&address=%s&apikey=%s", e.ApiGetSourceCode, address, e.ApiKey)
	if proxyDepth > 3 {
		// Should not be more than 3
		proxyDepth = 3
	}

	var response ApiResponse
	httpClient := resty.New()
	rawResp, err := httpClient.R().SetResult(&response).Get(fullUrl)
	if err != nil {
		return nil, err
	}
	if response.Status != "1" {
		return nil, errors.New(rawResp.String())
	}
	if len(response.Results) == 0 {
		return nil, errors.New("no result")
	}

	files := []ContractFile{
		{
			Name:    "address.txt",
			Content: address,
		},
	}

	// classify source code
	for _, result := range response.Results {
		if result.SourceCode == "" ||
			(result.ContractName == "" && result.ABI == "Contract source code not verified") {
			return nil, errors.New("the contract source code is not verified")
		}

		// source code has settings
		trimSourceCode := ""
		if sourceCodeContainsSetting(result.SourceCode) {
			trimSourceCode = strings.ReplaceAll(strings.ReplaceAll(result.SourceCode, "{{", "{"), "}}", "}")
		} else if sourceCodeNotContainsSetting(result.SourceCode) {
			trimSourceCode = result.SourceCode
		}

		if trimSourceCode != "" {
			parsed, err := parseSourceCodeString(trimSourceCode)
			if err != nil {
				return nil, err
			}

			// get settings
			v, _ := json.Marshal(parsed.Settings)
			files = append(files, ContractFile{
				Name:    "settings.json",
				Content: string(v),
			})

			// get remappings
			if strings.Contains(string(v), "remappings") {
				files = append(files, ContractFile{
					Name:    "remappings.txt",
					Content: parsed.Settings["remappings"].(string),
				})
			}

			// get contracts
			for _, key := range reflect.ValueOf(parsed.Sources).MapKeys() {
				files = append(files, ContractFile{
					Name:    key.String(),
					Content: parsed.Sources[key.String()]["content"],
				})
			}
		} else {
			ext := contractLanguage(result.CompilerVersion)
			files = append(files, ContractFile{
				Name:    result.ContractName + ext,
				Content: result.SourceCode,
			})
		}

		// If the contract is a proxy, we need to get the implementation contract
		if result.Implementation != "" && proxyDepth > 0 && result.Implementation != address {
			implCode, err := e.GetSourceCode(result.Implementation, proxyDepth-1)
			if err != nil {
				return nil, err
			}
			// return the implementation instead of
			return implCode, nil
		}
	}

	return files, nil
}
