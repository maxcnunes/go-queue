go-queue
========

A modest implementation of queue in go. The main advantage of this queue is that it allows running tasks concurrently.

## Usage

```go
import "github.com/maxcnunes/go-queue"

const CONCURRENT = 10

type TaskRunner struct {
  // include any context data
}

func (t TaskRunner) Run(item interface{}) error {
  // data := item.(YOUR_DATA_TYPE)
}

func main() {
  q := goqueue.NewQueue(
    CONCURRENT,
    MergeCardBranchesTaskRunner{},
  ).Run()

  // seed the queue
  go loadJobs(git, q)

  // wait finish processing jobs
  if err := q.Drain(); err != nil {
    log.Printf("Error %#v\n", err)
  }
}

func loadJobs(q *core.Queue) {
  q.Jobs <- YOUR_DATA_TYPE{}

  if err != nil {
    q.Errors <- err
    close(q.Jobs)
    return
  }

  // close the queue if your are not expecting more income data
  close(q.Jobs)
}
```
