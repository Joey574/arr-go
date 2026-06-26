package sonarr

import (
	"arr-go/v2/internal/fsu"
	"arr-go/v2/internal/log"
	"arr-go/v2/internal/qbit"
	"os"
	"strings"
)

const (
	_SONARR_SRC_PATH   = "sonarr_episodefile_sourcepath"
	_SONARR_DST_PATH   = "sonarr_episodefile_path"
	_SONARR_DWN_ID     = "sonarr_download_id"
	_SONARR_EVENT_TYPE = "sonarr_eventtype"

	_SONARR_TEST_EVENT     = "Test"
	_SONARR_DOWNLOAD_EVENT = "Download"
)

func IsSonarr() bool {
	_, ok := os.LookupEnv(_SONARR_EVENT_TYPE)
	return ok
}

func HandleEvent() {
	eventType, ok := os.LookupEnv(_SONARR_EVENT_TYPE)
	if !ok {
		log.Errorf("env variable not defined '%s'", _SONARR_EVENT_TYPE)
	}

	switch eventType {
	case _SONARR_TEST_EVENT:
		log.Infof("test event recieved, exiting")
		os.Exit(0)
	case _SONARR_DOWNLOAD_EVENT:
		log.Infof("download event recieved, calling handler")
		handleDownload()
	default:
		log.Fatalf("unknown event recieved '%s', exiting", eventType)
	}
}

func handleDownload() {
	src, ok := os.LookupEnv(_SONARR_SRC_PATH)
	if !ok {
		log.Fatalf("env variable not defined '%s'", _SONARR_SRC_PATH)
	}

	dst, ok := os.LookupEnv(_SONARR_DST_PATH)
	if !ok {
		log.Fatalf("env variable not defined '%s'", _SONARR_DST_PATH)
	}

	hash, ok := os.LookupEnv(_SONARR_DWN_ID)
	if !ok {
		log.Fatalf("env variable not defined '%s'", _SONARR_DWN_ID)
	}
	hash = strings.ToLower(hash)

	log.Infof("recieved args: src='%s', dst='%s', id='%s'", src, dst, hash)

	if err := fsu.Symlink(src, dst); err != nil {
		log.Fatalf("failed to create symlink: %v", err)
	}

	sid, err := qbit.Login()
	if err != nil {
		log.Errorf("failed to login to qbittorrent: %v", err)
	}

	if err = qbit.Recheck(sid, hash); err != nil {
		log.Errorf("failed to recheck torrent: %v", err)
	}
}
