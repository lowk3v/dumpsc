package main

import (
	"flag"
	"fmt"
	"github.com/lowk3v/dumpsc/internal"
	"os"
)
import global "github.com/lowk3v/dumpsc/config"

func _validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func parseFlags() (string, *internal.Options, error) {
	var configPath string
	var explorer string
	var address string
	var apikey string
	var output string
	var url string

	flag.StringVar(&configPath, "c", "./config.yml", "Optional. Path to config file")
	flag.StringVar(&explorer, "e", "etherscan", "Required. An explorer to use")
	flag.StringVar(&address, "a", "", "Required. Smart contract address to query")
	flag.StringVar(&apikey, "k", "", "Optional. api key of an explorer to use")
	flag.StringVar(&output, "o", "src", "Optional. Output directory")
	flag.StringVar(&url, "u", "", "Optional. Url to download")
	flag.Parse()

	// require options
	if explorer == "" || address == "" {
		if url == "" {
			flag.Usage()
			os.Exit(0)
		}
	}

	options := &internal.Options{
		Explorer: explorer,
		Address:  address,
		ApiKey:   apikey,
		Output:   output,
		Url:      url,
	}

	if err := _validateConfigPath(configPath); err != nil {
		return "", &internal.Options{}, err
	}

	// Return the configuration path
	return configPath, options, nil
}

func main() {
	cfgPath, options, err := parseFlags()
	if err != nil {
		os.Exit(0)
	}
	if global.NewConfig(cfgPath) != nil {
		os.Exit(0)
	}
	app := internal.New(options)
	app.Run()
}