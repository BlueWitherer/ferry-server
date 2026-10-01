package v1

import (
	"io"
	"net/http"

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

		userRes, code := authUser(r)
		if userRes.IsError() {
			utils.WriteWebErr(w, userRes.Error().Error(), code)
			return
		}
		user := userRes.MustGet()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.WriteWebErr(w, err.Error(), http.StatusRequestEntityTooLarge)
			return
		}
		defer r.Body.Close()

		res := access.ParseGameVarBody(body)
		if res.IsError() {
			utils.WriteWebErr(w, res.Error().Error(), http.StatusBadRequest)
			return
		}

		log.Debug("Uploading game settings save data for account of ID %v...", user.Account)

		gvRes := access.R2WriteGameVars(user.Account, body)
		if gvRes.IsError() {
			utils.WriteWebErr(w, gvRes.Error().Error(), http.StatusBadRequest)
			return
		}

		utils.WriteWebRes(w, mo.Some("Successfully uploaded game settings save data!"), http.StatusOK)
		log.Info("Successfully uploaded game settings save data for account of ID %v", user.Account)
	})

	http.HandleFunc("/api/v1/upload-geode", func(w http.ResponseWriter, r *http.Request) { // geode settings
		header := w.Header()
		utils.WriteHeaders(&header, http.MethodPost, false)

		if r.Method != http.MethodPost {
			utils.WriteWebErrMethod(w)
			return
		}

		userRes, code := authUser(r)
		if userRes.IsError() {
			utils.WriteWebErr(w, userRes.Error().Error(), code)
			return
		}
		user := userRes.MustGet()

		body, err := io.ReadAll(io.LimitReader(r.Body, 4<<10))
		if err != nil {
			utils.WriteWebErr(w, err.Error(), http.StatusRequestEntityTooLarge)
			return
		}
		defer r.Body.Close()

		log.Debug("Uploading Geode loader settings save data for account of ID %v...", user.Account)

		geodeRes := access.R2WriteGeodeSettings(user.Account, body)
		if geodeRes.IsError() {
			utils.WriteWebErr(w, geodeRes.Error().Error(), http.StatusBadRequest)
			return
		}

		utils.WriteWebRes(w, mo.Some("Successfully uploaded Geode loader settings save data!"), http.StatusOK)
		log.Info("Successfully uploaded Geode loader settings save data for account of ID %v", user.Account)
	})

	http.HandleFunc("/api/v1/upload-geode-mods", func(w http.ResponseWriter, r *http.Request) { // all mod settings
		header := w.Header()
		utils.WriteHeaders(&header, http.MethodPost, false)

		if r.Method != http.MethodPost {
			utils.WriteWebErrMethod(w)
			return
		}

		userRes, code := authUser(r)
		if userRes.IsError() {
			utils.WriteWebErr(w, userRes.Error().Error(), code)
			return
		}
		user := userRes.MustGet()

		body, err := io.ReadAll(io.LimitReader(r.Body, 2<<20))
		if err != nil {
			utils.WriteWebErr(w, err.Error(), http.StatusRequestEntityTooLarge)
			return
		}
		defer r.Body.Close()

		log.Debug("Uploading mods' settings save data for account of ID %v...", user.Account)

		modsRes := access.R2WriteModsSettings(user.Account, body)
		if modsRes.IsError() {
			utils.WriteWebErr(w, modsRes.Error().Error(), http.StatusBadRequest)
			return
		}

		utils.WriteWebRes(w, mo.Some("Successfully uploaded mods' settings save data!"), http.StatusOK)
		log.Info("Successfully uploaded mods' settings save data for account of ID %v", user.Account)
	})

	// coming soon, just lazy rn
	// http.HandleFunc("/api/v1/upload-geode-mods-saves", func(w http.ResponseWriter, r *http.Request) { // all mod save data (to be supporter-only)
	// 	header := w.Header()
	// 	utils.WriteHeaders(&header, http.MethodPost, false)
	// 	http.NotFound(w, r)
	// })
}
