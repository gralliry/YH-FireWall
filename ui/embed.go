package assets

import "embed"

//go:embed dist
var WebServerStaticFS embed.FS
