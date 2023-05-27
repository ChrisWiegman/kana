package mocks

type ReadCloser struct {
	ExpectedData []byte
	ExpectedErr  error
}

func (m *ReadCloser) Read(p []byte) (n int, err error) {
	copy(p, m.ExpectedData)
	return 0, m.ExpectedErr
}

func (m *ReadCloser) Close() error {
	return nil
}
