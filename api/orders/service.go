package orderHandler

import (
	"net/http"

	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
)

// each of these three request type add to requests pool to be processed

// func OpenOrderRequest(r *http.Request, parameters ...*interface{}) (interface{}, error) {
// 	return nil, nil
// }

// func CloseOrderRequest(r *http.Request, parameters ...*interface{}) (interface{}, error) {
// 	return nil, nil
// }

func GetOrdersRequest(r *http.Request, parameters ...*interface{}) (interface{}, error) {
	return nil, nil
}

func OrderRequest(r *http.Request, parameters ...*OrderRequestParams) (interface{}, error) {
	var params *OrderRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &OrderRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, err
		}
	}

	// if err := validateOrderRequest(params); err != nil {
	// 	return nil, err
	// }

	// validate table for perp-id
	// check if db has

	return nil, nil
}

func CloseOrder(r *http.Request, parameters ...*OrderCloseParams) (interface{}, error) {
	var params *OrderCloseParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &OrderCloseParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, err
		}
	}

	// if err := validateOrderClose(params); err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
