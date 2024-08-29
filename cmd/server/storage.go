package main

type Gauge float64
type Counter int64

type MemStorage struct {
	Gauges   map[string]Gauge
	Counters map[string]Counter
}

func NewMemStorage() MemStorage {
	var storage MemStorage
	storage.Gauges = make(map[string]Gauge, 0)
	storage.Counters = make(map[string]Counter, 0)
	return storage
}

func (st *MemStorage) UpdateGauge(name string, val float64) {
	st.Gauges[name] = Gauge(val)
}

func (st *MemStorage) UpdateCounter(name string, val int64) {
	st.Counters[name] += Counter(val)
}
