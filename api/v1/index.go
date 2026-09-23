package v1

import (
	"net/http"

	"github.com/BlueWitherer/ferry-server/log"
	"github.com/BlueWitherer/ferry-server/utils"

	"github.com/samber/mo"
)

func init() {
	http.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("v1 API pinged!")

		header := w.Header()
		utils.WriteHeaders(&header, http.MethodGet, false)

		utils.WriteWebRes(w, mo.Some("Pong!"), http.StatusOK)
	})
}
