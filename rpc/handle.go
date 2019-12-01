package rpc

import (
	"encoding/json"
	"io"
	"sync"
)

// A Handle holds the JSON encoder/decoder pair for RPC to a process.
type Handle struct {
	sync.Mutex
	r   io.ReadCloser
	w   io.WriteCloser
	e   *json.Encoder
	ch  chan *Message
	err error
}

type M interface {
	T() string
}

// HandleNotReady is returned when a Handle is used in an invalid state, such as before
// it has been created with NewHandle() or after a call to h.Stop().
const HandleNotReady = Error("Handle not ready")

// NewHandle returns a Handle initialized with the the Reader and Writer for a process.
func NewHandle(w io.WriteCloser, r io.ReadCloser) *Handle {
	h := new(Handle)
	h.r = r
	h.w = w
	h.e = json.NewEncoder(w)
	h.e.SetEscapeHTML(false)

	return h
}

// Err returns the last error encountered with this Handle, and clears it.
func (h *Handle) Err() error {
	h.Lock()
	defer h.Unlock()

	e := h.err
	h.err = nil
	return e
}

// SendRaw takes a value `v` and writes the JSON encoded form to stdin of the
// receiving process. SendRaw is safe to use from multiple goroutines.
func (h *Handle) SendRaw(v *Message) error {
	h.Lock()
	defer h.Unlock()

	if h.e == nil || h.w == nil {
		h.err = HandleNotReady
	}
	if h.err != nil {
		return h.err
	}

	return h.e.Encode(*v)
}

// Send wraps an arbitrary value into a Message with metadata.
func (h *Handle) Send(z M, optional bool) error {
	js, err := json.Marshal(z)
	if err != nil {
		return err
	}

	var m Message
	m.T = z.T()
	m.Optional = optional
	m.M = json.RawMessage(js)
	return h.SendRaw(&m)
}

func (h *Handle) reader() {
	if h.r == nil {
		panic(Error("handle reader is nil"))
	}
	if h.ch == nil {
		panic(Error("handle receiver channel is nil"))
	}
	d := json.NewDecoder(h.r)
	for {
		v := new(Message)
		err := d.Decode(v)
		if err != nil {
			h.Lock() // we modify h.ch and h.err
			defer h.Unlock()
			close(h.ch)
			h.ch = nil
			if err != io.EOF {
				h.err = err
			}
			return
		}
		h.ch <- v
	}
}

// Receive allocates and runs a JSON decoder for the stdout of the process,
// returning a channel for decoded messages. If a decoder is already running,
// the existing channel is reused.
func (h *Handle) Receive() (<-chan *Message, error) {
	h.Lock()
	defer h.Unlock()

	if h.r == nil {
		h.err = HandleNotReady
	}
	if h.err != nil {
		return nil, h.err
	}
	if h.ch != nil {
		// return existing decoder goroutine
		return h.ch, nil
	}
	h.ch = make(chan *Message)
	go h.reader()
	return h.ch, nil
}

// Stop closes stdin on the related process, discards the JSON encoder, and
// disallows future operations on this Handle.
func (h *Handle) Stop() {
	h.Lock()
	defer h.Unlock()

	h.e = nil
	h.w.Close()
	h.w = nil // cause any future writes to panic
	h.r = nil // disallow any future calls to h.Receive()
}
