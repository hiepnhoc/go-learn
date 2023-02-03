package esign

import (
	"acbs.com.vn/account-service/config"
	"acbs.com.vn/account-service/pkg/logger"
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
)

const (
	defaultContentType = "application/json; charset=utf-8"
)

const (
	post = "post"
)

const (
	EsignGetDetail = "/esign/"
)

const (
	AuthTokenKey = "Authorization"
	ClientIdKey  = "ClientID"
)

const (
	clientCredentials = "client_credentials"
)

type EsignResponseBase struct {
	ResultCode int    `json:"result_code"`
	Message    string `json:"message"`
	Data       string `json:"data"`
}

type EsignClient interface {
	GetDetail(ctx context.Context) (*EsignResponseBase, error)
}

type esignClient struct {
	log    logger.Logger
	config *config.Config
	client *resty.Client
}

func NewEsignClient(log logger.Logger, config *config.Config, client *resty.Client) EsignClient {
	return &esignClient{log: log, config: config, client: client}
}

func (h *esignClient) GetDetail(ctx context.Context) (*EsignResponseBase, error) {

	type Request struct{}
	ctx, err := h.getAuthToken(ctx)
	if err != nil {
		return nil, err
	}

	//request := &Request{}

	//invokeResponse, err := h.invoke(ctx, TrueIdPathCreateRequest, request)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//out := &TrueIdResponseBase{}
	//
	//if err := convert.JsonDecode(invokeResponse.Data, out); err != nil {
	//	return nil, convert.NewError("TrueIdClient-CreateRequest", "JsonDecode", invokeResponse.Data, out, err)
	//}

	return nil, nil

}

func (h *esignClient) getAuthToken(ctx context.Context) (context.Context, error) {

	type AuthRequest struct {
		GrantType    string `json:"grant_type"`
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}

	type AuthResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	//authRequest := &AuthRequest{
	//	GrantType:    clientCredentials,
	//	ClientId:     h.config.KeyCloak.ClientId,
	//	ClientSecret: h.config.KeyCloak.ClientSecret,
	//}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get("https://dapi-dev.acbs.tech/" + EsignGetDetail + "123123123")

	if err != nil {
		fmt.Println(err)
		return ctx, nil
	}

	fmt.Println(resp)

	resp1, err1 := client.R().
		SetFormData(map[string]string{
			"client_id":     h.config.KeyCloak.ClientId,
			"client_secret": h.config.KeyCloak.ClientSecret,
			"grant_type":    clientCredentials,
		}).
		Post("https://sso.acbs.tech/auth/realms/applications/protocol/openid-connect/token")

	if err1 != nil {
		fmt.Println(err)
		return ctx, nil
	}

	fmt.Println(resp1)

	//var authRequestData []byte
	//authRequestData, err = convert.JsonEncode(authRequest)
	//
	//authResponse := AuthResponse{}
	//
	//if err := convert.JsonDecode(invokeResponse.Data, &authResponse); err != nil {
	//	return nil, convert.NewError("TrueIdClient-WithTrueIdAuthToken", "JsonDecode", invokeResponse.Data, authResponse, err)
	//}
	//
	//token := fmt.Sprintf("%s %s", authResponse.TokenType, authResponse.AccessToken)
	//
	//ctx = context.WithValue(ctx, trueIdAuthToken, token)
	//
	//if err := h.client.SaveState(ctx, h.config.Dapr.StateStore, trueIdAuthToken, []byte(token), map[string]string{"ttlInSeconds": "300"}); err != nil {
	//	return nil, NewTrueIdStageError("TrueIdClient-WithTrueIdAuthToken", "SaveState", trueIdAuthToken, err)
	//}
	//
	//h.log.Debugf("trueIdClient withTrueIdAuthToken: cached  %s : %s", trueIdAuthToken, token)

	return ctx, nil
}
