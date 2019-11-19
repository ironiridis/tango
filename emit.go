package main

import (
	"fmt"
	"os"
)

func emitWarning(m string) {
	fmt.Fprint(os.Stderr, m)
}
