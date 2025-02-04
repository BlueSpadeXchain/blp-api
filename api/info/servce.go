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

func GetPairsRequest(r *http.Request, parameters ...interface{}) (interface{}, error) {
	return GetPairsResponse{
		Pairs: Pairs,
	}, nil
}

func GetPairsAndIdsRequest(r *http.Request, parameters ...interface{}) (interface{}, error) {
	return &PairAndIds, nil
}

func GetPairIdsRequest(r *http.Request, parameters ...interface{}) (interface{}, error) {
	return &Pairs, nil
}
