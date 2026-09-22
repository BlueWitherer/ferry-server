package v1

import (
	"net/http"

	"github.com/BlueWitherer/GDDataSyncServer/log"
	"github.com/BlueWitherer/GDDataSyncServer/utils"
	"github.com/samber/mo"
)

func init() {
	http.HandleFunc("/v1", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("v1 API pinged!")

		header := w.Header()
		utils.WriteHeaders(&header, http.MethodGet)

		utils.WriteWebRes(w, mo.Some("Pong!"), http.StatusOK)
	})
}
