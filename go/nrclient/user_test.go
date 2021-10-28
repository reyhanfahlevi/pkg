package nrclient

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/h2non/gock.v1"
)

func jsonRespCheck(t *testing.T, want interface{}, actual interface{}) {
	wantJSON, _ := json.Marshal(want)
	actualJSON, _ := json.Marshal(actual)

	assert.JSONEq(t, string(wantJSON), string(actualJSON))
}

func TestNRClient_BulkCreateUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		data []ParamCreateUser
	}
	tests := []struct {
		name string
		args args
		mock func(a args)
		want []BulkCreateUserSummary
	}{
		{
			name: "Test Bulk Create User Success",
			args: args{
				ctx: context.Background(),
				data: []ParamCreateUser{
					{
						FullName:  "Test",
						Email:     "test@test.com",
						AccountID: 1010101,
						UserType:  UserTypeBasic,
						BaseRole:  UserBaseRoleRestricted,
					},
				},
			},
			mock: func(a args) {
				gock.New(fmt.Sprintf(createUserUnderAccount, a.data[0].AccountID)).
					Post("").Reply(1).JSON(map[string]interface{}{"success": true})
			},
			want: []BulkCreateUserSummary{
				{
					Data: ParamCreateUser{
						FullName:  "Test",
						Email:     "test@test.com",
						AccountID: 1010101,
						UserType:  UserTypeBasic,
						BaseRole:  UserBaseRoleRestricted,
					},
					Success: true,
				},
			},
		}, {
			name: "Test Bulk Create User Partial Failed",
			args: args{
				ctx: context.Background(),
				data: []ParamCreateUser{
					{
						FullName:  "Test",
						Email:     "test@test.com",
						AccountID: 1010101,
						UserType:  UserTypeBasic,
						BaseRole:  UserBaseRoleRestricted,
					}, {
						FullName:  "Test2",
						Email:     "test2@test.com",
						AccountID: 1010101,
						UserType:  UserTypeBasic,
						BaseRole:  UserBaseRoleRestricted,
					},
				},
			},
			mock: func(a args) {
				gock.New(fmt.Sprintf(createUserUnderAccount, a.data[0].AccountID)).
					Post("").Reply(200).JSON(map[string]interface{}{"success": true})

				gock.New(fmt.Sprintf(createUserUnderAccount, a.data[1].AccountID)).
					Post("").Reply(200).BodyString(`{`)
			},
			want: []BulkCreateUserSummary{
				{
					Data: ParamCreateUser{
						FullName:  "Test",
						Email:     "test@test.com",
						AccountID: 1010101,
						UserType:  UserTypeBasic,
						BaseRole:  UserBaseRoleRestricted,
					},
					Success: true,
				}, {
					Data: ParamCreateUser{
						FullName:  "Test2",
						Email:     "test2@test.com",
						AccountID: 1010101,
						UserType:  UserTypeBasic,
						BaseRole:  UserBaseRoleRestricted,
					},
					Success: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nr := New(Option{})
			gock.Flush()

			tt.mock(tt.args)
			got := nr.BulkCreateUser(tt.args.ctx, tt.args.data...)
			for i, g := range got {
				assert.Equal(t, g.Data, tt.want[i].Data)
				assert.Equal(t, g.Success, tt.want[i].Success)
				assert.Equal(t, g.Success, tt.want[i].Success)
			}
		})
	}
}

func Test_userManagementError(t *testing.T) {
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
				jsonStr: []byte(""),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := userManagementError(tt.args.jsonStr); got != tt.want {
				t.Errorf("userManagementError() = %v, want %v", got, tt.want)
			}
		})
	}
}
