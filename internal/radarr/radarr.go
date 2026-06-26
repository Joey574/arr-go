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

func HandleEvent() error {
	eventType, ok := os.LookupEnv(_RADARR_EVENT_TYPE)
	if !ok {
		return log.AsError("env variable not defined '%s'", _RADARR_EVENT_TYPE)
	}

	switch eventType {
	case _RADARR_TEST_EVENT:
		log.Infof("test event recieved, exiting")
		return nil
	case _RADARR_DOWNLOAD_EVENT:
		log.Infof("download event recieved, calling handler")
		return handleDownload()
	default:
		return log.AsError("unknown event recieved '%s', exiting", eventType)
	}
}

func handleDownload() error {
	src, ok := os.LookupEnv(_RADARR_SRC_PATH)
	if !ok {
		return log.AsError("env variable not defined '%s'", _RADARR_SRC_PATH)
	}

	dst, ok := os.LookupEnv(_RADARR_DST_PATH)
	if !ok {
		return log.AsError("env variable not defined '%s'", _RADARR_DST_PATH)
	}

	hash, ok := os.LookupEnv(_RADARR_DWN_ID)
	if !ok {
		return log.AsError("env variable not defined '%s'", _RADARR_DWN_ID)
	}
	hash = strings.ToLower(hash)

	log.Infof("recieved args: src='%s', dst='%s', id='%s'", src, dst, hash)

	if err := fsu.Symlink(src, dst); err != nil {
		return log.AsError("failed to create symlink: %v", err)
	}

	sid, err := qbit.Login()
	if err != nil {
		log.Errorf("failed to login to qbittorrent: %v", err)
		return nil // symlink completed succesfully so we exit eithout error
	}

	if err = qbit.Recheck(sid, hash); err != nil {
		log.Errorf("failed to recheck torrent: %v", err)
	}

	if err = qbit.AddTag(sid, hash, []string{"movie"}); err != nil {
		log.Errorf("failed to add tag: %v", err)
	}

	return nil
}
