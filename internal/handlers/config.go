package handlers

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type eventType int

const (
	_UNKNOWN_EVENT_TYPE = eventType(iota)
	_RADARR_IMPORT_DOWNLOAD
	_SONARR_IMPORT_DOWNLOAD
	_LIDARR_IMPORT_DOWNLOAD
)

func getEventType() eventType {
	if event, ok := os.LookupEnv("radarr_eventtype"); ok && event == "Download" {
		return _RADARR_IMPORT_DOWNLOAD
	}

	if event, ok := os.LookupEnv("sonarr_eventtype"); ok && event == "Download" {
		return _SONARR_IMPORT_DOWNLOAD
	}

	if event, ok := os.LookupEnv("Lidarr_EventType"); ok && event == "AlbumDownload" {
		return _LIDARR_IMPORT_DOWNLOAD
	}

	return _UNKNOWN_EVENT_TYPE
}

func parseEventConfig(dest any) error {
	v := reflect.ValueOf(dest)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dest must be a pointer to a struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		if !fieldValue.CanSet() {
			continue
		}

		envKey := field.Tag.Get("env")
		if envKey == "" || envKey == "-" {
			continue
		}

		val, exists := os.LookupEnv(envKey)
		if !exists {
			val = field.Tag.Get("default")
		}

		if val == "" {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			parsed, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid integer '%s' for field '%s'", val, field.Name)
			}
			fieldValue.SetInt(parsed)
		case reflect.Bool:
			parsed, err := strconv.ParseBool(val)
			if err != nil {
				return fmt.Errorf("invalid boolean '%s' for field '%s'", val, field.Name)
			}
			fieldValue.SetBool(parsed)
		default:
			return fmt.Errorf("unsupported type '%s' for field '%s'", fieldValue.Kind(), field.Name)
		}
	}

	return nil
}

type radarrImportDownload struct {
	EventType                           string `env:"radarr_eventtype"`
	InstanceName                        string `env:"radarr_instancename"`
	ApplicationURL                      string `env:"radarr_applicationurl"`
	IsUpgrade                           bool   `env:"radarr_isupgrade"`
	MovieID                             string `env:"radarr_movie_id"`
	MovieTitle                          string `env:"radarr_movie_title"`
	MovieYear                           int    `env:"radarr_movie_year"`
	MovieOriginalLanguage               string `env:"radarr_movie_originallanguage"`
	MovieGenres                         string `env:"radarr_movie_genres"`
	MovieTags                           string `env:"radarr_movie_tags"`
	MoviePath                           string `env:"radarr_movie_path"`
	MovieIMDbID                         string `env:"radarr_movie_imdbid"`
	MovieTMDbID                         string `env:"radarr_movie_tmdbid"`
	MovieCinemaReleaseDate              string `env:"radarr_movie_in_cinemas_date"`
	MoviePhysicalReleaseDate            string `env:"radarr_movie_physical_release_date"`
	MovieOverview                       string `env:"radarr_movie_overview"`
	MovieFileID                         string `env:"radarr_moviefile_id"`
	MovieFileRelativePath               string `env:"radarr_moviefile_relativepath"`
	MovieFilePath                       string `env:"radarr_moviefile_path"`
	MovieFileQuality                    string `env:"radarr_moviefile_quality"`
	MovieFileQualityVersion             int    `env:"radarr_moviefile_qualityversion"`
	MovieFileReleaseGroup               string `env:"radarr_moviefile_releasegroup"`
	MovieFileSceneName                  string `env:"radarr_moviefile_scenename"`
	MovieFileSourcePath                 string `env:"radarr_moviefile_sourcepath"`
	MovieFileSourceFolder               string `env:"radarr_moviefile_sourcefolder"`
	MovieFileMediaInfoAudioChannels     string `env:"radarr_moviefile_mediainfo_audiochannels"`
	MovieFileMediaInfoAudioCodec        string `env:"radarr_moviefile_mediainfo_audiocodec"`
	MovieFileMediaInfoAudioLanguages    string `env:"radarr_moviefile_mediainfo_audiolanguages"`
	MovieFileMediaInfoLanguages         string `env:"radarr_moviefile_mediainfo_languages"`
	MovieFileMediaInfoHeight            int    `env:"radarr_moviefile_mediainfo_height"`
	MovieFileMediaInfoWidth             int    `env:"radarr_moviefile_mediainfo_width"`
	MovieFileMediaInfoSubtitles         string `env:"radarr_moviefile_mediainfo_subtitles"`
	MovieFileMediaInfoVideoCodec        string `env:"radarr_moviefile_mediainfo_videocodec"`
	MovieFileMediaInfoVideoDynamicRange string `env:"radarr_moviefile_mediainfo_videodynamicrangetype"`
	MovieFileCustomFormat               string `env:"radarr_moviefile_customformat"`
	MovieFileCustomFormatScore          int    `env:"radarr_moviefile_customformatscore"`
	ReleaseIndexer                      string `env:"radarr_release_indexer"`
	ReleaseSize                         int64  `env:"radarr_release_size"`
	ReleaseTitle                        string `env:"radarr_release_title"`
	DownloadClient                      string `env:"radarr_download_client"`
	DownloadClientType                  string `env:"radarr_download_client_type"`
	DownloadID                          string `env:"radarr_download_id"`
	DeletedRelativePaths                string `env:"radarr_deletedrelativepaths"`
	DeletedPaths                        string `env:"radarr_deletedpaths"`
	DeletedDateAdded                    string `env:"radarr_deleteddateadded"`
	DeletedRecycleBinPaths              string `env:"radarr_deletedrecyclebinpaths"`
}

type sonarrImportDownload struct {
	EventType                          string `env:"sonarr_eventtype"`
	InstanceName                       string `env:"sonarr_instancename"`
	ApplicationURL                     string `env:"sonarr_applicationurl"`
	IsUpgrade                          bool   `env:"sonarr_isupgrade"`
	SeriesID                           string `env:"sonarr_series_id"`
	SeriesTitle                        string `env:"sonarr_series_title"`
	SeriesTitleSlug                    string `env:"sonarr_series_titleslug"`
	SeriesPath                         string `env:"sonarr_series_path"`
	SeriesTVDBID                       string `env:"sonarr_series_tvdbid"`
	SeriesTVMazeID                     string `env:"sonarr_series_tvmazeid"`
	SeriesTMDBID                       string `env:"sonarr_series_tmdbid"`
	SeriesIMDBID                       string `env:"sonarr_series_imdbid"`
	SeriesType                         string `env:"sonarr_series_type"`
	SeriesYear                         int    `env:"sonarr_series_year"`
	SeriesOriginalLanguage             string `env:"sonarr_series_originallanguage"`
	SeriesGenres                       string `env:"sonarr_series_genres"`
	SeriesTags                         string `env:"sonarr_series_tags"`
	EpisodeFileID                      string `env:"sonarr_episodefile_id"`
	EpisodeFileEpisodeCount            int    `env:"sonarr_episodefile_episodecount"`
	EpisodeFileRelativePath            string `env:"sonarr_episodefile_relativepath"`
	EpisodeFilePath                    string `env:"sonarr_episodefile_path"`
	EpisodeFileEpisodeIDs              string `env:"sonarr_episodefile_episodeids"`
	EpisodeFileSeasonNumber            int    `env:"sonarr_episodefile_seasonnumber"`
	EpisodeFileEpisodeNumbers          string `env:"sonarr_episodefile_episodenumbers"`
	EpisodeFileEpisodeAirDates         string `env:"sonarr_episodefile_episodeairdates"`
	EpisodeFileEpisodeAirDatesUTC      string `env:"sonarr_episodefile_episodeairdatesutc"`
	EpisodeFileEpisodeTitles           string `env:"sonarr_episodefile_episodetitles"`
	EpisodeFileEpisodeOverviews        string `env:"sonarr_episodefile_episodeoverviews"`
	EpisodeFileQuality                 string `env:"sonarr_episodefile_quality"`
	EpisodeFileQualityVersion          int    `env:"sonarr_episodefile_qualityversion"`
	EpisodeFileReleaseGroup            string `env:"sonarr_episodefile_releasegroup"`
	EpisodeFileSceneName               string `env:"sonarr_episodefile_scenename"`
	EpisodeFileSourcePath              string `env:"sonarr_episodefile_sourcepath"`
	EpisodeFileSourceFolder            string `env:"sonarr_episodefile_sourcefolder"`
	EpisodeFileMediaInfoAudioChannels  string `env:"sonarr_episodefile_mediainfo_audiochannels"`
	EpisodeFileMediaInfoAudioCodec     string `env:"sonarr_episodefile_mediainfo_audiocodec"`
	EpisodeFileMediaInfoAudioLanguages string `env:"sonarr_episodefile_mediainfo_audiolanguages"`
	EpisodeFileMediaInfoLanguages      string `env:"sonarr_episodefile_mediainfo_languages"`
	EpisodeFileMediaInfoHeight         int    `env:"sonarr_episodefile_mediainfo_height"`
	EpisodeFileMediaInfoWidth          int    `env:"sonarr_episodefile_mediainfo_width"`
	EpisodeFileMediaInfoSubtitles      string `env:"sonarr_episodefile_mediainfo_subtitles"`
	EpisodeFileMediaInfoVideoCodec     string `env:"sonarr_episodefile_mediainfo_videocodec"`
	EpisodeFileMediaInfoVideoDynamicR  string `env:"sonarr_episodefile_mediainfo_videodynamicrangetype"`
	EpisodeFileCustomFormat            string `env:"sonarr_episodefile_customformat"`
	EpisodeFileCustomFormatScore       int    `env:"sonarr_episodefile_customformatscore"`
	DownloadClient                     string `env:"sonarr_download_client"`
	DownloadClientType                 string `env:"sonarr_download_client_type"`
	DownloadID                         string `env:"sonarr_download_id"`
	ReleaseIndexer                     string `env:"sonarr_release_indexer"`
	ReleaseSize                        int64  `env:"sonarr_release_size"`
	ReleaseTitle                       string `env:"sonarr_release_title"`
	ReleaseType                        string `env:"sonarr_release_releasetype"`
	DeletedRelativePaths               string `env:"sonarr_deletedrelativepaths"`
	DeletedPaths                       string `env:"sonarr_deletedpaths"`
	DeletedDateAdded                   string `env:"sonarr_deleteddateadded"`
	DeletedRecycleBinPaths             string `env:"sonarr_deletedrecyclebinpaths"`
}

type lidarrImportDownload struct {
	EventType          string `env:"Lidarr_EventType"`
	ArtistID           string `env:"Lidarr_Artist_Id"`
	ArtistName         string `env:"Lidarr_Artist_Name"`
	ArtistPath         string `env:"Lidarr_Artist_Path"`
	ArtistMBID         string `env:"Lidarr_Artist_MBId"`
	ArtistType         string `env:"Lidarr_Artist_Type"`
	ArtistGenres       string `env:"Lidarr_Artist_Genres"`
	ArtistTags         string `env:"Lidarr_Artist_Tags"`
	AlbumID            string `env:"Lidarr_Album_Id"`
	AlbumTitle         string `env:"Lidarr_Album_Title"`
	AlbumOverview      string `env:"Lidarr_Album_Overview"`
	AlbumMBID          string `env:"Lidarr_Album_MBId"`
	AlbumReleaseMBID   string `env:"Lidarr_AlbumRelease_MBId"`
	AlbumReleaseDate   string `env:"Lidarr_Album_ReleaseDate"`
	DownloadClient     string `env:"Lidarr_Download_Client"`
	DownloadClientType string `env:"Lidarr_Download_Client_Type"`
	DownloadID         string `env:"Lidarr_Download_Id"`
	AddedTrackPaths    string `env:"Lidarr_AddedTrackPaths"`
	DeletedPaths       string `env:"Lidarr_DeletedPaths"`
	DeletedDateAdded   string `env:"Lidarr_DeletedDateAdded"`
}
