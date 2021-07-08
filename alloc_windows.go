package stealthpool

func alloc(size int) ([]byte, error) {
	return make([]byte, size), nil
}

func dealloc(b []byte) error {
	return nil
}
