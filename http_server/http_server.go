package http_server

import (
	"fmt"
	"github.com/prusiksoftware/monorepo/obscurio/analytics"
	"net/http"
)

type Server struct {
	isReady          bool
	isLive           bool
	port             int
	analyticsTracker *analytics.Analytics
}

func (s *Server) healthzReady(w http.ResponseWriter, r *http.Request) {
	if !s.isReady {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Not ready"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) healthzLive(w http.ResponseWriter, r *http.Request) {
	if !s.isLive {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Not live"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) debug(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if s.analyticsTracker != nil {
		fmt.Fprintf(w, "Ready: %t\n", s.isReady)
		fmt.Fprintf(w, "Live: %t\n", s.isLive)
		fmt.Fprintf(w, "Port: %d\n\n", s.port)
		for _, event := range s.analyticsTracker.Events {
			fmt.Fprintf(w, "\n\n")
			fmt.Fprintf(w, "original: %s\n", event.OriginalQuery)
			fmt.Fprintf(w, "created: %s\n", event.Created)
			fmt.Fprintf(w, "modified: %s\n", event.ModifiedQuery)
			fmt.Fprintf(w, "profile: %s\n", event.Profile)
			fmt.Fprintf(w, "events:\n")
			for k, v := range event.Durations {
				fmt.Fprintf(w, "  - %s: %s\n", k, v)
			}
		}
	}
}

func (s *Server) SetReady(ready bool) {
	s.isReady = ready
}

func (s *Server) SetLive(live bool) {
	s.isLive = live
}

func (s *Server) Start() {
	http.HandleFunc("/healthz/ready", s.healthzReady)
	http.HandleFunc("/healthz/live", s.healthzLive)
	http.HandleFunc("/debug", s.debug)
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("HTTP server listening on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("HTTP server stopped\n")
}

func NewHTTPServer(at *analytics.Analytics) *Server {
	return &Server{
		isReady:          false,
		isLive:           false,
		port:             80,
		analyticsTracker: at,
	}
}
