package config

import (
	"github.com/namsral/flag"
)

// DataDirectory is the path used for loading templates/database migrations
var DataDirectory = flag.String("data-directory", "", "Path for loading templates and migration scripts.")
