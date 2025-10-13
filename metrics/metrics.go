package metrics

import (
	"fmt"
	"sync"
	"time"
)

// metrics routines
// init: for each request, track latency and give overall successs rate

type MetricStruct struct {
	Identifier uint8
	StartTime  time.Time
	EndTIme    time.Time
}

type Metrics struct {
	lock          sync.Mutex
	timemap       map[uint8]*MetricStruct
	eChan         chan *MetricStruct
	avglatency    time.Duration
	requestCount  int64
	responseCount int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		timemap:       make(map[uint8]*MetricStruct),
		eChan:         make(chan *MetricStruct, 1000),
		requestCount:  0,
		responseCount: 0,
	}
}

func (m *Metrics) Start(identifier uint8) {
	m.lock.Lock()
	m.requestCount += 1
	m.timemap[identifier] = &MetricStruct{Identifier: identifier, StartTime: time.Now()}
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
	m.lock.Lock()
	defer m.lock.Unlock()
	return len(m.timemap) == 0

}

// consumer
func (m *Metrics) StartCollection() {
	var tlatency time.Duration
	var rescount int64
	for mstruct := range m.eChan {
		latency := mstruct.EndTIme.Sub(mstruct.StartTime)
		rescount++
		tlatency += latency
		alatency := tlatency / time.Duration(rescount)
		m.avglatency = alatency
		m.responseCount = rescount // modify this so that only proper responses are counted
		fmt.Printf("[METRICS]: ID: %d - Latency: %d | Avg Latency: %d\n", mstruct.Identifier, latency.Milliseconds(), alatency.Milliseconds())
	}

}

func (m *Metrics) GetSummary() {
	// average latency
	// average error rate
	fmt.Printf("\n===[METRICS SUMMARY]===\nRequest Count\t\t: %d\nResponse Count\t\t: %d\nAvg Latency\t\t: %dms\n", m.requestCount, m.responseCount, m.avglatency.Milliseconds())
}
