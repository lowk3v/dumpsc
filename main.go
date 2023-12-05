package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
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

func _banner() {
	// https://patorjk.com/software/taag/#p=display&f=ANSI%20Shadow&t=%20dumpsc
	_, _ = fmt.Fprintf(os.Stderr, "%s by %s\n%s\nCredits: https://github.com/lowk3v/%s\nTwitter: https://twitter.com/%s\n\n",
		color.HiBlueString(`
    ██████╗ ██╗   ██╗███╗   ███╗██████╗ ███████╗ ██████╗
    ██╔══██╗██║   ██║████╗ ████║██╔══██╗██╔════╝██╔════╝
    ██║  ██║██║   ██║██╔████╔██║██████╔╝███████╗██║     
    ██║  ██║██║   ██║██║╚██╔╝██║██╔═══╝ ╚════██║██║     
    ██████╔╝╚██████╔╝██║ ╚═╝ ██║██║     ███████║╚██████╗
    ╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝     ╚══════╝ ╚═════╝`),
		color.BlueString("@LowK"),
		"A tool is used to download a verified source code of smart contracts from an explorer.",
		"dumpsc",
		"Lowk3v_",
	)
	_, _ = fmt.Fprintf(os.Stderr, "Usage of: %s <options> <args>\n", os.Args[0])
	flag.PrintDefaults()
}

func parseFlags() (string, *internal.Options, error) {
	action := internal.NONE
	var configPath string
	var explorer string
	var address string
	var apikey string
	var output string
	var url string
	var listExplorer bool

	flag.StringVar(&configPath, "c", "./config.yml", "Optional. Path to config file")
	flag.StringVar(&explorer, "e", "etherscan", "Required. An explorer to use")
	flag.StringVar(&address, "a", "", "Required. Smart contract address to query")
	flag.StringVar(&apikey, "k", "", "Optional. api key of an explorer to use")
	flag.StringVar(&output, "o", "src", "Optional. Output directory")
	flag.StringVar(&url, "u", "", "Optional. Url to download")
	flag.BoolVar(&listExplorer, "l", false, "Optional. Url to download")
	flag.Usage = _banner
	flag.Parse()

	if listExplorer {
		action = internal.LISTEXPLORER
	} else if url != "" {
		action = internal.GETSOURCECODEBYURL
	} else if explorer != "" && address != "" {
		action = internal.GETSOURCECODE
	} else {
		_banner()
		flag.PrintDefaults()
		os.Exit(0)
	}

	options := &internal.Options{
		Action:   action,
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
