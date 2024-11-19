package handler

type Handler struct {
	// handler is struct with 3 pointers to clients.
	// this is how we bring together all the data.
}

type HandlerOptions func(*Handler) error

// initialize new handler instance
func New(opts ...HandlerOptions) (*Handler, error) {
	h := &Handler{}

	// error checking for all options
	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, err
		}
	}

	println("Successfully created handler instance!")

	return h, nil
}
