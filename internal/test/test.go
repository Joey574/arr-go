package test

import (
	"arr-go/v2/internal/cli"
	"arr-go/v2/internal/log"
	"arr-go/v2/internal/qbit"
)

func TestQbitLogin(args *cli.Args) {
	err := qbit.SourceEnv(args.Env)
	if err != nil {
		log.Errorf("failed to source env: %v", err)
	}

	sid, err := qbit.Login()
	if err != nil {
		log.Errorf("failed to get sid: %v", err)
	}

	log.Infof("got sid='%s'", sid)
	version, err := qbit.Version(sid)
	if err != nil {
		log.Errorf("failed to get version: %v", err)
	}

	log.Infof("got qbit version: %s", version)
}
