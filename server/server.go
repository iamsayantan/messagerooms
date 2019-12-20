package server

import "github.com/iamsayantan/MessageRooms/user"
import "github.com/iamsayantan/MessageRooms/room"
import "github.com/go-chi/chi"
import chiware "github.com/go-chi/chi/middleware"
import "net/http"
import "encoding/json"

// Server holds the dependencies for handling all the interactions.
type Server struct {
	User user.Service
	Room room.Service

	router chi.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// NewServer returns a new HTTP server.
func NewServer(us user.Service, rs room.Service) *Server {
	s := &Server{
		User: us,
		Room: rs,
	}

	r := chi.NewRouter()
	r.Use(chiware.AllowContentType("application/json"))
	r.Use(ServeJson)

	r.Method("GET", "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Message string `json:"message"`
		}{Message: "Welcome to Message Rooms"}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			_, _ = w.Write([]byte(err.Error()))
		}
	}))

	r.Route("/user", func(r chi.Router) {
		h := NewUserHandler(us)
		r.Mount("/v1", h.Route())
	})

	s.router = r

	return s
}

// WebHandler interface provides general interface for web request handlers.
type WebHandler interface {
	Route() chi.Router
}

func encodeError(w http.ResponseWriter, errorCode int, message string) {
	err := struct {
		Error string `json:"error"`
	}{Error: message}

	resp, _ := json.Marshal(err)
	w.WriteHeader(errorCode)
	_, _ = w.Write(resp)
}

func sendResponse(w http.ResponseWriter, statusCode int, v interface{}) {
	resp, _ := json.Marshal(v)

	w.WriteHeader(statusCode)
	_, _ = w.Write(resp)
}

func ServeJson(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Powerd-By", "golang")
		h.ServeHTTP(w, r)
	})
}
