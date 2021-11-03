package docs

import (
	_ "embed"
)

//go:embed openapi.json
var Openapi string
