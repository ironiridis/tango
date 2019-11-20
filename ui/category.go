package ui

import (
	"sync"

	"github.com/ironiridis/tango/rpc"
)

type Category struct {
	sync.RWMutex
	UUID     rpc.UUID
	Label    string
	CSSClass string
	Children []Category
	Elements []Element
}
