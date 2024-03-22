package explorer

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	global "github.com/lowk3v/dumpsc/config"
	"reflect"
	"strings"
)

type Explorer struct {
	Chain            string `yaml:"chain"`
	ApiGetSourceCode string `yaml:"apiGetSourceCode"`
	ApiKey           string `yaml:"apiKey"`
}

type AddressInfo struct {
	Address  string `json:"address"`
	Contract string `json:"contract"`
	Chain    string `json:"chain"`
}

func New(explorerName string) (*Explorer, error) {
	val := reflect.Indirect(reflect.ValueOf(global.Config))

	c, is := val.Type().FieldByNameFunc(func(s string) bool {
		if strings.ToLower(s) == strings.ToLower(explorerName) {
			return true
		}
		return false
	})
	if !is {
		return &Explorer{}, errors.New("explorer is not supported")
	}

	e := val.FieldByName(c.Name).Interface().(global.ExplorerConfig)
	return &Explorer{
		Chain:            explorerName,
		ApiGetSourceCode: e.ApiGetSourceCode,
		ApiKey:           e.ApiKey,
	}, nil
}

func _isDataError(apiResponse *resty.Response, response *ApiResponse) error {
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

func (e *Explorer) GetSourceCode(address string, proxyDepth int) ([]ContractFile, error) {
	address = ChecksumAddress(address)
	if address == "" || address == "0x" {
		return nil, errors.New("address is empty")
	}

	fullApi := strings.ReplaceAll(e.ApiGetSourceCode, "{address}", address)
	fullApi = strings.ReplaceAll(fullApi, "{apiKey}", e.ApiKey)
	if proxyDepth > 3 {
		// Should not be more than 3
		proxyDepth = 3
	}

	var response ApiResponse
	httpClient := resty.New()
	rawResp, err := httpClient.R().SetResult(&response).Get(fullApi)
	if err != nil {
		return nil, err
	}
	if err := _isDataError(rawResp, &response); err != nil {
		return nil, err
	}

	var files []ContractFile
	var addresses []AddressInfo

	var contractInfos []ContractInfo
	marshal, _ := json.Marshal(response.Results)
	err = json.Unmarshal(marshal, &contractInfos)
	if err != nil {
		return nil, err
	}

	// classify source code
	for _, result := range contractInfos {
		if result.SourceCode == "" ||
			(result.ContractName == "" && result.ABI == "Contract source code not verified") {
			return nil, errors.New("the contract source code is not verified")
		}

		addresses = append(addresses, AddressInfo{
			Address:  address,
			Contract: result.ContractName,
			Chain:    e.Chain,
		})

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
				remapping, _ := json.Marshal(parsed.Settings["remappings"])
				files = append(files, ContractFile{
					Name:    "remappings.txt",
					Content: string(remapping),
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
			global.Log.Infof("Is proxy, get implementation: %s", result.Implementation)
			implCode, err := e.GetSourceCode(result.Implementation, proxyDepth-1)
			if err != nil {
				return nil, err
			}
			// return the implementation instead of
			return implCode, nil
		}
	}

	addressJson, _ := json.Marshal(addresses)
	files = append(files, ContractFile{
		Name:    "addresses.json",
		Content: string(addressJson),
	})

	return files, nil
}
