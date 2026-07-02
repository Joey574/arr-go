package handlers

import (
	"arr-go/v2/internal/log"
	"arr-go/v2/internal/qbit"
)

type Handler struct {
	qclient *qbit.Client
}

func NewHandler(qclient *qbit.Client) *Handler {
	return &Handler{
		qclient: qclient,
	}
}

func (h *Handler) HandleEvent() error {
	et := getEventType()

	err := h.qclient.Login()
	if err != nil {
		return err
	}

	switch et {
	case _RADARR_IMPORT_DOWNLOAD:
		var cfg radarrImportDownload
		if err := parseEventConfig(cfg); err != nil {
			return err
		}

		return h.radarrImportDownloadHandler(cfg)
	case _SONARR_IMPORT_DOWNLOAD:
		var cfg sonarrImportDownload
		if err := parseEventConfig(cfg); err != nil {
			return err
		}

		return h.sonarrImportDownloadHandler(cfg)
	case _LIDARR_IMPORT_DOWNLOAD:
		var cfg lidarrImportDownload
		if err := parseEventConfig(cfg); err != nil {
			return err
		}

		return h.lidarrImportDownloadHandler(cfg)
	default:
		return log.AsError("unrecognized event")
	}
}
