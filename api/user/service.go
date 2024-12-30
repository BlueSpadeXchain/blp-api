package userHandler

import (
	"net/http"

	"github.com/BlueSpadeXchain/blp-api/pkg/db"
	"github.com/BlueSpadeXchain/blp-api/pkg/utils"
	"github.com/supabase-community/supabase-go"
)

func WithdrawRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*WithdrawRequestParams) (interface{}, error) {
	var params *WithdrawRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &WithdrawRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	return nil, nil
}

func DespositRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*DespositRequestParams) (interface{}, error) {
	var params *DespositRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &DespositRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	return nil, nil
}

func UserDataRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*UserDataRequestParams) (interface{}, error) {
	var params *UserDataRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &UserDataRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	return nil, nil
}

func GetUserByIdRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetUserByIdRequestParams) (interface{}, error) {
	var params *GetUserByIdRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetUserByIdRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	user, err := db.GetUserByUserId(supabaseClient, params.UserId)
	if err != nil {
		return nil, utils.ErrInternal(err.Error())
	}

	return user, nil
}

func GetUserByAddressRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*GetUserByAddressRequestParams) (interface{}, error) {
	var params *GetUserByAddressRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &GetUserByAddressRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	user, err := db.GetOrCreateUser(supabaseClient, params.Address, params.AddressType)
	if err != nil {
		utils.LogError("db GetOrCreateUser failed", err.Error())
		return nil, utils.ErrInternal(err.Error())
	}

	return user, nil
}

func AddAuthorizedWalletRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*AddAuthorizedWalletRequestParams) (interface{}, error) {
	var params *AddAuthorizedWalletRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &AddAuthorizedWalletRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	return nil, nil
}

func RemoveAuthorizedWalletRequest(r *http.Request, supabaseClient *supabase.Client, parameters ...*RemoveAuthorizedWalletRequestParams) (interface{}, error) {
	var params *RemoveAuthorizedWalletRequestParams

	if len(parameters) > 0 {
		params = parameters[0]
	} else {
		params = &RemoveAuthorizedWalletRequestParams{}
	}

	if r != nil {
		if err := utils.ParseAndValidateParams(r, &params); err != nil {
			utils.LogError("failed to parse params", err.Error())
			return nil, utils.ErrInternal(err.Error())
		}
	}

	return nil, nil
}
