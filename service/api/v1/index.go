package v1

import (
	"net/http"
	"strconv"

	"ferry-srv/access"
	"ferry-srv/log"
	"ferry-srv/utils"

	"github.com/samber/mo"
)

func authUser(r *http.Request) (mo.Result[utils.ArgonUser], int) {
	q := r.URL.Query()

	accStr := q.Get("account_id")

	acc, err := strconv.Atoi(accStr)
	if err != nil {
		return mo.Errf[utils.ArgonUser]("Failed to parse account ID"), http.StatusBadRequest
	}

	token := q.Get("authtoken")

	userRes := access.ValidateArgonUser(&utils.ArgonUser{Account: acc, Token: token}, false)
	if userRes.IsError() {
		return userRes, http.StatusUnauthorized
	}
	return mo.Ok(userRes.MustGet()), http.StatusOK
}

func init() {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("API endpoint pinged!")

		header := w.Header()
		utils.WriteHeaders(&header, http.MethodGet, false)

		utils.WriteWebRes(w, mo.Some("Pong!"), http.StatusOK)
	})

	http.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("v1 API pinged!")

		header := w.Header()
		utils.WriteHeaders(&header, http.MethodGet, false)

		utils.WriteWebRes(w, mo.Some("Pong!"), http.StatusOK)
	})
}
