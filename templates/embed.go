package templates

import "embed"

//go:embed services/*.yaml
var ServicesFS embed.FS
