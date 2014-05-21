package baseworker

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/mikespook/gearman-go/client"
	"io"
	"log"
	"net"
	"os"
	"testing"
	"time"
)

var GearmanPort = os.Getenv("GEARMAN_PORT")
var GearmanHost = os.Getenv("GEARMAN_HOST")

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

// TestJobFuncConversion tests that our JobFunc is called when 'worker.fn' is called with a job.
func TestJobFuncConversion(t *testing.T) {
	payload := "I'm a payload!"
	jobFunc := func(job Job) ([]byte, error) {
		if string(job.Data()) != payload {
			t.Fatalf("expected payload %s, received %s", payload, string(job.Data()))
		}
		return []byte{}, nil
	}
	worker := NewWorker("test", jobFunc)
	worker.fn(&MockJob{payload: payload})
}

func makeTCPServer(addr string, handler func(conn net.Conn) error) (net.Listener, chan error) {
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

	return listener, channel
}

func readBytes(reader io.Reader, size uint32) ([]byte, error) {
	buf := make([]byte, size)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func fromBigEndianBytes(buf []byte) (uint32, error) {
	var num uint32
	if err := binary.Read(bytes.NewReader(buf), binary.BigEndian, &num); err != nil {
		return 0, err
	}
	return num, nil
}

func toBigEndianBytes(num uint32) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, num); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func readGearmanHeader(reader io.Reader) (uint32, uint32, error) {
	header, err := readBytes(reader, 12)
	if err != nil {
		return 0, 0, err
	}
	cmd, err := fromBigEndianBytes(header[4:8])
	if err != nil {
		return 0, 0, err
	}
	cmdLen, err := fromBigEndianBytes(header[8:12])
	if err != nil {
		return 0, 0, err
	}
	return cmd, cmdLen, nil
}

func readGearmanCommand(reader io.Reader) (uint32, string, error) {
	cmd, dataSize, err := readGearmanHeader(reader)
	if err != nil {
		return 0, "", err
	}
	body, err := readBytes(reader, dataSize)
	if err != nil {
		return 0, "", err
	}
	return cmd, string(body), nil
}

// TestCanDo tests that Listen properly sends a 'CAN_DO worker_name' packet to the TCP server.
func TestCanDo(t *testing.T) {

	var channel chan error
	var listener net.Listener

	name := "worker_name"

	listener, channel = makeTCPServer(":1338", func(conn net.Conn) error {
		cmd, body, err := readGearmanCommand(conn)
		if err != nil {
			return err
		}
		// 1 = CAN_DO
		if cmd != 1 {
			return fmt.Errorf("expected command 1 (CAN_DO), received command %d", cmd)
		}
		if body != "worker_name" {
			return fmt.Errorf("expected '%s', received '%s'", name, body)
		}
		close(channel)
		return nil
	})
	defer listener.Close()

	worker := NewWorker(name, func(job Job) ([]byte, error) {
		return []byte{}, nil
	})
	go worker.Listen("localhost", "1338")
	defer worker.Close()

	for err := range channel {
		t.Fatal(err)
	}
}

func makeGearmanCommand(cmd uint32, body []byte) ([]byte, error) {
	header := []byte{'\x00', 'R', 'E', 'S'}
	// 11 is JOB_ASSIGN
	cmdBytes, err := toBigEndianBytes(cmd)
	if err != nil {
		return nil, err
	}
	header = append(header, cmdBytes...)
	bodySize, err := toBigEndianBytes(uint32(len(body)))
	if err != nil {
		return nil, err
	}
	header = append(header, bodySize...)
	response := append(header, body...)
	return response, nil
}

// TestJobAssign tests that the worker runs the JOB_FUNC if the server sends a 'JOB_ASSIGN' packet.
// TODO: something in this job throws EOFs when run multiple times, fix plz
func TestJobAssign(t *testing.T) {

	name := "worker_name"
	workload := "the_workload"

	var channel chan error
	var listener net.Listener

	listener, channel = makeTCPServer(":1337", func(conn net.Conn) error {
		handle := "job_handle"
		function := name
		body := []byte(handle + string('\x00') + function + string('\x00') + workload)

		response, err := makeGearmanCommand(11, body)
		if err != nil {
			return err
		}
		if _, err := conn.Write(response); err != nil {
			return err
		}
		return nil
	})
	defer listener.Close()

	worker := NewWorker(name, func(job Job) ([]byte, error) {
		if string(job.Data()) != workload {
			close(channel)
			t.Fatalf("expected workload of '%s', received '%s'", workload, string(job.Data()))
		}
		close(channel)
		return []byte{}, nil
	})
	go worker.Listen("localhost", "1337")
	defer worker.Close()

	for err := range channel {
		t.Fatal(err)
	}
}

func GetClient() (c *client.Client) {
	c, err := client.New(client.Network, fmt.Sprintf("%s:%s", GearmanHost, GearmanPort))
	if err != nil {
		log.Fatalf("'%s', are you sure gearmand is running?", err)
	}
	c.ErrorHandler = func(e error) {
		log.Fatalln(e)
	}
	return c
}

// makes a job function that waits before completing
func getShutdownJobFn(doneChan chan string, readyChan chan int, workload string, sleepTime time.Duration) func(job Job) ([]byte, error) {
	return func(job Job) ([]byte, error) {
		readyChan <- 1
		time.Sleep(sleepTime)
		doneChan <- workload
		return []byte{}, nil
	}
}

func TestShutdownNoJob(t *testing.T) {
	worker := NewWorker("shutdown_no_job", func(job Job) ([]byte, error) {
		t.Fatalf("should not have invoked worker!")
		return []byte{}, nil
	})

	go worker.Listen(GearmanHost, GearmanPort)
	worker.Shutdown()
	return
}

// TestShutdown tests that the worker completes after worker.Shutdown is called
// make sure the next job is the second workload
func TestShutdown(t *testing.T) {
	c := GetClient()
	defer c.Close()

	// add jobs to client
	name := "shutdown_worker"

	_, err1 := c.DoBg(name, []byte("1"), client.JobNormal)
	if err1 != nil {
		log.Fatalln(err1)
	}

	_, err2 := c.DoBg(name, []byte("2"), client.JobNormal)
	if err2 != nil {
		t.Fatal(err2)
	}

	readyChan := make(chan int, 1)
	doneChan := make(chan string, 1)
	workload1 := "1"
	workload2 := "2"

	worker1 := NewWorker(name, getShutdownJobFn(doneChan, readyChan, workload1, 2*time.Second))
	go worker1.Listen(GearmanHost, GearmanPort)
	<-readyChan
	worker1.Shutdown()

	out1 := <-doneChan
	if out1 != workload1 {
		t.Fatalf("expected return of '%s', received '%s'", out1, workload1)
	}
	doneChan = make(chan string, 1)
	worker2 := NewWorker(name, getShutdownJobFn(doneChan, readyChan, workload2, 0))
	go worker2.Listen(GearmanHost, GearmanPort)
	out2 := <-doneChan
	if out2 != workload2 {
		t.Fatalf("expected return of '%s', received '%s'", out2, workload1)
	}
	return
}
