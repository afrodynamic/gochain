package adapter

type Registry struct {
	adapters map[string]ChainAdapter
}

func NewRegistry() *Registry {
	return &Registry{
		adapters: make(map[string]ChainAdapter),
	}
}

func (r *Registry) Register(id string, adapter ChainAdapter) {
	r.adapters[id] = adapter
}

func (r *Registry) Get(id string) (ChainAdapter, bool) {
	adapter, exists := r.adapters[id]

	return adapter, exists
}
