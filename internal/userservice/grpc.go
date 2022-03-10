// Created by Hisen at 2022/3/2.
package userservice

import (
	"github.com/hanxinhisen/moss/pkg/log"
	"google.golang.org/grpc"
	"net"
)

type grpcAPIServer struct {
	*grpc.Server
	address string
}

func (s *grpcAPIServer) Run() {
	listen, err := net.Listen("tcp", s.address)
	if err != nil {

		log.Fatalf("fail to listen:%s\n", err.Error())
	}
	go func() {

		if err := s.Serve(listen); err != nil {

			log.Fatalf("failed to start grpc server:%s", err.Error())
		}

	}()
	log.Infof("start grpc server at %s", s.address)

}

func (s *grpcAPIServer) Close() {
	s.GracefulStop()
	log.Infof("GRPC server on %s stopped", s.address)
}
