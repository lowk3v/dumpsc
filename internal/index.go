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

// _parseUrl parse an url
// parsing url if url is valid download
// else print error
func _parseUrl(url string) (string, string, error) {
	if url == "" {
		return "", "", errors.New("the url is empty")
	}
	regExplorer := regexp.MustCompile(`https://([a-z.\-]+)\.`)
	regAddress := regexp.MustCompile(`https://.*/address/(0x[a-zA-Z0-9]{40}).*`)

	name := regExplorer.FindStringSubmatch(url)
	address := regAddress.FindStringSubmatch(url)

	if len(name) < 2 || len(address) < 2 {
		return "", "", errors.New("the url is not valid")
	}

	return name[1], address[1], nil
}

func (o Options) Run() {
	var expl *explorer.Explorer
	var err error

	// parse url
	if o.Url != "" {
		o.Explorer, o.Address, err = _parseUrl(o.Url)
		if utils.HandleError(err, "") {
			return
		}
	}

	expl, err = explorer.New(o.Explorer)
	if utils.HandleError(err, "") {
		return
	}

	// if empty, use default api key
	if o.ApiKey != "" {
		expl.ApiKey = o.ApiKey
	}

	// use default if empty
	utils.DirExists(o.Output, true)

	// download data
	fileContents, err := expl.GetSourceCode(o.Address, 3)
	if utils.HandleError(err, "") || fileContents == nil {
		return
	}

	// store data
	for _, fileContent := range fileContents {
		global.Log.Infof("downloaded: %s", fileContent.Name)
		if err := utils.WriteFile(
			fmt.Sprintf("%s/%s", o.Output, fileContent.Name),
			fileContent.Content,
		); err != nil {
			global.Log.Errorf("Error write file: %s", err.Error())
			continue
		}
	}
}

func (o Options) RunTest() (string, string, error) {
	if o.Url != "" {
		return _parseUrl(o.Url)
	}
	return "", "", nil
}
