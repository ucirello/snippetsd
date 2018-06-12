package web // import "cirello.io/snippetsd/pkg/ui/web"

import (
	"encoding/json"
	"log"
	"net/http"

	"cirello.io/snippetsd/pkg/models/snippet"
	"github.com/jmoiron/sqlx"
)

// Server implements the web interface.
type Server struct {
	db  *sqlx.DB
	mux *http.ServeMux
}

// New creates a web interface handler.
func New(db *sqlx.DB) *Server {
	s := &Server{
		db:  db,
		mux: http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/state", s.state)
	s.mux.HandleFunc("/", http.NotFound)
}

// ServeHTTP process HTTP requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) state(w http.ResponseWriter, r *http.Request) {
	// TODO: handle Access-Control-Allow-Origin correctly
	w.Header().Set("Access-Control-Allow-Origin", "*")
	snippets, err := snippet.LoadAll(s.db)
	if err != nil {
		log.Println("cannot load all snippets:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(snippets); err != nil {
		log.Println("cannot marshal snippets:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}
