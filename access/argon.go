package access

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/BlueWitherer/GDDataSyncServer/log"
	"github.com/BlueWitherer/GDDataSyncServer/utils"
	"github.com/samber/mo"

	"github.com/patrickmn/go-cache"
)

var argonCache = cache.New(15*time.Minute, 10*time.Minute)
var invalids = cache.New(5*time.Minute, 10*time.Minute)

func getToken() (string, error) {
	token := os.Getenv("ARGON_TOKEN")
	if token == "" {
		return "", fmt.Errorf("env for argon token is not defined!")
	} else {
		return token, nil
	}
}

func ValidateArgonUser(user *utils.ArgonUser, strong bool) mo.Result[bool] {
	if val, found := invalids.Get(fmt.Sprintf("%d", user.Account)); found {
		return mo.Errf[bool]("Argon token %s is invalid", val.(string))
	}

	if _, found := argonCache.Get(fmt.Sprintf("%d", user.Account)); found {
		return mo.Ok(found)
	}

	var authUrl string
	if strong {
		authUrl = "https://argon.globed.dev/v1/validation/check-strong"
	} else {
		authUrl = "https://argon.globed.dev/v1/validation/check"
	}

	u, err := url.Parse(authUrl)
	if err != nil {
		return mo.Err[bool](err)
	} else {
		log.Trace("Argon URL parsed for account of ID %v", user.Account)
	}

	q := u.Query()
	q.Set("account_id", fmt.Sprintf("%v", user.Account))
	q.Set("authtoken", user.Token)
	if strong {
		q.Set("user_id", fmt.Sprintf("%v", user.User))
		q.Set("username", user.Username)
	}
	u.RawQuery = q.Encode()

	log.Trace("Argon validation parameters: account_id=%d (type check: %T), authtoken length=%d", user.Account, user.Account, len(user.Token))
	log.Trace("Full Argon URL being requested: %s", u.String())

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return mo.Err[bool](err)
	} else {
		log.Trace("Argon request object constructed for account of ID %v", user.Account)
	}

	req.Header.Set("User-Agent", "GDDataSync/1.0")

	argon, err := getToken()
	if err != nil {
		log.Warn("Failed to get Argon API token: %s", err.Error())
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", argon))
	}

	log.Trace("Sending request to Argon server: %s", u.String())
	client := &http.Client{Timeout: 15 * time.Second}

	resp, reqErr := client.Do(req)
	if reqErr != nil {
		return mo.Err[bool](reqErr)
	}
	defer resp.Body.Close()

	log.Trace("Argon status code received: %d, for account of ID %v", resp.StatusCode, user.Account)

	if resp.StatusCode != http.StatusOK {
		return mo.Errf[bool]("argon server returned status code %d", resp.StatusCode)
	}

	var valid utils.ArgonValidation
	if err := json.NewDecoder(resp.Body).Decode(&valid); err != nil {
		return mo.Errf[bool]("failed to parse argon response: %v", err)
	} else {
		log.Trace("Argon status of account of ID %v retrieved", user.Account)
	}

	if valid.Valid {
		log.Info("Argon status of account of ID %v is valid", user.Account)
		return mo.Ok(valid.Valid)
	}

	invalids.Set(fmt.Sprintf("%d", user.Account), user.Token, cache.DefaultExpiration)
	return mo.Errf[bool]("cause: %s", valid.Cause)
}
