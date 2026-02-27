package models

import "time"

type Directive struct {
	DropRatio float64   `json:"drop_ratio"`
	TTL       int       `json:"ttl"`
	Version   int64     `json:"version"`
	ExpiresAt time.Time `json:"-"`
}
