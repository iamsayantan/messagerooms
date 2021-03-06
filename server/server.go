package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/iamsayantan/messagerooms/room"
	"github.com/iamsayantan/messagerooms/user"
)

// maximum bytes allowed in request body. 1MB
var maxAllowedLimit int64 = 1048576

// Server holds the dependencies for handling all the interactions.
type Server struct {
	User user.Service
	Room room.Service

	Hub    *SSEHub
	router chi.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received HTTP Request: %s", r.URL)
	s.router.ServeHTTP(w, r)
}

// NewServer returns a new HTTP server.
func NewServer(us user.Service, rs room.Service, hub *SSEHub) *Server {
	s := &Server{
		User: us,
		Room: rs,
		Hub:  hub,
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	})
	am := newAuthMiddleware(us)

	r := chi.NewRouter()
	r.Use(chiware.AllowContentType("application/json"))
	r.Use(corsHandler.Handler)

	r.Method("GET", "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Message string `json:"message"`
		}{Message: "Welcome to Message Rooms"}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			_, _ = w.Write([]byte(err.Error()))
		}
	}))

	r.Route("/sse", func(r chi.Router) {
		// Authentication middleware
		r.Use(am.Register)
		r.Get("/connect", http.HandlerFunc(s.Hub.HandleSSE))
	})

	r.Route("/user", func(r chi.Router) {
		h := NewUserHandler(us, am)
		r.Mount("/v1", h.Route())
	})

	r.Route("/rooms", func(r chi.Router) {
		r.Use(am.Register)
		h := newRoomHandler(rs)
		r.Mount("/v1", h.Route())
	})

	r.Method("GET", "/metrics", promhttp.Handler())

	s.router = r

	return s
}

// WebHandler interface provides general interface for web request handlers.
type WebHandler interface {
	Route() chi.Router
}

func sendResponse(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp, _ := json.Marshal(v)

	w.WriteHeader(statusCode)
	_, _ = w.Write(resp)
}

// malformedRequest represents errors generated while parsing JSON requests to go structs.
type malformedRequest struct {
	HTTPStatus int
	Message    string
}

func (mr *malformedRequest) Error() string {
	return mr.Message
}

// decodeJSONBody decodes an incoming request body and returns any error that occurs while parsing the request with the given struct.
func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxAllowedLimit)
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{HTTPStatus: http.StatusBadRequest, Message: msg}
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{HTTPStatus: http.StatusBadRequest, Message: msg}
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{HTTPStatus: http.StatusBadRequest, Message: msg}
		case errors.Is(err, io.EOF):
			msg := fmt.Sprintf("Request body must not be empty")
			return &malformedRequest{HTTPStatus: http.StatusBadRequest, Message: msg}
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{HTTPStatus: http.StatusRequestEntityTooLarge, Message: msg}
		default:
			return err
		}
	}
	return nil
}

// ErrResponse  renderer type is for rendering all sorts of errors.
type ErrResponse struct {
	Err            error  `json:"-"`      // low-level runtime error
	HTTPStatusCode int    `json:"-"`      // http response status code
	StatusText     string `json:"status"` // user level status message
	ErrorText      string `json:"error"`  // application level error message
}

// Render set the http status code before the response marshalling.
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest returns error response struct with appropiate status and message.
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Bad Request",
		ErrorText:      err.Error(),
	}
}

// ErrNotFound returns error response with appropiate status.
func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "Not Found",
		ErrorText:      err.Error(),
	}
}

// ErrUnAuthorized returns error response with appropiate status.
func ErrUnAuthorized(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Authentication Required",
		ErrorText:      err.Error(),
	}
}

// ErrInternalServer returns error response with appropiate status.
func ErrInternalServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
	}
}
