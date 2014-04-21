# baseworker-go

`baseworker-go` is a wrapper around [gearman-go](https://github.com/mikespook/gearman-go) to provide
a simpler API around creating a worker.

## Example

```go
package main

import(
    "github.com/Clever/baseworker-go/worker"
    "log"
)

func jobFunc(job gearman.Job) ([]byte, error) {
    log.Printf("Got job with data %s", job.Data())
    return []byte{}, nil
}

func main() {
    worker := gearman.New("test", jobFunc)
    worker.Listen("localhost", "4730")
}
```

## API

View API documentation with the `godoc` command, by running `godoc .` in the project root.
