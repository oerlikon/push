package main

import (
	"context"
	"fmt"

	pb "push/proto/uplink"
)

//
// MessagingServer

type MessagingServer struct {
	echo   bool
	script []string
	nextid uint64
}

func NewMessagingServer() *MessagingServer {
	return &MessagingServer{nextid: 1}
}

func (s *MessagingServer) SetEcho(echo bool) {
	s.echo = echo
}

func (s *MessagingServer) SetScript(script []string) {
	s.script = make([]string, len(script))
	copy(s.script, script)
}

func (s *MessagingServer) Push(ctx context.Context, request *pb.Messages) (*pb.Messages, error) {
	reply := &pb.Messages{}
	for _, msg := range request.Messages {
		fmt.Printf("[%s][%d] %s\n", msg.GetTopic(), msg.GetId(), msg.GetText())
		if s.echo {
			rep := &pb.Message{
				Id:    s.nextid,
				Topic: msg.GetTopic(),
				Text:  msg.GetText(),
			}
			reply.Messages = append(reply.Messages, rep)
			s.nextid++
		}
		for len(s.script) > 0 {
			l := s.script[0]
			s.script = s.script[1:]
			if l == "" {
				break
			}
			rep := &pb.Message{
				Id:    s.nextid,
				Topic: msg.GetTopic(),
				Text:  l,
			}
			reply.Messages = append(reply.Messages, rep)
			s.nextid++
		}
	}
	return reply, nil
}
