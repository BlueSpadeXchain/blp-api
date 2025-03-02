package infoHandler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/BlueSpadeXchain/blp-api/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
	"github.com/supabase-community/supabase-go"
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

func GetPairIdRequest(r *http.Request, parameters ...*GetPairRequestParams) (interface{}, error) {
	var params *GetPairRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetPairRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	pairId, err := getPairId(params.Pair)
	if err != nil {
		err_ := utils.ErrInternal(fmt.Sprintf("%v", err.Error()))
		utils.ErrInternal(err.Error())
		utils.LogError(err_.Message, err_.Details)
		return nil, err_
	}

	return &GetPairResponse{
		Pair:   params.Pair,
		PairId: pairId,
	}, nil

}

func GetLatestMetricSnapshotRequest(r *http.Request, supabaseClient *supabase.Client) (interface{}, error) {
	user, err := db.GetLatestMetricSnapshot(supabaseClient)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}

	return user, nil
}

type GetMetricsInRangeRequestParams struct {
	StartTime string `query:"start-time"`
	EndTime   string `query:"end-time"`
}

func GetMetricsInRangeRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetMetricsInRangeRequestParams) (interface{}, error) {
	var params *GetMetricsInRangeRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetMetricsInRangeRequestParams{}
	}

	// Parse query parameters
	query := r.URL.Query()
	startTime := query.Get("start_time")
	endTime := query.Get("end_time")

	// Validate time parameters
	if startTime == "" || endTime == "" {
		return nil, utils.ErrInternal("Missing required parameters: start_time or end_time")
	}

	// Parse timestamps to ensure validity
	_, err := time.Parse(time.RFC3339, params.StartTime)
	if err != nil {
		return nil, utils.ErrInternal("Invalid start_time format")
	}

	_, err = time.Parse(time.RFC3339, params.EndTime)
	if err != nil {
		return nil, utils.ErrInternal("Invalid end_time format")
	}

	// Fetch metrics from database
	metrics, err := db.GetMetricsInRange(supabaseClient, params.StartTime, params.EndTime)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}

	return metrics, nil
}
