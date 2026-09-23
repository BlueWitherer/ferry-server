package api

import (
	"fmt"
	"net/http"

	"github.com/BlueWitherer/ferry-server/log"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		log.Debug("API server pinged!")

		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", http.MethodGet)
		header.Set("Access-Control-Allow-Headers", "Content-Type")
		header.Set("Content-Type", "text/plain")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Game Data Sync master API service")
	})
}
