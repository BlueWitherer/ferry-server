package v1

import (
	"net/http"
	"strconv"

	"ferry-srv/access"
	"ferry-srv/log"
	"ferry-srv/utils"
)

func init() {
	http.HandleFunc("/api/v1/download", func(w http.ResponseWriter, r *http.Request) { // geometry dash gamevars
		header := w.Header()
		utils.WriteHeaders(&header, http.MethodGet, true)

		if r.Method != http.MethodGet {
			utils.WriteWebErrMethod(w)
			return
		}

		q := r.URL.Query()

		accStr := q.Get("account_id")

		acc, err := strconv.Atoi(accStr)
		if err != nil {
			utils.WriteWebErr(w, "Failed to parse account ID", http.StatusBadRequest)
			return
		}

		token := q.Get("authtoken")

		userRes := access.ValidateArgonUser(&utils.ArgonUser{Account: acc, Token: token}, false)
		if userRes.IsError() {
			utils.WriteWebErr(w, userRes.Error().Error(), http.StatusUnauthorized)
			return
		}

		user := userRes.MustGet()

		gvRes := access.R2ReadGameVarsRaw(user.Account)
		if gvRes.IsError() {
			utils.WriteWebErr(w, gvRes.Error().Error(), http.StatusBadRequest)
			return
		}

		log.Debug("Streaming game settings save data for account of ID %v...", user.Account)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(gvRes.MustGet()); err != nil {
			utils.WriteWebErr(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info("Streamed game settings save data for account of ID %v", user.Account)
	})
}
