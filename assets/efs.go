package assets

import (
	"embed"
)

//go:embed "emails" "migrations" "templates"
var EmbeddedFiles embed.FS
