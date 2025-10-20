package core

type Registry struct {
	adapters map[string]ChainAdapter
}

func NewRegistry(adapters ...ChainAdapter) *Registry {
	m := make(map[string]ChainAdapter, len(adapters))
	for _, a := range adapters {
		m[a.Network()] = a
	}
	return &Registry{adapters: m}
}

func (r *Registry) Get(name string) (ChainAdapter, bool) {
	a, ok := r.adapters[name]
	return a, ok
}

func (r *Registry) List() []string {
	out := make([]string, 0, len(r.adapters))
	for k := range r.adapters {
		out = append(out, k)
	}
	return out
}
