package v1

import (
	"encoding/json"
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

		log.Trace(
			"Game settings upload for %v size: %d bytes (%.2f KiB)",
			user.Account,
			len(body),
			float64(len(body))/1024,
		)

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

		log.Trace(
			"Geode settings upload for %v size: %d bytes (%.2f KiB)",
			user.Account,
			len(body),
			float64(len(body))/1024,
		)

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

		body, err := io.ReadAll(io.LimitReader(r.Body, 64<<10))
		if err != nil {
			utils.WriteWebErr(w, err.Error(), http.StatusRequestEntityTooLarge)
			return
		}
		defer r.Body.Close()

		log.Trace(
			"Mod settings upload for %v size: %d bytes (%.2f KiB)",
			user.Account,
			len(body),
			float64(len(body))/1024,
		)

		log.Debug("Uploading mods' settings save data for account of ID %v...", user.Account)

		jsonRes := access.ParseModsSettingsBody(body)
		if jsonRes.IsError() {
			utils.WriteWebErr(w, jsonRes.Error().Error(), http.StatusBadRequest)
			return
		}

		existing := access.R2ReadModsSettings(user.Account)
		if existing.IsOk() {
			log.Debug("Found existing mod settings save for user %v", user.Account)

			saved := jsonRes.MustGet()
			newSaveRes := access.ParseModsSettingsBody(existing.MustGet())

			if newSaveRes.IsOk() {
				for key, val := range newSaveRes.MustGet() {
					_, ok := saved[key]
					if !ok {
						log.Trace("Adding missing '%s' mod settings for user %v", key, user.Account)
						saved[key] = val
					}
				}
			}

			b, err := json.Marshal(saved)
			if err != nil {
				utils.WriteWebErr(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Debug("Merging settings for user %v", user.Account)

			size := len(b)
			if size < 64<<10 {
				log.Debug("Mod settings save is within limits (%.2f KiB)", float64(size)/1024)
				body = b
			} else {
				log.Warn("Mod settings save exceeds 64 KiB limit (%.2f KiB) for user %v, skipping merge", float64(size)/1024, user.Account)
			}
		}

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
