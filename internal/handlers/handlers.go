package handlers

import (
	"arr-go/v2/internal/fsu"
	"arr-go/v2/internal/log"
	"arr-go/v2/internal/qbit"
	"os"
	"strings"
)

func isRadarr() bool {
	_, ok := os.LookupEnv(_RADARR_EVENT_TYPE)
	return ok
}

func isSonarr() bool {
	_, ok := os.LookupEnv(_SONARR_EVENT_TYPE)
	return ok
}

func isLidarr() bool {
	_, ok := os.LookupEnv(_LIDARR_EVENT_TYPE)
	return ok
}

func rsImportHandler(srcEnv, dstEnv, dwnIdEnv string, tags []string) error {
	src, ok := os.LookupEnv(srcEnv)
	if !ok {
		return log.AsError("env variable not defined '%s'", srcEnv)
	}

	dst, ok := os.LookupEnv(dstEnv)
	if !ok {
		return log.AsError("env variable not defined '%s'", dstEnv)
	}

	hash, ok := os.LookupEnv(dwnIdEnv)
	if !ok {
		return log.AsError("env variable not defined '%s'", dwnIdEnv)
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

	if err = qbit.AddTag(sid, hash, tags); err != nil {
		log.Errorf("failed to add tag: %v", err)
	}

	return nil
}

func albumDownloadHandler(srcEnv, dstEnv, dwnIdEnv string, tags []string) error {
	return nil
}

func HandleEvent() error {
	var radarr, sonarr, lidarr bool
	var event string

	if isRadarr() {
		event, _ = os.LookupEnv(_RADARR_EVENT_TYPE)
		radarr = true
	} else if isSonarr() {
		event, _ = os.LookupEnv(_SONARR_EVENT_TYPE)
		sonarr = true
	} else if isLidarr() {
		event, _ = os.LookupEnv(_LIDARR_EVENT_TYPE)
		lidarr = true
	} else {
		return log.AsError("cannot determine source")
	}

	switch event {
	case _DOWNLOAD_EVENT:
		if radarr {
			return rsImportHandler(_RADARR_SRC_PATH, _RADARR_DST_PATH, _RADARR_DWN_ID, []string{"movie"})
		} else if sonarr {
			return rsImportHandler(_SONARR_SRC_PATH, _SONARR_DST_PATH, _SONARR_DWN_ID, []string{"show"})
		}
	case _ALBUM_DOWNLOAD_EVENT:
		if lidarr {
			return albumDownloadHandler(_LIDARR_SRC_PATH, _LIDARR_DST_PATH, _LIDARR_DWN_ID, []string{"music"})
		}
	default:
		return log.AsError("unrecognized event: %s", event)
	}

	return nil
}
