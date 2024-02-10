package static

import "embed"

//go:embed index.html help.html
var Assets embed.FS
