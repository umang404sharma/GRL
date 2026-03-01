package aggregator

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"umang404sharma/GRL/internal/models"
)
type Server struct {
	zone       string
	controller string
	collector  *Collector
}

func NewServer(zone, controller string) *Server {
	return &Server {
		zone: zone,
		controller: controller,
		collector: NewCollector(),
	}
}

func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	var report models.HostReport

	err := json.NewDecoder(r.Body).Decode(&report)

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	s.collector.Update(report.Host, report.RPS)

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleDirective(w http.ResponseWriter, r *http.Request) {
	var directive models.Directive

	err := json.NewDecoder(r.Body).Decode(&directive)

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	log.Println("received directive:", directive.DropRatio)
	
	hosts := s.collector.Hosts()
	log.Println("hostss:", hosts)

	for _, host := range hosts {
		url := "http://" + host + "/directive"

		body, _ := json.Marshal(directive)

		_, err := http.Post(url, "application/json", bytes.NewBuffer(body))

		if err != nil {
			log.Println("failed to forward to", host)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) Start(port string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/report", s.handleReport)
	mux.HandleFunc("/directive", s.handleDirective)

	s.startReporter()

	log.Println("aggregator running on", port)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}