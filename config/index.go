package config

import (
	_ "embed"
	"github.com/lowk3v/dumpsc/pkg/log"
	"gopkg.in/yaml.v3"
	"os"
)

var Config Yaml
var Log log.Logger

//go:embed config.yml
var configYml string

type Yaml struct {
	EtherScan   ExplorerConfig `yaml:"etherscan"`
	BscScan     ExplorerConfig `yaml:"bscscan"`
	ArbiScan    ExplorerConfig `yaml:"arbiscan"`
	PolygonScan ExplorerConfig `yaml:"polygonscan"`
	CronoScan   ExplorerConfig `yaml:"cronoscan"`
	MoonScan    ExplorerConfig `yaml:"moonscan"`
	BaseScan    ExplorerConfig `yaml:"basescan"`
	SnowTrace   ExplorerConfig `yaml:"snowtrace"`
}

type ExplorerConfig struct {
	ApiGetSourceCode string `yaml:"apiGetSourceCode"`
	ApiKey           string `yaml:"apiKey"`
}

func init() {
	Log = *log.New("debug")

	// Load Config yml
	err := yaml.Unmarshal([]byte(configYml), &Config)
	if err != nil {
		Log.Errorf("Error loading config: %s", err)
	}
}

func CustomConfig(cfgPath string) error {
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
