package services

import (
	"github.com/Qwepo/rusprofile/gen/rusprof"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Services struct {
	log zerolog.Logger
}

func NewServices(log zerolog.Logger) *Services {
	return &Services{log: log}
}

func (s *Services) RegisterRusprof(server *grpc.Server) {
	rusprof.RegisterRusprofServer(server, s)
}
