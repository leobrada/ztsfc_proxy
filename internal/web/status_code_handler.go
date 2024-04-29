package web

import (
	"fmt"
	"net/http"
)

func Handle404(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	responseMessage := "<html><body><h1>404 Not Found</h1><p>The requested resource could not be found.</p></body></html>"
	fmt.Fprint(w, responseMessage)
}

func Handle500(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	responseMessage := "<html><body><h1>500 Internal Server Error</h1><p>Sorry, something went wrong on our end.</p></body></html>"
	fmt.Fprint(w, responseMessage)
}

func Handle501(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotImplemented)
	responseMessage := "<html><body><h1>501 Not Implemented</h1><p>Sorry, the requested functionality is not supported.</p></body></html>"
	fmt.Fprint(w, responseMessage)
}
