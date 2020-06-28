module github.com/egeback/download_media_api/internal

go 1.14

replace (
	github.com/egeback/download_media_api/internal/controllers => ./controllers
	github.com/egeback/download_media_api/internal/docs => ./docs
	github.com/egeback/download_media_api/internal/models => ./models
	github.com/egeback/download_media_api/internal/version => ./version
)

require (
	github.com/DispatchMe/go-work v0.6.0
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/egeback/download_media_api/internal/controllers v1.0.1
	github.com/egeback/download_media_api/internal/models v1.0.1
	github.com/egeback/download_media_api/internal/version v1.0.1
	github.com/egeback/play_media_api/internal/version v1.0.1
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-cmd/cmd v1.2.0
	github.com/gomodule/redigo v1.8.2
	github.com/google/uuid v1.1.1 // indirect
	github.com/hashicorp/go-version v1.2.1 // indirect
	github.com/liip/sheriff v0.0.0-20190308094614-91aa83a45a3d // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.5.1
	github.com/tomwei7/gin-jsonp v0.0.0-20191103091125-e5236eb5393d
)
