package main

import (
	"arr-go/v2/internal/cli"
	"arr-go/v2/internal/handlers"
	"arr-go/v2/internal/log"
	"arr-go/v2/internal/qbit"

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

	if err = log.SetLogFile(arg.Log); err != nil {
		log.Errorf("opening log file '%s': %v", arg.Log, err)
	}

	qclient := qbit.NewClient(arg.Host)
	if err = qclient.SourceEnv(arg.Env); err != nil {
		log.Errorf("sourcing env file '%s': '%v'", arg.Env, err)
	}

	handler := handlers.NewHandler(qclient)
	if err = handler.HandleEvent(); err != nil {
		log.Fatalf("%v", err)
	}
}
