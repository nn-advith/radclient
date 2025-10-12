package metrics

import (
	"fmt"
	"sync"
	"time"
)

// metrics routines
// init: for each request, track latency and give overall successs rate

type TimeStruct struct {
	Identifier uint8
	StartTime  time.Time
	EndTIme    time.Time
}

type Metrics struct {
	lock    sync.Mutex
	timemap map[uint8]*TimeStruct
	eChan   chan *TimeStruct
}

func NewMetrics() *Metrics {
	return &Metrics{
		timemap: make(map[uint8]*TimeStruct),
		eChan:   make(chan *TimeStruct, 50),
	}
}

func (m *Metrics) Start(identifier uint8) {
	m.lock.Lock()
	m.timemap[identifier] = &TimeStruct{Identifier: identifier, StartTime: time.Now()}
	m.lock.Unlock()
}

func (m *Metrics) End(identifier uint8) {
	m.lock.Lock()
	// check if id is present
	if val, exists := m.timemap[identifier]; exists {
		// set end time
		val.EndTIme = time.Now()
		delete(m.timemap, identifier)
		// producer
		m.eChan <- val
	}
	m.lock.Unlock()
}

func (m *Metrics) IsEmpty() bool {
	return len(m.eChan) == 0
}

// consumer
func (m *Metrics) StartCollection() {
	for timestruct := range m.eChan {
		latency := timestruct.EndTIme.Sub(timestruct.StartTime)
		fmt.Printf("[METRICS]: ID: %d - Latency: %d\n", timestruct.Identifier, latency.Milliseconds())
	}
}
