package main

import (
	"context"
	"fmt"

	pb "push/proto/uplink"
)

//
// Server

type Server struct {
	Echo   bool
	nextid uint64
}

func NewServer(echo bool) *Server {
	return &Server{
		Echo:   echo,
		nextid: 1,
	}
}

func (s *Server) Push(ctx context.Context, request *pb.Messages) (*pb.Messages, error) {
	reply := &pb.Messages{}
	for _, msg := range request.Messages {
		fmt.Printf("[%s][%d] %s\n", msg.GetTopic(), msg.GetId(), msg.GetText())
		if s.Echo {
			rep := &pb.Message{
				Id:    s.nextid,
				Topic: msg.GetTopic(),
				Text:  msg.GetText(),
			}
			reply.Messages = append(reply.Messages, rep)
			s.nextid++
		}
	}
	return reply, nil
}
