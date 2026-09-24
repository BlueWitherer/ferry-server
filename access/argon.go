package access

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"ferry-srv/log"
	"ferry-srv/utils"

	"github.com/samber/mo"

	"github.com/patrickmn/go-cache"
)

// i dont wanna spam the server lowkey...
var argonCache = cache.New(3*time.Minute, 5*time.Minute)
var invalids = cache.New(2*time.Minute, 3*time.Minute)

var rlToken string

func getToken() (string, error) {
	if rlToken == "" {
		return "", fmt.Errorf("env for argon ratelimit token is not defined!")
	} else {
		return rlToken, nil
	}
}

func ValidateArgonUser(user *utils.ArgonUser, strong bool) mo.Result[utils.ArgonUser] {
	if _, found := invalids.Get(user.Token); found {
		return mo.Errf[utils.ArgonUser]("Argon ratelimit token %s is invalid", user.Token)
	}

	if u, found := argonCache.Get(fmt.Sprintf("%v", user.Account)); found {
		return mo.Ok(u.(utils.ArgonUser))
	}

	var authUrl string
	if strong {
		authUrl = "https://argon.globed.dev/v1/validation/check-strong"
	} else {
		authUrl = "https://argon.globed.dev/v1/validation/check"
	}

	u, err := url.Parse(authUrl)
	if err != nil {
		return mo.Err[utils.ArgonUser](err)
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

	log.Trace("Argon validation parameters: account_id=%v (type check: %T), authtoken length=%v", user.Account, user.Account, len(user.Token))
	log.Trace("Full Argon URL being requested: %s", u.String())

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return mo.Err[utils.ArgonUser](err)
	} else {
		log.Trace("Argon request object constructed for account of ID %v", user.Account)
	}

	req.Header.Set("User-Agent", "GDDataSync/1.0")

	argon, err := getToken()
	if err != nil {
		log.Warn("Failed to get Argon API rlToken: %s", err.Error())
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", argon))
	}

	log.Trace("Sending request to Argon server: %s", u.String())
	client := &http.Client{Timeout: 15 * time.Second}

	resp, reqErr := client.Do(req)
	if reqErr != nil {
		return mo.Err[utils.ArgonUser](reqErr)
	}
	defer resp.Body.Close()

	log.Trace("Argon status code received: %v, for account of ID %v", resp.StatusCode, user.Account)

	if resp.StatusCode != http.StatusOK {
		return mo.Errf[utils.ArgonUser]("argon server returned status code %v", resp.StatusCode)
	}

	var valid utils.ArgonValidation
	if err := json.NewDecoder(resp.Body).Decode(&valid); err != nil {
		return mo.Errf[utils.ArgonUser]("failed to parse argon response: %v", err)
	} else {
		log.Trace("Argon status of account of ID %v retrieved", user.Account)
	}

	if valid.Valid {
		log.Info("Argon status of account of ID %v is valid", user.Account)
		user.Token = ""

		argonCache.Set(fmt.Sprintf("%v", user.Account), *user, cache.DefaultExpiration)
		return mo.Ok(*user)
	}

	invalids.Set(user.Token, true, cache.DefaultExpiration)
	return mo.Errf[utils.ArgonUser]("cause: %s", valid.Cause)
}

func init() {
	rlToken = os.Getenv("ARGON_TOKEN")
}
