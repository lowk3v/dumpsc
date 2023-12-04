package config

import (
	"errors"
	"github.com/lowk3v/dumpsc/internal/explorer"
	"github.com/lowk3v/dumpsc/pkg/log"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strings"
)

var Config Yaml
var Log log.Logger

type Yaml struct {
	EtherScan   explorer.Explorer `yaml:"etherscan"`
	BscScan     explorer.Explorer `yaml:"bscscan"`
	PolygonScan explorer.Explorer `yaml:"polygonscan"`
	FtmScan     explorer.Explorer `yaml:"ftmscan"`
	HecoInfo    explorer.Explorer `yaml:"hecoinfo"`
	SnowTrace   explorer.Explorer `yaml:"snowtrace"`
	ArbiScan    explorer.Explorer `yaml:"arbiscan"`
	AvaxScan    explorer.Explorer `yaml:"avaxscan"`
	CronoScan   explorer.Explorer `yaml:"cronoscan"`
	MoonBean    explorer.Explorer `yaml:"moonbean"`
	AuroraScan  explorer.Explorer `yaml:"aurorascan"`
	BaseScan    explorer.Explorer `yaml:"basescan"`
}

func NewConfig(cfgPath string) error {
	Log = *log.New("debug")
	Config = Yaml{}

	// Open config file
	file, err := os.Open(cfgPath)
	if err != nil {
		Log.Errorf("Error opening config file: %s", err)
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&Config); err != nil {
		return err
	}

	return nil
}

func (y *Yaml) GetExplorerConfig(name string) (*explorer.Explorer, error) {
	val := reflect.Indirect(reflect.ValueOf(y))
	c, is := val.Type().FieldByNameFunc(func(s string) bool {
		if strings.ToLower(s) == strings.ToLower(name) {
			return true
		}
		return false
	})
	if !is {
		return &explorer.Explorer{}, errors.New("explorer is not supported")
	}

	e := val.FieldByName(c.Name).Interface().(explorer.Explorer)
	return &e, nil
}
