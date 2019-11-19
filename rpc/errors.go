package rpc

// Error is an error interface type for various constant RPC error values
type Error string

func (e Error) Error() string {
	return string(e)
}

func (e Error) String() string {
	return string(e)
}
