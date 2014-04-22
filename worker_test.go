package gearman

import (
	"fmt"
	"io"
	"net"
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

func makeTCPServer(addr string, handler func(conn net.Conn) error) chan error {
	channel := make(chan error)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		if err := handler(conn); err != nil {
			channel <- err
		}
	}()

	return channel
}

func readBytes(reader io.Reader, size uint) ([]byte, error) {
	buf := make([]byte, size)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func readGearmanHeader(reader io.Reader) (uint, uint, error) {
	header, err := readBytes(reader, 12)
	if err != nil {
		return 0, 0, err
	}
	cmd := (uint(header[4]) << 24) | (uint(header[5]) << 16) |
		(uint(header[6]) << 8) | uint(header[7])
	cmdLen := (uint(header[8]) << 24) | (uint(header[9]) << 16) |
		(uint(header[10]) << 8) | uint(header[11])
	return cmd, cmdLen, nil
}

func readGearmanCommand(reader io.Reader) (uint, string, error) {
	cmd, dataSize, err := readGearmanHeader(reader)
	if err != nil {
		return 0, "", err
	}
	data, err := readBytes(reader, dataSize)
	if err != nil {
		return 0, "", err
	}
	return cmd, string(data), nil
}

func TestCanDo(t *testing.T) {

	var channel chan error

	name := "worker_name"

	channel = makeTCPServer(":1337", func(conn net.Conn) error {
		cmd, data, err := readGearmanCommand(conn)
		if err != nil {
			return err
		}
		// 1 = CAN_DO
		if cmd != 1 {
			return fmt.Errorf("expected command 1 (CAN_DO), received command %d", cmd)
		}
		if data != "worker_name" {
			return fmt.Errorf("expected '%s', received '%s'", name, data)
		}
		close(channel)
		return nil
	})

	worker := New(name, func(job Job) ([]byte, error) {
		return []byte{}, nil
	})
	go worker.Listen("localhost", "1337")

	for err := range channel {
		t.Fatal(err)
	}
}
