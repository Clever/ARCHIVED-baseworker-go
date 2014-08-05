package mock

type MockJob struct {
	payload, name, handle, id string
	err                       error
	warnings                  [][]byte
	data                      [][]byte
	numerator, denominator    int
}

// CreateMockJob creates an object that implements the gearman-go/worker#Job interface
func CreateMockJob(payload string) *MockJob {
	return &MockJob{payload: payload}
}

func (m MockJob) Err() error {
	return m.err
}
func (m MockJob) Data() []byte {
	return []byte(m.payload)
}
func (m MockJob) Fn() string {
	return m.name
}
func (m MockJob) Handle() string {
	return m.handle
}
func (m MockJob) UniqueId() string {
	return m.id
}
func (m *MockJob) SendWarning(warning []byte) {
	m.warnings = append(m.warnings, warning)
}
func (m *MockJob) GetWarnings() [][]byte {
	return m.warnings
}
func (m *MockJob) SendData(data []byte) {
	m.data = append(m.data, data)
}
func (m *MockJob) UpdateStatus(numerator, denominator int) {
	m.numerator = numerator
	m.denominator = denominator
}
