package assets

import (
	"embed"
)

//go:embed images
var ImagesDir embed.FS

//go:embed fonts
var Fonts embed.FS
