package cli

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/jessevdk/go-flags"
)

type Args struct {
	Log       string `long:"log" default:"/var/log/arr/arr-go.log" description:"specify the log path, an empty path will disable logging"`
	Env       string `long:"env" default:"/etc/arr/arr.env" description:"specify the env file to source for creds"`
	Host      string `long:"host" default:"http://localhost:8080" description:"qbit host to connect to"`
	NoCleanup bool   `long:"no-cleanup" description:"skip clean up phase for media"`
	Version   bool   `short:"v" long:"version" description:"print version and exit"`
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
