package launch

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ironiridis/tango/rpc"
)

// NewTask immediately starts a process and begins capturing any decoded messages
// into the provided chan. This is intended for fan-in style processing.
func NewTask(p string, rx chan<- *TaskMessage) (*Task, error) {
	t := new(Task)
	t.c = exec.Command(p)
	stdin, err := t.c.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("NewTask StdinPipe() failed: %w", err)
	}
	stdout, err := t.c.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("NewTask StdoutPipe() failed: %w", err)
	}
	t.h = rpc.NewHandle(stdin, stdout)
	t.c.Stderr = os.Stderr
	t.c.Start()
	go func() {
		defer t.h.Stop()
		msgch, err := t.h.Receive()
		if err != nil {
			return
		}
		for m := range msgch {
			tm := new(TaskMessage)
			tm.pid = t.c.Process.Pid
			tm.msg = m
			rx <- tm
		}
	}()
	return t, nil
}
