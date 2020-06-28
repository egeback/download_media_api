module github.com/egeback/download_media_api/internal/controllers

go 1.14

replace (
	github.com/egeback/download_media_api/internal/controllers => ../controllers
	github.com/egeback/download_media_api/internal/docs => ../docs
	github.com/egeback/download_media_api/internal/models => ../models
	github.com/egeback/download_media_api/internal/version => ../version
)

require (
	github.com/egeback/download_media_api/internal/actions v0.0.0-20200622165410-818df93be324
	github.com/egeback/download_media_api/internal/models v0.0.0-20200622165410-818df93be324
	github.com/gin-gonic/gin v1.6.3
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-version v1.2.1
	github.com/liip/sheriff v0.0.0-20190308094614-91aa83a45a3d
)
