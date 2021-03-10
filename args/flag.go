package args

import (
	"flag"
	"os"
)

var (
	Help       bool
	ConfigPath string
)

func init() {
	flag.BoolVar(&Help, "h", false, "Show Help")
	flag.StringVar(&ConfigPath, "c", os.Args[0]+".yaml", "Configure path")
}

func Parse()  {
	flag.Parse()
	if Help {
		flag.Usage()
		os.Exit(0)
	}
}