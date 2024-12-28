package userHandler

import (
	"net/http"

	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
)

func WithdrawRequest(r *http.Request, parameters ...*WithdrawRequestParams) (interface{}, error) {
	var params *WithdrawRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &WithdrawRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func DespositRequest(r *http.Request, parameters ...*DespositRequestParams) (interface{}, error) {
	var params *DespositRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &DespositRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func UserDataRequest(r *http.Request, parameters ...*UserDataRequestParams) (interface{}, error) {
	var params *UserDataRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &UserDataRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func UserIdRequest(r *http.Request, parameters ...*UserIdRequestParams) (interface{}, error) {
	var params *UserIdRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &UserIdRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func AddAuthorizedWalletRequest(r *http.Request, parameters ...*AddAuthorizedWalletRequestParams) (interface{}, error) {
	var params *AddAuthorizedWalletRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &AddAuthorizedWalletRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func RemoveAuthorizedWalletRequest(r *http.Request, parameters ...*RemoveAuthorizedWalletRequestParams) (interface{}, error) {
	var params *RemoveAuthorizedWalletRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &RemoveAuthorizedWalletRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
