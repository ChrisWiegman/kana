package mocks

type ReadCloser struct {
	expectedData []byte
	expectedErr  error
}

func (m *ReadCloser) Read(p []byte) (n int, err error) {
	copy(p, m.expectedData)
	return 0, m.expectedErr
}

func (m *ReadCloser) Close() error {
	return nil
}
