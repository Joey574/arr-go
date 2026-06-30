package main

import (
	"arr-go/v2/internal/cli"
	"arr-go/v2/internal/handlers"
	"arr-go/v2/internal/log"
	"arr-go/v2/internal/qbit"
	"arr-go/v2/internal/test"
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	arg := cli.NewArgs()
	_, err := arg.Parse()
	if flags.WroteHelp(err) {
		return
	}

	if err != nil {
		log.Fatalf("%v", err)
	}

	if arg.TestQbitLogin {
		test.TestQbitLogin(arg)
		os.Exit(0)
	}

	if err = log.SetLogFile(arg.Log); err != nil {
		log.Errorf("opening log file '%s': %v", arg.Log, err)
	}

	if err = qbit.SourceEnv(arg.Env); err != nil {
		log.Errorf("sourcing env file '%s': '%v'", arg.Env, err)
	}

	if err = handlers.HandleEvent(); err != nil {
		log.Fatalf("%v", err)
	}
}
