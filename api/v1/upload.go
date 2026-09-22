package v1

import (
	"net/http"

	"github.com/BlueWitherer/GDDataSyncServer/utils"
)

func init() {
	http.HandleFunc("/v1/upload", func(w http.ResponseWriter, r *http.Request) { // geometry dash gamevars
		header := w.Header()
		utils.WriteHeaders(&header, http.MethodPost)
	})

	http.HandleFunc("/v1/upload-mods", func(w http.ResponseWriter, r *http.Request) { // all mod settings
		header := w.Header()
		utils.WriteHeaders(&header, http.MethodPost)
	})

	http.HandleFunc("/v1/upload-mods-saves", func(w http.ResponseWriter, r *http.Request) { // all mod save data (supporter only probably)
		header := w.Header()
		utils.WriteHeaders(&header, http.MethodPost)
	})
}
