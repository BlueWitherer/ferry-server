package v1

import (
	"net/http"

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

		userRes, code := authUser(r)
		if userRes.IsError() {
			utils.WriteWebErr(w, userRes.Error().Error(), code)
			return
		}
		user := userRes.MustGet()

		gvRes := access.R2ReadGameVarsRaw(user.Account)
		if gvRes.IsError() {
			utils.WriteWebErr(w, gvRes.Error().Error(), http.StatusBadRequest)
			return
		}

		bytes := gvRes.MustGet()

		log.Debug("Streaming game settings save data for account of ID %v...", user.Account)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(bytes); err != nil {
			utils.WriteWebErr(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info("Successfully streamed game settings save data for account of ID %v", user.Account)
	})

	http.HandleFunc("/api/v1/download-geode", func(w http.ResponseWriter, r *http.Request) { // geode settings
		header := w.Header()
		utils.WriteHeaders(&header, http.MethodGet, true)

		if r.Method != http.MethodGet {
			utils.WriteWebErrMethod(w)
			return
		}

		userRes, code := authUser(r)
		if userRes.IsError() {
			utils.WriteWebErr(w, userRes.Error().Error(), code)
			return
		}
		user := userRes.MustGet()

		geodeRes := access.R2ReadGeodeSettings(user.Account)
		if geodeRes.IsError() {
			utils.WriteWebErr(w, geodeRes.Error().Error(), http.StatusBadRequest)
			return
		}

		bytes := geodeRes.MustGet()

		log.Debug("Streaming Geode loader settings save data for account of ID %v...", user.Account)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(bytes); err != nil {
			utils.WriteWebErr(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info("Successfully streamed Geode loader settings save data for account of ID %v", user.Account)
	})

	http.HandleFunc("/api/v1/download-geode-mods", func(w http.ResponseWriter, r *http.Request) { // geode mods
		header := w.Header()
		utils.WriteHeaders(&header, http.MethodGet, true)

		if r.Method != http.MethodGet {
			utils.WriteWebErrMethod(w)
			return
		}

		userRes, code := authUser(r)
		if userRes.IsError() {
			utils.WriteWebErr(w, userRes.Error().Error(), code)
			return
		}
		user := userRes.MustGet()

		modsRes := access.R2ReadModsSettings(user.Account)
		if modsRes.IsError() {
			utils.WriteWebErr(w, modsRes.Error().Error(), http.StatusBadRequest)
			return
		}

		bytes := modsRes.MustGet()

		log.Debug("Streaming mods' settings save data for account of ID %v...", user.Account)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(bytes); err != nil {
			utils.WriteWebErr(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info("Successfully streamed mods' settings save data for account of ID %v", user.Account)
	})

	// coming soon, just lazy rn
	// http.HandleFunc("/api/v1/download-geode-mods-saves", func(w http.ResponseWriter, r *http.Request) { // all mod save data (to be supporter-only)
	// 	header := w.Header()
	// 	utils.WriteHeaders(&header, http.MethodPost, false)
	// 	http.NotFound(w, r)
	// })
}
