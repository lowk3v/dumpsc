package config

import (
	"github.com/lowk3v/dumpsc/pkg/log"
	"gopkg.in/yaml.v3"
	"os"
)

var Config Yaml
var Log log.Logger

type Yaml struct {
	EtherScan   ExplorerConfig `yaml:"etherscan"`
	BscScan     ExplorerConfig `yaml:"bscscan"`
	PolygonScan ExplorerConfig `yaml:"polygonscan"`
	FtmScan     ExplorerConfig `yaml:"ftmscan"`
	HecoInfo    ExplorerConfig `yaml:"hecoinfo"`
	SnowTrace   ExplorerConfig `yaml:"snowtrace"`
	ArbiScan    ExplorerConfig `yaml:"arbiscan"`
	AvaxScan    ExplorerConfig `yaml:"avaxscan"`
	CronoScan   ExplorerConfig `yaml:"cronoscan"`
	MoonBean    ExplorerConfig `yaml:"moonbean"`
	AuroraScan  ExplorerConfig `yaml:"aurorascan"`
	BaseScan    ExplorerConfig `yaml:"basescan"`
}

type ExplorerConfig struct {
	ApiGetSourceCode string `yaml:"apiGetSourceCode"`
	ApiKey           string `yaml:"apiKey"`
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
