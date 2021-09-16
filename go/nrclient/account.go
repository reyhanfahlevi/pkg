package nrclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NRAccount struct
type NRAccount struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetAccountList get all available nr accounts
func (nr *NRClient) GetAccountList(ctx context.Context) ([]NRAccount, error) {
	var (
		resp = struct {
			Data struct {
				Actor struct {
					Accounts []NRAccount `json:"accounts"`
				} `json:"actor"`
			} `json:"data"`
		}{}
	)

	body := map[string]interface{}{
		"query": `{
		  actor {
			accounts {
			  id
			  name
			}
		  }
		}`,
	}

	rawBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", graphqlURL, bytes.NewBuffer(rawBody))
	if err != nil {
		return resp.Data.Actor.Accounts, err
	}

	req.Header.Add("api-key", nr.apiKey)
	req.Header.Add("content-type", "application/json")

	httpResp, err := nr.c.Do(req)
	if err != nil {
		return resp.Data.Actor.Accounts, err
	}
	defer httpResp.Body.Close()

	data, _ := ioutil.ReadAll(httpResp.Body)
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp.Data.Actor.Accounts, err
	}

	if httpResp.StatusCode > 399 {
		return resp.Data.Actor.Accounts, fmt.Errorf("error code %v: %s", httpResp.StatusCode, graphqlError(data))
	}

	return resp.Data.Actor.Accounts, err
}

// GetAddOnRoles get available addon roles under account
func (nr *NRClient) GetAddOnRoles(ctx context.Context, nrAccountID int64) ([]NRUserRoles, error) {
	var (
		roles []NRUserRoles
	)

	req, err := http.NewRequest("GET", fmt.Sprintf(getListOfUserRoles, nrAccountID), nil)
	if err != nil {
		return roles, err
	}

	req.Header.Add("cookie", fmt.Sprintf("login_service_login_newrelic_com_tokens=%s", nr.loginCookies))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")

	httpResp, err := nr.c.Do(req)
	if err != nil {
		return roles, err
	}
	defer httpResp.Body.Close()

	data, _ := ioutil.ReadAll(httpResp.Body)

	if httpResp.StatusCode > 399 {
		return roles, fmt.Errorf("error code %v: %s", httpResp.StatusCode, userManagementError(data))
	}

	err = json.Unmarshal(data, &roles)
	if err != nil {
		return roles, err
	}

	return []NRUserRoles{}, nil
}

type errResp struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func graphqlError(jsonStr []byte) string {
	resp := errResp{}

	err := json.Unmarshal(jsonStr, &resp)
	if err != nil {
		return ""
	}

	if len(resp.Errors) == 0 {
		return ""
	}

	return resp.Errors[0].Message
}
