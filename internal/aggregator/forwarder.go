package aggregator

import (
	"bytes"
	"encoding/json"
	"log"
	"time"
	"umang404sharma/GRL/internal/models"
)

func (s *Server) startReporter() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)

		for range ticker.C {
			total := s.collector.Total()

			report := models.ZoneReport{
				Zone:     s.zone,
				TotalRPS: total,
			}

			body, _ := json.Marshal(report)
			log.Println("REPORT BODY", bytes.NewBuffer(body))
			// _, err := http.Post(
			// 	s.controller+"/zone-report",
			// 	"application/json",
			// 	bytes.NewBuffer(body),
			// )

			// if err != nil {
			// 	log.Println("failed to report to controller:", err)
			// } else {
			// 	log.Println("zone total RPS:", total)
			// }

		}
	}()
}
