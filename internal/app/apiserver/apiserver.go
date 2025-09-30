package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Manabreaker/Calendar/internal/app/store"
	"github.com/Manabreaker/Calendar/internal/app/store/eventStore"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	ctxKeyRequestID = iota
)

type Server struct {
	router *mux.Router
	config *ServerConfig
	logger *zap.Logger
	store  store.Store
}

func NewServer() *Server {
	logger, _ := zap.NewProduction()
	s := &Server{
		router: mux.NewRouter(),
		logger: logger,
		store:  eventStore.NewStore(), // *EventStore реализует store.Store
	}
	s.configureRouter()
	return s
}

func (s *Server) Start(config *ServerConfig) error {
	s.config = config
	return http.ListenAndServe(config.Host+":"+config.Port, s.router)
}

func (s *Server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)

	s.router.HandleFunc("/create_event", s.createEvent).Methods(http.MethodPost)
	s.router.HandleFunc("/update_event", s.updateEvent).Methods(http.MethodPost)
	s.router.HandleFunc("/delete_event", s.deleteEvent).Methods(http.MethodPost)

	s.router.HandleFunc("/events_for_day", s.GetEventsToday).Methods(http.MethodGet)
	s.router.HandleFunc("/events_for_week", s.GetEventsWeek).Methods(http.MethodGet)
	s.router.HandleFunc("/events_for_month", s.GetEventsMonth).Methods(http.MethodGet)
	s.router.HandleFunc("/events_for_all", s.GetAllEvents).Methods(http.MethodGet)
}

func (s *Server) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events := s.store.GetAllEvents()
	s.respond(w, http.StatusOK, events)

}
func (s *Server) GetEventsMonth(w http.ResponseWriter, r *http.Request) {
	// Implementation for getting events for Month
	events := s.store.GetEventsMonth()
	s.respond(w, http.StatusOK, events)

}
func (s *Server) GetEventsWeek(w http.ResponseWriter, r *http.Request) {
	// Implementation for getting events for Week
	events := s.store.GetEventsWeek()
	s.respond(w, http.StatusOK, events)

}
func (s *Server) GetEventsToday(w http.ResponseWriter, r *http.Request) {
	// Implementation for getting events for today
	events := s.store.GetEventsToday()
	s.respond(w, http.StatusOK, events)

}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// setRequestID Middleware to set a unique request ID for each request
// It checks for an existing X-Request-ID header and uses it if present
// Otherwise, it generates a new UUID and sets it in the response header
func (s *Server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.With(
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("request_id", r.Context().Value(ctxKeyRequestID).(string)),
		)
		logger.Info(fmt.Sprintf("started %s %s", r.Method, r.RequestURI))
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		logger.Info(
			"completed",
			zap.Int("status_code", rw.code),
			zap.String("status_text", http.StatusText(rw.code)),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

func (s *Server) error(w http.ResponseWriter, code int, err error) {
	s.logger.Error("error", zap.Error(err))
	s.respond(w, code, map[string]string{"error": err.Error()})
}

func (s *Server) respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.error(w, http.StatusBadRequest, err)
		return
	}
	if err := s.store.Create(body); err != nil {
		s.error(w, http.StatusBadRequest, err)
		return
	}
	s.respond(w, http.StatusCreated, map[string]string{"result": "success"})
}

func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.error(w, http.StatusBadRequest, err)
		return
	}
	if err := s.store.Update(body); err != nil {
		s.error(w, http.StatusBadRequest, err)
		return
	}
	s.respond(w, http.StatusCreated, map[string]string{"result": "success"})
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if err := s.store.Delete(id); err != nil {
		s.error(w, http.StatusNotFound, err)
		return
	}
	s.respond(w, http.StatusCreated, map[string]string{"result": "success"})
}
