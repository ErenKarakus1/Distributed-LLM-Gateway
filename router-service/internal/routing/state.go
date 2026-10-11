package routing

import "sync"

type RouterState struct {
	mu       sync.Mutex
	counters map[string]int
}

func NewRouterState() *RouterState {
	return &RouterState{
		counters: make(map[string]int),
	}
}

func (r *RouterState) OrderProviders(model string, providers []string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(providers) == 0 {
		return providers
	}

	index := r.counters[model] % len(providers)
	r.counters[model]++

	ordered := append([]string{}, providers[index:]...)
	ordered = append(ordered, providers[:index]...)

	return ordered
}
