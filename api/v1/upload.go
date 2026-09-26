package v1

import (
	"io"
	"net/http"
	"strconv"

	"ferry-srv/access"
	"ferry-srv/log"
	"ferry-srv/utils"

	"github.com/samber/mo"
)

func init() {
	http.HandleFunc("/api/v1/upload", func(w http.ResponseWriter, r *http.Request) { // geometry dash gamevars
		header := w.Header()
		utils.WriteHeaders(&header, http.MethodPost, false)

		if r.Method != http.MethodPost {
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

		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.WriteWebErr(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		res := access.ParseGameVarBody(body)
		if res.IsError() {
			utils.WriteWebErr(w, res.Error().Error(), http.StatusBadRequest)
			return
		}

		log.Debug("Uploading game settings save data for account of ID %v...", user.Account)

		bytes := utils.CompressZstd(body)

		gvRes := access.R2WriteGameVars(user.Account, bytes)
		if gvRes.IsError() {
			utils.WriteWebErr(w, gvRes.Error().Error(), http.StatusBadRequest)
			return
		}

		utils.WriteWebRes(w, mo.Some("Successfully uploaded game settings save data!"), http.StatusOK)
		log.Info("Successfully uploaded game settings save data for account of ID %v", user.Account)
	})

	// coming soon, just lazy rn

	http.HandleFunc("/api/v1/upload-mods", func(w http.ResponseWriter, r *http.Request) { // all mod settings
		header := w.Header()
		utils.WriteHeaders(&header, http.MethodPost, false)
		http.NotFound(w, r)
	})

	http.HandleFunc("/api/v1/upload-mods-saves", func(w http.ResponseWriter, r *http.Request) { // all mod save data (to be supporter-only)
		header := w.Header()
		utils.WriteHeaders(&header, http.MethodPost, false)
		http.NotFound(w, r)
	})
}
