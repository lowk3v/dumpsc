package explorer

type ApiResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Results []ContractInfo `json:"result"`
}

type ContractInfo struct {
	SourceCode           string `json:"SourceCode"`
	ABI                  string `json:"ABI"`
	ContractName         string `json:"ContractName"`
	CompilerVersion      string `json:"CompilerVersion"`
	OptimizationUsed     string `json:"OptimizationUsed"`
	Runs                 string `json:"Runs"`
	ConstructorArguments string `json:"ConstructorArguments"`
	EVMVersion           string `json:"EVMVersion"`
	Library              string `json:"Library"`
	LicenseType          string `json:"LicenseType"`
	Proxy                string `json:"Proxy"`
	Implementation       string `json:"Implementation"`
	SwarmSource          string `json:"SwarmSource"`
}

type SourceCodeInfo struct {
	Language string                       `json:"language"`
	Sources  map[string]map[string]string `json:"sources"`
	Settings map[string]interface{}       `json:"settings"`
}

type ContractFile struct {
	Name    string
	Content string
}
