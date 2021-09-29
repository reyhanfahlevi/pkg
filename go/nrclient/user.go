package nrclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type BaseRole string
type UserType int

const (
	UserTypeBasic UserType = 1
	UserTypeFull  UserType = 0

	UserBaseRoleAdmin      BaseRole = "admin"
	UserBaseRoleUser       BaseRole = "user"
	UserBaseRoleRestricted BaseRole = "restricted"
)

var (
	ErrUserNotFound = errors.New("user not found")
	validate        = validator.New()
)

// NRUser new relic user
type NRUser struct {
	UserID         int64         `json:"user_id"`
	AccountID      int64         `json:"account_id"`
	AccountName    string        `json:"account_name,omitempty"`
	FullName       string        `json:"full_name"`
	Email          string        `json:"email"`
	LastAccessAt   int64         `json:"last_access_at"`
	LastAccessTime time.Time     `json:"last_access_time"`
	Title          interface{}   `json:"title"`
	Status         string        `json:"status"`
	Roles          []NRUserRoles `json:"roles"`
	UserTierID     int           `json:"user_tier_id"`
}

// NRUserRoles struct
type NRUserRoles struct {
	ID          int64       `json:"id"`
	AccountID   interface{} `json:"account_id"`
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name"`
	Type        string      `json:"type"`
	BatchIds    []int       `json:"batch_ids"`
	GrantCount  int         `json:"grant_count"`
}

// GetUserUnderAccount get user under specific account
func (nr *NRClient) GetUserUnderAccount(ctx context.Context, email string, nrAccountID int64) (NRUser, error) {
	var (
		resp = []NRUser{}

		user = NRUser{}
	)

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(getListOfUserUnderAccount, nrAccountID), nil)
	if err != nil {
		return user, err
	}

	req.Header.Add("cookie", fmt.Sprintf("login_service_login_newrelic_com_tokens=%s", nr.loginCookies))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")

	httpResp, err := nr.c.Do(req)
	if err != nil {
		return user, err
	}
	defer httpResp.Body.Close()

	data, _ := ioutil.ReadAll(httpResp.Body)

	if httpResp.StatusCode > 399 {
		return user, fmt.Errorf("error code %v: %s", httpResp.StatusCode, userManagementError(data))
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return user, err
	}

	for _, r := range resp {
		if r.Email == email {
			return r, nil
		}
	}

	return user, ErrUserNotFound
}

// FindUserAccount find user from all available account
func (nr *NRClient) FindUserAccount(ctx context.Context, email string) ([]NRUser, error) {
	var (
		users   = []NRUser{}
		wg      sync.WaitGroup
		mux     sync.Mutex
		errChan = make(chan error)
	)

	account, err := nr.GetAccountList(ctx)
	if err != nil {
		return users, err
	}

	for _, a := range account {

		wg.Add(1)
		go func(acc NRAccount) {
			defer wg.Done()

			user, err := nr.GetUserUnderAccount(ctx, email, acc.ID)
			if err != nil && err != ErrUserNotFound {
				errChan <- err
				return
			}

			if err == ErrUserNotFound {
				return
			}

			mux.Lock()
			user.AccountName = acc.Name
			user.LastAccessTime = time.Unix(user.LastAccessAt, 0)
			users = append(users, user)
			mux.Unlock()
		}(a)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	// catch any error and blocking
	for err := range errChan {
		if err != nil {
			return users, err
		}
	}

	return users, nil
}

// ParamCreateUser struct
type ParamCreateUser struct {
	FullName   string `validate:"required"`
	Email      string `validate:"required,email"`
	AccountID  int64  `validate:"required"`
	AddOnRoles []int64
	UserType   UserType
	BaseRole   BaseRole `validate:"required"`
}

// CreateUser create new user
func (nr *NRClient) CreateUser(ctx context.Context, param ParamCreateUser) error {
	var (
		resp struct {
			Success        bool   `json:"success"`
			WelcomeMessage string `json:"welcome_message"`
			UserID         int64  `json:"user_id"`
		}
	)

	err := validate.Struct(param)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	body := map[string]interface{}{
		"account_view": map[string]interface{}{
			"user": map[string]interface{}{
				"full_name": param.FullName,
				"email":     param.Email,
			},
			"level":        param.BaseRole,
			"user_tier_id": param.UserType,
		},
	}

	rawBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf(createUserUnderAccount, param.AccountID), bytes.NewBuffer(rawBody))
	if err != nil {
		return err
	}

	req.Header.Add("cookie", fmt.Sprintf("login_service_login_newrelic_com_tokens=%s", nr.loginCookies))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	httpResp, err := nr.c.Do(req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	data, _ := ioutil.ReadAll(httpResp.Body)

	if httpResp.StatusCode > 399 {
		return fmt.Errorf("error code %v: %s", httpResp.StatusCode, userManagementError(data))
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}

	if len(param.AddOnRoles) > 0 {
		return nr.UpdateUserAddOnRoles(ctx, param.AccountID, resp.UserID, param.AddOnRoles)
	}

	return nil
}

// UpdateUserAddOnRoles update user addons role
func (nr *NRClient) UpdateUserAddOnRoles(ctx context.Context, nrAccountID int64, userID int64, roles []int64) error {
	body := map[string]interface{}{
		"roles": roles,
	}

	rawBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf(updateUsers, nrAccountID, userID), bytes.NewBuffer(rawBody))
	if err != nil {
		return err
	}

	req.Header.Add("cookie", fmt.Sprintf("login_service_login_newrelic_com_tokens=%s", nr.loginCookies))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	httpResp, err := nr.c.Do(req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	data, _ := ioutil.ReadAll(httpResp.Body)

	if httpResp.StatusCode > 399 {
		return fmt.Errorf("error code %v: %s", httpResp.StatusCode, userManagementError(data))
	}

	return nil
}

// RemoveUserFromAccount remove user from nr account
func (nr *NRClient) RemoveUserFromAccount(ctx context.Context, email string, nrAccountID int64) error {
	user, err := nr.GetUserUnderAccount(ctx, email, nrAccountID)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf(updateUsers, nrAccountID, user.UserID), nil)
	if err != nil {
		return err
	}

	req.Header.Add("cookie", fmt.Sprintf("login_service_login_newrelic_com_tokens=%s", nr.loginCookies))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	httpResp, err := nr.c.Do(req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	data, _ := ioutil.ReadAll(httpResp.Body)

	if httpResp.StatusCode > 399 {
		return fmt.Errorf("error code %v: %s", httpResp.StatusCode, userManagementError(data))
	}

	return nil
}

// BulkCreateUserSummary struct
type BulkCreateUserSummary struct {
	Data    ParamCreateUser
	Success bool
	Err     error
}

// BulkCreateUser create new user in bulk
func (nr *NRClient) BulkCreateUser(ctx context.Context, data ...ParamCreateUser) []BulkCreateUserSummary {
	var (
		summary []BulkCreateUserSummary
	)
	for _, d := range data {
		tmp := BulkCreateUserSummary{
			Data:    d,
			Success: true,
		}

		err := nr.CreateUser(ctx, d)
		if err != nil {
			tmp.Err = err
			tmp.Success = false
		}

		summary = append(summary, tmp)
	}

	return summary
}

func userManagementError(jsonStr []byte) string {
	resp := struct {
		Error string `json:"error"`
	}{}

	err := json.Unmarshal(jsonStr, &resp)
	if err != nil {
		return ""
	}

	return resp.Error
}
