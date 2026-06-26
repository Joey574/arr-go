package cli

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/jessevdk/go-flags"
)

type Args struct {
	Log     string `long:"log" default:"/var/log/radarr/radarr-go.log" description:"specify the log path, an empty path will disable logging"`
	Env     string `long:"env" default:"/etc/radarr/radarr.env" description:"specify the env file to source for creds"`
	Version bool   `short:"v" long:"version" description:"print version and exit"`

	TestQbitLogin bool `long:"test-qbit"`
}

func NewArgs() *Args {
	return &Args{}
}

func (a *Args) Parse() ([]string, error) {
	args, err := flags.Parse(a)
	if err != nil {
		return nil, err
	}

	if a.Version {
		var version string
		info, ok := debug.ReadBuildInfo()
		if !ok {
			version = "unknown"
		} else {
			version = info.Main.Version
		}

		fmt.Println(version)
		os.Exit(0)
	}

	return args, nil
}
