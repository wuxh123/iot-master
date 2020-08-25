package flag

import "flag"

var (
	help bool
	SyncTables bool
	ConfigPath string
)

func init() {
	flag.BoolVar(&help, "h", false, "Show help")
	flag.BoolVar(&SyncTables, "s", false, "Sync tables")
	flag.StringVar(&ConfigPath, "c", "config.yaml", "Configure path")
}


func Parse() bool {
	flag.Parse()
	if help {
		flag.Usage()
		return false
	}
	return true
}
