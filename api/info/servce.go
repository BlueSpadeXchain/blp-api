package infoHandler

import (
	"net/http"

	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
)

func VersionRequest(r *http.Request, parameters ...interface{}) (interface{}, error) {
	return utils.VersionResponse{
		Version: Version,
	}, nil
}
