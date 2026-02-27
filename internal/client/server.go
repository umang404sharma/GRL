package client

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
	"time"
	"umang404sharma/GRL/internal/models"
)

type Server struct {
	hostID     string
	dropper    *Dropper
	reqCount   atomic.Int64
	aggregator string
	directive  atomic.Pointer[models.Directive]
}

func NewServer(hostID, aggregator string) *Server {
	s := &Server{
		hostID:     hostID,
		dropper:    NewDropper(),
		aggregator: aggregator,
	}
	return s
}

func (s *Server) handleAPI(w http.ResponseWriter, r *http.Request) {
	s.reqCount.Add(1)

	if s.dropper.ShouldDrop() {
		http.Error(w, "rate limited", http.StatusTooManyRequests)
		return
	}

	time.Sleep(10 * time.Millisecond)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) handleDirective(w http.ResponseWriter, r *http.Request) {
	var d models.Directive

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	d.ExpiresAt = time.Now().Add(time.Duration(d.TTL) * time.Second)

	s.directive.Store(&d)
	s.dropper.SetRatio(d.DropRatio)

	w.WriteHeader(http.StatusOK)
}

func (s *Server) startDirectiveWatcher() {
	go func() {
		for {
			time.Sleep(1 * time.Second)

			v := s.directive.Load()
			if v == nil {
				continue
			}

			if time.Now().After(v.ExpiresAt) {
				log.Println("directive expired, failing open")

				s.dropper.SetRatio(0)
				s.directive.Store(nil)
			}
		}
	}()
}

func (s *Server) startReporter() {
	go func() {
		ticker := time.NewTicker((1 * time.Second))
		for range ticker.C {
			count := s.reqCount.Swap(0)

			report := models.HostReport{
				Host: s.hostID,
				RPS:  count,
			}

			body, _ := json.Marshal(report)
			http.Post(s.aggregator+"/report", "application/json", bytes.NewBuffer(body))
		}
	}()
}

func (s *Server) Start(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api", s.handleAPI)
	mux.HandleFunc("/directive", s.handleDirective)

	s.startReporter()
	s.startDirectiveWatcher()

	log.Println("client running on", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
