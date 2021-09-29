package nrclient

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/pkg/errors"
	"gopkg.in/h2non/gock.v1"
)

func TestNRClient_GetAccountList(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		mock    func(a args)
		want    []NRAccount
		wantErr bool
	}{
		{
			name: "Test Success",
			args: args{
				ctx: context.Background(),
			},
			mock: func(a args) {
				got := gqlGetAccountResp{}
				got.Data.Actor.Accounts = []NRAccount{
					{
						ID:   1,
						Name: "NRAccount1",
					},
				}

				gock.New("https://api.newrelic.com").Post("/graphql").Reply(200).JSON(got)
			},
			want: []NRAccount{
				{
					ID:   1,
					Name: "NRAccount1",
				},
			},
			wantErr: false,
		}, {
			name: "Test 404",
			args: args{
				ctx: context.Background(),
			},
			mock: func(a args) {
				got := gqlGetAccountResp{}
				got.Data.Actor.Accounts = []NRAccount{
					{
						ID:   1,
						Name: "NRAccount1",
					},
				}

				gock.New("https://api.newrelic.com").Post("/graphql").Reply(404).JSON(errResp{Errors: []struct {
					Message string `json:"message"`
				}{
					{
						Message: "some err",
					},
				}})
			},
			want:    []NRAccount{},
			wantErr: true,
		}, {
			name: "Test Error",
			args: args{
				ctx: context.Background(),
			},
			mock: func(a args) {
				got := gqlGetAccountResp{}
				got.Data.Actor.Accounts = []NRAccount{
					{
						ID:   1,
						Name: "NRAccount1",
					},
				}

				gock.New("https://api.newrelic.com").Post("/graphql").
					ReplyError(errors.New("some err"))
			},
			want:    []NRAccount{},
			wantErr: true,
		}, {
			name: "Test Unmarshal Failed",
			args: args{
				ctx: context.Background(),
			},
			mock: func(a args) {
				got := gqlGetAccountResp{}
				got.Data.Actor.Accounts = []NRAccount{
					{
						ID:   1,
						Name: "NRAccount1",
					},
				}

				gock.New("https://api.newrelic.com").Post("/graphql").Reply(200).
					BodyString("{")
			},
			want:    []NRAccount{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nr := New(Option{})
			gock.Flush()

			tt.mock(tt.args)
			got, err := nr.GetAccountList(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAccountList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNRClient_GetAddOnRoles(t *testing.T) {
	type fields struct {
		c            *httpclient.Client
		loginCookies string
		apiKey       string
	}
	type args struct {
		ctx         context.Context
		nrAccountID int64
	}
	tests := []struct {
		name    string
		args    args
		mock    func(a args)
		want    []NRUserRoles
		wantErr bool
	}{
		{
			name: "Test Success",
			args: args{
				ctx:         context.Background(),
				nrAccountID: 1,
			},
			mock: func(a args) {
				gock.New(fmt.Sprintf(getListOfUserRoles, a.nrAccountID)).Get("").Reply(200).JSON([]NRUserRoles{})
			},
			want:    []NRUserRoles{},
			wantErr: false,
		}, {
			name: "Test 404",
			args: args{
				ctx:         context.Background(),
				nrAccountID: 1,
			},
			mock: func(a args) {
				gock.New(fmt.Sprintf(getListOfUserRoles, a.nrAccountID)).Get("").Reply(http.StatusUnprocessableEntity).JSON(struct {
					Error string `json:"error"`
				}{
					Error: "some err",
				})
			},
			want:    []NRUserRoles{},
			wantErr: true,
		}, {
			name: "Test Some Error",
			args: args{
				ctx:         context.Background(),
				nrAccountID: 1,
			},
			mock: func(a args) {
				gock.New(fmt.Sprintf(getListOfUserRoles, a.nrAccountID)).Get("").
					ReplyError(errors.New("some err"))
			},
			want:    []NRUserRoles{},
			wantErr: true,
		}, {
			name: "Test Unmarshal Failed",
			args: args{
				ctx:         context.Background(),
				nrAccountID: 1,
			},
			mock: func(a args) {
				gock.New(fmt.Sprintf(getListOfUserRoles, a.nrAccountID)).Get("").Reply(200).BodyString("{")
			},
			want:    []NRUserRoles{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nr := New(Option{})
			gock.Flush()

			tt.mock(tt.args)
			got, err := nr.GetAddOnRoles(tt.args.ctx, tt.args.nrAccountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddOnRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddOnRoles() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_graphqlError(t *testing.T) {
	type args struct {
		jsonStr []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1",
			args: args{
				jsonStr: []byte{},
			},
			want: "",
		}, {
			name: "Test 2",
			args: args{
				jsonStr: []byte("{}"),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := graphqlError(tt.args.jsonStr); got != tt.want {
				t.Errorf("graphqlError() = %v, want %v", got, tt.want)
			}
		})
	}
}
