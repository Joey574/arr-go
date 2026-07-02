package handlers

import (
	"arr-go/v2/internal/fsu"
	"arr-go/v2/internal/log"
	"strings"
)

func (h *Handler) radarrImportDownloadHandler(conf radarrImportDownload) error {
	src := conf.MovieFileSourcePath
	dst := conf.MovieFilePath
	hash := strings.ToLower(conf.DownloadID)

	log.Infof("processing radarr import/download event: src='%s', dst='%s', id='%s'", src, dst, hash)

	if err := fsu.Symlink(src, dst); err != nil {
		return log.AsError("failed to create symlink: %v", err)
	}

	if err := h.qclient.Recheck(hash); err != nil {
		log.Errorf("failed to recheck torrent: %v", err)
	}

	if err := h.qclient.AddTags(hash, []string{"movie"}); err != nil {
		log.Errorf("failed to add tags: %v", err)
	}

	return nil
}

func (h *Handler) sonarrImportDownloadHandler(conf sonarrImportDownload) error {
	src := conf.EpisodeFileSourcePath
	dst := conf.EpisodeFilePath
	hash := strings.ToLower(conf.DownloadID)

	log.Infof("processing sonarr import/download event: src='%s', dst='%s', id='%s'", src, dst, hash)

	if err := fsu.Symlink(src, dst); err != nil {
		return log.AsError("failed to create symlink: %v", err)
	}

	if err := h.qclient.Recheck(hash); err != nil {
		log.Errorf("failed to recheck torrent: %v", err)
	}

	if err := h.qclient.AddTags(hash, []string{"movie"}); err != nil {
		log.Errorf("failed to add tags: %v", err)
	}

	return nil
}

func (h *Handler) lidarrImportDownloadHandler(conf lidarrImportDownload) error { return nil }
