package mock

// MockJob is a fake Gearman job for tests
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

// Err returns the job's error
func (m MockJob) Err() error {
	return m.err
}

// Data returns the job's data
func (m MockJob) Data() []byte {
	return []byte(m.payload)
}

// Fn returns the job's name
func (m MockJob) Fn() string {
	return m.name
}

// Handle returns the job's handle
func (m MockJob) Handle() string {
	return m.handle
}

// UniqueId returns the unique id for the job
func (m MockJob) UniqueId() string {
	return m.id
}

// SendWarning appends to the array of job warnings
func (m *MockJob) SendWarning(warning []byte) {
	m.warnings = append(m.warnings, warning)
}

// Warnings returns the array of jobs warnings
func (m *MockJob) Warnings() [][]byte {
	return m.warnings
}

// SendData appends to the array of job data
func (m *MockJob) SendData(data []byte) {
	m.data = append(m.data, data)
}

// UpdateStatus updates the progress of job
func (m *MockJob) UpdateStatus(numerator, denominator int) {
	m.numerator = numerator
	m.denominator = denominator
}
