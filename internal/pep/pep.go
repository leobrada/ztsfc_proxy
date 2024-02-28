package pep

import (
	"fmt"
	"log"
	"net/http"

	"github.com/leobrada/ztsfc_proxy/internal/configs"
)

type PEP struct {
	dpLogger *log.Logger
}

func NewPEP(config *configs.Config, dataPlaneLogger *log.Logger) (*PEP, error) {
	return &PEP{dpLogger: dataPlaneLogger}, nil
}

func (pep *PEP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Setting content type to text/plain
	w.Header().Set("Content-Type", "text/plain")

	// Writing "Hello World" to the response
	_, err := fmt.Fprint(w, "Hello World")
	if err != nil {
		http.Error(w, "Error writing string", http.StatusInternalServerError)
	}
}
