package internal

import (
	"errors"
	"fmt"
	"github.com/fatih/structs"
	global "github.com/lowk3v/dumpsc/config"
	"github.com/lowk3v/dumpsc/internal/explorer"
	"github.com/lowk3v/dumpsc/utils"
	"regexp"
	"strings"
)

type Options struct {
	Action   ACTION
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
	regAddress := regexp.MustCompile(`https://.*/(address|token)/(0x[a-zA-Z0-9]{40}).*`)

	name := regExplorer.FindStringSubmatch(url)
	address := regAddress.FindStringSubmatch(url)

	if len(name) < 2 || len(address) < 2 {
		return "", "", errors.New("the url is not valid")
	}

	return name[1], address[2], nil
}

func _getSourceCode(o Options) {
	global.Log.Infof("output to %s", o.Output)
	expl, err := explorer.New(o.Explorer)
	if utils.HandleError(err, "") {
		return
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

func _listExplorer() {
	explorers := structs.Names(&global.Config)
	global.Log.Infof("List explorer: %s", strings.ToLower(strings.Join(explorers, ", ")))
}

func (o Options) Run() {
	var expl *explorer.Explorer
	var err error

	// if empty, use default api key
	if o.ApiKey != "" {
		expl.ApiKey = o.ApiKey
	}

	switch o.Action {
	case GETSOURCECODEBYURL:
		// parse url
		if o.Url != "" {
			o.Explorer, o.Address, err = _parseUrl(o.Url)
			if utils.HandleError(err, "") {
				return
			}
		}
		_getSourceCode(o)
	case GETSOURCECODE:
		_getSourceCode(o)
		break
	case LISTEXPLORER:
		_listExplorer()
		break
	case SHOWVERSION:
		global.Log.Infof("Version: %s", global.Version)
		break
	default:
		global.Log.Errorf("Action not found")
	}
}

func (o Options) RunTest() (string, string, error) {
	if o.Url != "" {
		return _parseUrl(o.Url)
	}
	return "", "", nil
}
