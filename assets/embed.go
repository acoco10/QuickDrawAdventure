package assets

import (
	"embed"
)

//go:embed images
var ImagesDir embed.FS

//go:embed fonts
var Fonts embed.FS

//go:embed map
var Map embed.FS

//go:embed dialogueData
var Dialogue embed.FS
