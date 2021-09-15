package nrclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
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
}

// GetUserUnderAccount get user under specific account
func (nr *NRClient) GetUserUnderAccount(ctx context.Context, email string, nrAccountID int64) (NRUser, error) {
	var (
		resp = []NRUser{}

		user = NRUser{}
	)

	req, err := http.NewRequest("GET", fmt.Sprintf(getListOfUserUnderAccount, nrAccountID), nil)
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
