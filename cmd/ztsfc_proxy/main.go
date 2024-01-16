package main

import (
	"github.com/leobrada/ztsfc_proxy/internal/logger"
)

func init() {
	logger.Logger_init()
}

func main() {
	logger.Log.Info("Hello...")
}

/*
func main() {

	server := &http.Server{
		Addr:     "134.60.77.40:8081",
		ErrorLog: log.New(logger.Out, "", 0),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Info("Request received")

		// Respond with "hello world"
		fmt.Fprint(w, "hello world")
	})

	if err := server.ListenAndServe(); err != nil {
		logger.Error(err)
	}
}
*/
