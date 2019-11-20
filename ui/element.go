package ui

import (
	"html/template"
	"sync"

	"github.com/ironiridis/tango/rpc"
)

type Element struct {
	sync.RWMutex
	UUID     rpc.UUID
	Label    string
	CSSClass string
	Template *template.Template
}
