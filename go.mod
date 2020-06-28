module gitlab.com/egeback/download_media_api

go 1.14

replace (
	github.com/egeback/play_media_api/internal/controllers => ./internal/controllers
	github.com/egeback/play_media_api/internal/docs => ./internal/docs
	github.com/egeback/play_media_api/internal/models => ./internal/models
	github.com/egeback/play_media_api/internal/utils => ./internal/utils
)
