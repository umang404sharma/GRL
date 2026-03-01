package aggregator

import "sync"

type Collector struct {
	mu      sync.Mutex
	metrics map[string]int64
}

func NewCollector() *Collector {
	return &Collector{
		metrics: make(map[string]int64),
	}
}

func (c *Collector) Update(host string, rps int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics[host] = rps
}

func (c *Collector) Total() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	var total int64
	for _, rps := range c.metrics {
		total += rps
	}

	return total
}

func (c *Collector) Hosts() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	hosts := make([]string, 0, len(c.metrics))
	for host := range c.metrics {
		hosts = append(hosts, host)
	}

	return hosts
}