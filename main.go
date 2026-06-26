package main

import (
	"arr-go/v2/internal/cli"
	"arr-go/v2/internal/log"
	"arr-go/v2/internal/qbit"
	"arr-go/v2/internal/radarr"
	"arr-go/v2/internal/sonarr"
	"arr-go/v2/internal/test"
	"os"
)

func main() {
	arg := cli.NewArgs()
	arg.Parse()

	if arg.TestQbitLogin {
		test.TestQbitLogin(arg)
		os.Exit(0)
	}

	if err := log.SetLogFile(arg.Log); err != nil {
		log.Errorf("opening log file '%s': %v", arg.Log, err)
	}

	if err := qbit.SourceEnv(arg.Env); err != nil {
		log.Errorf("sourcing env file '%s': '%v'", arg.Env, err)
	}

	if radarr.IsRadarr() {
		radarr.HandleEvent()
	} else if sonarr.IsSonarr() {
		sonarr.HandleEvent()
	} else {
		log.Errorf("unable to determine source")
	}
}
