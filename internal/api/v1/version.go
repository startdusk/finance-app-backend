package v1

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/startdusk/finance-app-backend/internal/config"
)

// API for returning version
// When server starts, we set version and than use it if necesary.

// ServerVersion represents the server version.
type ServerVersion struct {
	Version string `json:"version"`
}

// Marshaled JSON
var versionJSON []byte

func init() {
	var err error
	versionJSON, err = json.Marshal(ServerVersion{
		Version: config.Version,
	})
	if err != nil {
		panic(err)
	}
}

// VersionHandler serves version information
func VersionHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(versionJSON); err != nil {
		logrus.WithError(err).Debug("error writing version")
	}
}
