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
	nextid uint64
}

func NewEchoServer(echo bool) *EchoServer {
	return &EchoServer{
		echo:   echo,
		nextid: 1,
	}
}

func (s *EchoServer) Push(ctx context.Context, request *pb.Messages) (*pb.Messages, error) {
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
	}
	return reply, nil
}

//
// ScriptServer

type ScriptServer struct {
	script []string
	nextid uint64
}

func NewScriptServer(script []string) *ScriptServer {
	s := &ScriptServer{
		script: make([]string, len(script)),
		nextid: 1,
	}
	copy(s.script, script)
	return s
}

func (s *ScriptServer) Push(ctx context.Context, request *pb.Messages) (*pb.Messages, error) {
	reply := &pb.Messages{}
	topic := ""
	if len(request.Messages) == 0 {
		return reply, nil
	}
	for _, msg := range request.Messages {
		if topic == "" {
			topic = msg.GetTopic()
		}
		fmt.Printf("[%s][%d] %s\n", msg.GetTopic(), msg.GetId(), msg.GetText())
	}
	for len(s.script) > 0 {
		l := s.script[0]
		s.script = s.script[1:]
		if l == "" {
			break
		}
		rep := &pb.Message{
			Id:    s.nextid,
			Topic: topic,
			Text:  l,
		}
		reply.Messages = append(reply.Messages, rep)
		fmt.Println(l)
		s.nextid++
	}
	return reply, nil
}
