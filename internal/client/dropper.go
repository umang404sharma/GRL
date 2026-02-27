package client

import (
	"math/rand"
	"sync/atomic"
)

type Dropper struct {
	dropRatio atomic.Value
}

func NewDropper() *Dropper {
	d := &Dropper{}
	d.dropRatio.Store(float64(0))
	return d
}

func (d *Dropper) SetRatio(r float64) {
	if r < 0 {
		r = 0
	}
	if r > 1 {
		r = 1
	}

	d.dropRatio.Store(r)
}

func (d *Dropper) GetRatio() float64 {
	return d.dropRatio.Load().(float64)
}

func (d *Dropper) ShouldDrop() bool {
	ratio := d.GetRatio()
	return rand.Float64() < ratio
}