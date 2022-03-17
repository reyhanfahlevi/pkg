# New Relic Client

A custom package to aggregate newrelic API.

```go
import "github.com/ennobelprakoso/pkg/go/nrclient"

nr := nrclient.New(nrclient.Option{
    NRLoginCookies: `login cookies`,
    APIKey:         "api key",
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
```
