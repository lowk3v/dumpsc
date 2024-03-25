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
	Address        string `json:"address"`
	Contract       string `json:"contract"`
	Implementation string `json:"implementation"`
	Chain          string `json:"chain"`
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

func (e *Explorer) pullSourceCodeRaw(address string) ([]ContractInfo, error) {
	fullApi := strings.ReplaceAll(e.ApiGetSourceCode, "{address}", address)
	fullApi = strings.ReplaceAll(fullApi, "{apiKey}", e.ApiKey)

	var response ApiResponse
	httpClient := resty.New()
	rawResp, err := httpClient.R().SetResult(&response).Get(fullApi)
	if err != nil {
		return nil, err
	}
	if err := isDataError(rawResp, &response); err != nil {
		return nil, err
	}

	var contractInfos []ContractInfo
	marshal, _ := json.Marshal(response.Results)
	err = json.Unmarshal(marshal, &contractInfos)
	if err != nil {
		return nil, err
	}
	return contractInfos, nil
}

func (e *Explorer) GetSourceCode(address string, proxyDepth int) ([]AddressInfo, []ContractFile, error) {
	addresses := make([]AddressInfo, 0)
	files := make([]ContractFile, 0)

	// pre-check address
	address = ChecksumAddress(address)
	if address == "" || address == "0x" {
		return nil, nil, errors.New("address is empty")
	}

	if proxyDepth > 3 {
		// Should not be more than 3
		proxyDepth = 3
	}

	// pulling source code
	contractInfos, err := e.pullSourceCodeRaw(address)
	if err != nil {
		return nil, nil, err
	}

	// classify source code
	for _, result := range contractInfos {
		// source code not found
		if result.SourceCode == "" ||
			(result.ContractName == "" && result.ABI == "Contract source code not verified") {
			return nil, nil, errors.New("the contract source code is not verified")
		}

		addressInfo := AddressInfo{
			Address:  address,
			Contract: result.ContractName,
			Chain:    e.Chain,
		}

		// source code has settings
		trimSourceCode := parseSourceCodeRaw(result.SourceCode)

		if trimSourceCode == "" {
			// other source codes
			ext := contractLanguage(result.CompilerVersion)
			files = append(files, ContractFile{
				Name:    result.ContractName + ext,
				Content: result.SourceCode,
			})
		} else {
			parsed, err := parseSourceCodeToStruct(trimSourceCode)
			if err != nil {
				return nil, nil, err
			}

			// get settings
			v, _ := json.Marshal(parsed.Settings)
			files = append(files, ContractFile{
				Name:    "settings.json",
				Content: string(v),
			})

			// get remapping
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
		}

		// If the contract is a proxy, we need to get the implementation contract
		if result.Implementation != "" && proxyDepth > 0 && result.Implementation != address {
			global.Log.Infof("Is proxy, get implementation: %s", result.Implementation)
			implAddrs, implCode, err := e.GetSourceCode(result.Implementation, proxyDepth-1)
			if err != nil {
				return nil, nil, err
			}

			// add the implementation to the addresses
			for i := range implAddrs {
				implAddrs[i].Address = address
				implAddrs[i].Implementation = result.Implementation
			}
			addresses = append(addresses, implAddrs...)

			// return the implementation instead of
			return addresses, implCode, nil
		}

		addresses = append(addresses, addressInfo)
	}

	return addresses, files, nil
}
