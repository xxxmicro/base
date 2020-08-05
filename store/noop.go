package store

type noopStore struct{}

func (n *noopStore) Init(opts ...Option) error {
	return nil
}

func (n *noopStore) Options() Options {
	return Options{}
}

func (n *noopStore) String() string {
	return "noop"
}

func (n *noopStore) Get(key string, opts ...ReadOption) (*Record, error) {
	return nil, nil
}

func (n *noopStore) Set(r *Record, opts ...WriteOption) error {
	return nil
}

func (n *noopStore) Delete(key string, opts ...DeleteOption) error {
	return nil
}

func (n *noopStore) List(opts ...ListOption) ([]string, error) {
	return []string{}, nil
}

func (n *noopStore) Close() error {
	return nil
}