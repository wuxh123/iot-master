package args

import (
	"flag"
	"fmt"
	"os"
)

var (
	showHelp    bool
	showVersion bool
	ConfigPath  string
)

var (
	gitHash   string
	buildTime string
	goVersion string
)

func init() {
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.StringVar(&ConfigPath, "c", os.Args[0]+".yaml", "Configure path")
}

func Parse() {
	flag.Parse()
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}
	if showVersion {
		fmt.Printf("Git Commit Hash: %s \n", gitHash)
		fmt.Printf("Build TimeStamp: %s \n", buildTime)
		fmt.Printf("GoLang Version: %s \n", goVersion)
		os.Exit(0)
	}
}
