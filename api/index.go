package api

import (
	"net/http"

	"ferry-srv/log"
	"ferry-srv/utils"

	"github.com/samber/mo"

	_ "ferry-srv/api/v1"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		log.Trace("API server pinged!")

		header := w.Header()
		utils.WriteHeaders(&header, http.MethodGet, false)
		utils.WriteWebRes(w, mo.Some("Ferry master API service"), http.StatusOK)
	})
}
