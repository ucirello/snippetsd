package web // import "cirello.io/snippetsd/pkg/ui/web"

import (
	"encoding/json"
	"log"
	"net/http"

	"cirello.io/snippetsd/pkg/infra/repositories"
	"cirello.io/snippetsd/pkg/models/user"
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
	s.mux.HandleFunc("/snippetsByUser", s.snippetsByUser)
	s.mux.HandleFunc("/state", s.state)
	s.mux.HandleFunc("/", http.NotFound)
}

func (s *Server) unauthorized(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", `Basic realm="snippetsd"`)
	w.WriteHeader(http.StatusUnauthorized)
}

// ServeHTTP process HTTP requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: handle Access-Control-Allow-Origin correctly
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5200")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	email, password, ok := r.BasicAuth()
	if !ok {
		s.unauthorized(w)
		return
	}

	u, err := repositories.Users(s.db).GetByEmail(email)
	if err != nil {
		s.unauthorized(w)
		return
	}

	if err := user.Authenticate(u, email, password); err != nil {
		s.unauthorized(w)
		return
	}

	r = r.WithContext(user.WithContext(r.Context(), u))
	s.mux.ServeHTTP(w, r)
}

func (s *Server) state(w http.ResponseWriter, r *http.Request) {
	snippets, err := repositories.Snippets(s.db).All()
	if err != nil {
		log.Println("cannot load all snippets:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(snippets); err != nil {
		log.Println("cannot marshal snippets:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}

func (s *Server) snippetsByUser(w http.ResponseWriter, r *http.Request) {
	user := user.WhoAmI(r.Context())
	snippets, err := repositories.Snippets(s.db).GetByUser(user)
	if err != nil {
		log.Println("cannot load user's snippets:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(snippets); err != nil {
		log.Println("cannot marshal snippets:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}
