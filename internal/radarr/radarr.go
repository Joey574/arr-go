package radarr

import (
	"arr-go/v2/internal/fsu"
	"arr-go/v2/internal/log"
	"arr-go/v2/internal/qbit"
	"os"
	"strings"
)

const (
	_RADARR_SRC_PATH   = "radarr_moviefile_sourcepath"
	_RADARR_DST_PATH   = "radarr_moviefile_path"
	_RADARR_DWN_ID     = "radarr_download_id"
	_RADARR_EVENT_TYPE = "radarr_eventtype"

	_RADARR_TEST_EVENT     = "Test"
	_RADARR_DOWNLOAD_EVENT = "Download"
)

func IsRadarr() bool {
	_, ok := os.LookupEnv(_RADARR_EVENT_TYPE)
	return ok
}

func HandleEvent() {
	eventType, ok := os.LookupEnv(_RADARR_EVENT_TYPE)
	if !ok {
		log.Errorf("env variable not defined '%s'", _RADARR_EVENT_TYPE)
	}

	switch eventType {
	case _RADARR_TEST_EVENT:
		log.Infof("test event recieved, exiting")
		os.Exit(0)
	case _RADARR_DOWNLOAD_EVENT:
		log.Infof("download event recieved, calling handler")
		handleDownload()
	default:
		log.Fatalf("unknown event recieved '%s', exiting", eventType)
	}
}

func handleDownload() {
	src, ok := os.LookupEnv(_RADARR_SRC_PATH)
	if !ok {
		log.Fatalf("env variable not defined '%s'", _RADARR_SRC_PATH)
	}

	dst, ok := os.LookupEnv(_RADARR_DST_PATH)
	if !ok {
		log.Fatalf("env variable not defined '%s'", _RADARR_DST_PATH)
	}

	hash, ok := os.LookupEnv(_RADARR_DWN_ID)
	if !ok {
		log.Fatalf("env variable not defined '%s'", _RADARR_DWN_ID)
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
