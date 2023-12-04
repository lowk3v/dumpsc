package internal

import (
	"errors"
	"fmt"
	global "github.com/lowk3v/dumpsc/config"
	"github.com/lowk3v/dumpsc/internal/explorer"
	"github.com/lowk3v/dumpsc/utils"
	"regexp"
)

type Options struct {
	Explorer string
	Address  string
	ApiKey   string
	Output   string
	Url      string
}

func New(opt *Options) *Options {
	return opt
}

// _storeContent store content to file
func _storeContent(contents string, output string) error {
	return nil
}

// _parseUrl parse an url
// parsing url if url is valid download
// else print error
func _parseUrl(url string) (*explorer.Explorer, string, error) {
	var err error
	if url == "" {
		err = errors.New("url is empty")
	}
	regExplorer := regexp.MustCompile(`https://([a-z.\-]+)\.`)
	regAddress := regexp.MustCompile(`https://.*/address/([a-zA-Z0-9]{40}).*`)

	var expl *explorer.Explorer

	preExplorer := regExplorer.FindStringSubmatch(url)
	address := regAddress.FindStringSubmatch(url)

	if len(preExplorer) < 2 || len(address) < 2 {
		return &explorer.Explorer{}, "", errors.New("the url is not valid")
	}

	expl, err = global.Config.GetExplorerConfig(preExplorer[1])
	if err != nil {
		return &explorer.Explorer{}, "", err
	}

	global.Log.Infof("expl: %s", expl)
	global.Log.Infof("address: %s", address[1])
	return &explorer.Explorer{}, "", nil
}

func (o Options) Run() {
	var expl *explorer.Explorer
	var address string
	var err error

	if o.Url != "" {
		expl, address, err = _parseUrl(o.Url)
		if utils.HandleError(err, "") {
			return
		}
	} else {
		// checked explorer and address are not empty
		expl, err = global.Config.GetExplorerConfig(o.Explorer)
		if utils.HandleError(err, "") {
			return
		}
		address = o.Address
	}
	if o.ApiKey != "" {
		// if empty, use default api key
		expl.ApiKey = o.ApiKey
	}

	// use default if empty
	utils.DirExists(o.Output, true)

	// download data
	fileContents, err := expl.GetSourceCode(address, 3)
	if len(fileContents) == 0 {
		return
	}

	// store data
	for _, fileContent := range fileContents {
		global.Log.Infof("file: %s", fileContent.Name)
		if err := utils.WriteFile(
			fmt.Sprintf("%s/%s", o.Output, fileContent.Name),
			fileContent.Content,
		); err != nil {
			global.Log.Errorf("Error write file: %s", err.Error())
			continue
		}
	}
}

func (o Options) RunTest() (*explorer.Explorer, string, error) {
	if o.Url != "" {
		return _parseUrl(o.Url)
	}
	return &explorer.Explorer{}, "", nil
}
