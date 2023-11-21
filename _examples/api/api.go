// Example - demonstrates REST API server implementation tests.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/regen-network/gocuke"
)

// Server implements http.Server
type Server struct {
	*http.Server
}

// NewServer creates a new instance of Server
func NewServer(port int, handler http.Handler) *Server {
	addr := fmt.Sprintf(":%d", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return &Server{
		Server: srv,
	}
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fail(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data := struct {
		Version string `json:"version"`
	}{Version: gocuke.Version}
	ok(w, data)
}

func main() {
	// Define your handler
	mux := http.NewServeMux()
	mux.HandleFunc("/version", getVersion)

	// Create a new server instance
	server := NewServer(8080, mux)

	// Start the server
	fmt.Printf("Server running on port %s...\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
