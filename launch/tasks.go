package launch

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/ironiridis/tango/rpc"
)

type Task struct {
	c *exec.Cmd
	h *rpc.Handle
}

func (t *Task) Pid() int {
	return t.c.Process.Pid
}

type TaskMessage struct {
	pid int
	msg *rpc.Message
}

type TaskList struct {
	sync.Mutex
	m map[int]*Task
}

func ipcstuff(tl *TaskList, ch <-chan *TaskMessage) {
	for m := range ch {
		fmt.Printf("tango msg: %#v\n", *m)
		if m.msg.T != "ipc" {
			continue
		}
		if m.msg.Optional {
			continue
		}

		tl.Lock()
		for p := range tl.m {
			if p != m.pid {
				err := tl.m[p].h.Send(m.msg)
				if err != nil {
					panic(err)
				}
			}
		}
		tl.Unlock()
	}
}

func StartTasks() {
	tl := new(TaskList)
	tl.m = make(map[int]*Task)
	tmch := make(chan *TaskMessage)
	go ipcstuff(tl, tmch)
	t1, err := NewTask(`C:\Users\Chris Harrington\go\src\github.com\ironiridis\tango\sandbox\echotest\echotest.exe`, tmch)
	if err != nil {
		panic(err)
	}
	tl.m[t1.Pid()] = t1

	t2, err := NewTask(`C:\Users\Chris Harrington\go\src\github.com\ironiridis\tango\sandbox\echotest\echotest.exe`, tmch)
	if err != nil {
		panic(err)
	}
	tl.m[t2.Pid()] = t2

	t3, err := NewTask(`C:\Users\Chris Harrington\go\src\github.com\ironiridis\tango\sandbox\talkertest\talkertest.exe`, tmch)
	if err != nil {
		panic(err)
	}

	tl.m[t3.Pid()] = t3
	time.Sleep(time.Second * 30)
	fmt.Printf("tl: %+v\n", tl)
}
