package main

import (
	"context"

	"github.com/ennobelprakoso/pkg/go/log"
	"github.com/ennobelprakoso/pkg/go/nrclient"
)

func main() {
	nr := nrclient.New(nrclient.Option{
		NRLoginCookies: ``,
		APIKey:         "",
	})

	users, err := nr.FindUserAccount(context.Background(), "test@test.com")
	if err != nil {
		log.Error(err)
		return
	}

	for _, u := range users {
		log.Info(u.AccountID)
		log.Info(u.AccountName)
		log.Info(u.LastAccessTime)
	}
}
