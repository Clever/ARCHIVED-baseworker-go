package gearman

import (
	"testing"
)

type MockJob struct {
	payload, name, handle, id string
	err                       error
	warnings                  [][]byte
	data                      [][]byte
	numerator, denominator    int
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
func (m *MockJob) SendData(data []byte) {
	m.data = append(m.data, data)
}
func (m *MockJob) UpdateStatus(numerator, denominator int) {
	m.numerator = numerator
	m.denominator = denominator
}

func TestJobFuncConversion(t *testing.T) {
	payload := "I'm a payload!"
	jobFunc := func(job Job) ([]byte, error) {
		if string(job.Data()) != payload {
			t.Fatalf("expected payload %s, received %s", payload, string(job.Data()))
		}
		return []byte{}, nil
	}
	worker := New("test", jobFunc)
	worker.fn(&MockJob{payload: payload})
}
