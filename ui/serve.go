package ui

import "net/http"

type Server struct{}

func (s *Server) startHTTP() {
	r := http.ListenAndServe("0.0.0.0:80", nil)
	panic(r)
}

func (s *Server) startHTTPS() {
	// We need to generate an ad-hoc certificate. tls.Config.GetCertificate seems to be the way
	// to do this; we can cache the generated certificate and reuse it for future traffic.
	// Avoid doing this at startup since it may be resource intensive.
	//TODO: As above
}

func (s *Server) Start() {
	go s.startHTTP()
	go s.startHTTPS()
}
