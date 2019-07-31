package main

import (
	"context"
	"fmt"

	pb "push/proto/uplink"
)

//
// EchoServer

type EchoServer struct {
	echo   bool
	lastid uint64
}

func NewEchoServer(echo bool) *EchoServer {
	return &EchoServer{echo: echo}
}

func (s *EchoServer) Push(ctx context.Context, request *pb.Messages) (*pb.Messages, error) {
	reply := &pb.Messages{}
	for _, msg := range request.Messages {
		fmt.Printf("[%s][%d] %s\n", msg.GetTopic(), msg.GetId(), msg.GetText())
		if s.echo {
			s.lastid++
			reply.Messages = append(reply.Messages, &pb.Message{
				Id:    s.lastid,
				Topic: msg.GetTopic(),
				Text:  msg.GetText(),
			})
		}
	}
	return reply, nil
}

//
// ScriptServer

type ScriptServer struct {
	script []string
	lastid uint64
}

func NewScriptServer(script []string) *ScriptServer {
	s := &ScriptServer{script: make([]string, len(script))}
	copy(s.script, script)
	return s
}

func (s *ScriptServer) Push(ctx context.Context, request *pb.Messages) (*pb.Messages, error) {
	reply := &pb.Messages{}
	if len(request.Messages) == 0 {
		return reply, nil
	}
	topic := ""
	for _, msg := range request.Messages {
		fmt.Printf("[%s][%d] %s\n", msg.GetTopic(), msg.GetId(), msg.GetText())
		topic = msg.GetTopic()
	}
	for len(s.script) > 0 {
		l := s.script[0]
		s.script = s.script[1:]
		if l == "" {
			break
		}
		fmt.Println(l)
		s.lastid++
		reply.Messages = append(reply.Messages, &pb.Message{
			Id:    s.lastid,
			Topic: topic,
			Text:  l,
		})
	}
	return reply, nil
}
