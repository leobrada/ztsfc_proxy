package hashutil

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"time"
)

// CalcRequestHash calculates the SHA256 hash of a given http.Request including the current time and returns the hash as string
func CalcRequestHash(r *http.Request) string {
	// Get the request method (GET, POST, etc.)
	method := r.Method

	// Get the request URL
	url := r.URL.String()

	// Get the headers as a string
	headers := ""
	for name, values := range r.Header {
		headers += name + ": " + strings.Join(values, ",") + "\n"
	}

	// Get the body of the request
	body, _ := io.ReadAll(r.Body)
	// Reset the body so it can be read again
	r.Body = io.NopCloser(strings.NewReader(string(body)))

	// Get the current time
	currentTime := time.Now().String()

	// Concatenate all parts to a single string
	toHash := method + r.Host + r.RemoteAddr + url + headers + string(body) + currentTime

	// Calculate the SHA256 hash
	hash := sha256.Sum256([]byte(toHash))

	// Return the first 8 byte of the hash as a hexadecimal string
	return hex.EncodeToString(hash[:7])
}
