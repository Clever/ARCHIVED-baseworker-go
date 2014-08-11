# mock
--
    import "github.com/Clever/baseworker-go/mock"


## Usage

#### type MockJob

```go
type MockJob struct {
}
```


#### func  CreateMockJob

```go
func CreateMockJob(payload string) *MockJob
```
CreateMockJob creates an object that implements the gearman-go/worker#Job
interface

#### func (MockJob) Data

```go
func (m MockJob) Data() []byte
```
Data returns the job's data

#### func (MockJob) Err

```go
func (m MockJob) Err() error
```
Err returns the job's error

#### func (MockJob) Fn

```go
func (m MockJob) Fn() string
```
Fn returns the name of the worker's job

#### func (MockJob) Handle

```go
func (m MockJob) Handle() string
```
Handle returns the job handle

#### func (*MockJob) SendData

```go
func (m *MockJob) SendData(data []byte)
```
SendData appends to the array of job data

#### func (*MockJob) SendWarning

```go
func (m *MockJob) SendWarning(warning []byte)
```
SendWarning appends to the array of job warnings

#### func (MockJob) UniqueId

```go
func (m MockJob) UniqueId() string
```
UniqueId returns the unique id for the job

#### func (*MockJob) UpdateStatus

```go
func (m *MockJob) UpdateStatus(numerator, denominator int)
```
UpdateStatus updates the progress of job

#### func (*MockJob) Warnings

```go
func (m *MockJob) Warnings() [][]byte
```
Warnings returns the array of warnings for the job
