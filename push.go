package main

import (
	"context"
	"flag"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"

	pb "push/proto/uplink"
)

type Options struct {
	Listen  bool
	Echo    bool
	Script  string
	Address string
	Timeout time.Duration
	Count   int
	Sleep   time.Duration
}

var (
	options Options
	flags   flag.FlagSet
)

func init() {
	flags.SetOutput(ioutil.Discard)
	flags.BoolVar(&options.Listen, "L", false, "listen")
	flags.BoolVar(&options.Echo, "E", false, "echo")
	flags.StringVar(&options.Script, "S", "", "script")
	flags.StringVar(&options.Address, "a", "", "server address")
	flags.DurationVar(&options.Timeout, "w", time.Second, "request timeout")
	flags.IntVar(&options.Count, "n", 1, "repeat count")
	flags.DurationVar(&options.Sleep, "z", time.Second, "sleep delay")
}

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		LogEcholn("Error: %s", err)
		os.Exit(1)
	}
	if options.Address == "" {
		LogEcholn("Address? (-a ...)")
		os.Exit(1)
	}
	if options.Echo && options.Script != "" {
		LogEcholn("Echo or script? (-E|-S ...)")
		os.Exit(1)
	}

	if options.Script != "" || options.Echo || options.Listen {
		ln, err := net.Listen("tcp", options.Address)
		if err != nil {
			LogEcholn("Error: %s", err)
			os.Exit(2)
		}
		var srv pb.MessagingServer
		if options.Script != "" {
			content, err := ioutil.ReadFile(options.Script)
			if err != nil {
				LogEcholn("Error: %s", err)
				os.Exit(2)
			}
			srv = NewScriptServer(strings.Split(string(content), "\n"))
		} else {
			srv = NewEchoServer(options.Echo)
		}
		s := grpc.NewServer()
		pb.RegisterMessagingServer(s, srv)
		s.Serve(ln)
	}

	conn, err := grpc.Dial(options.Address, grpc.WithInsecure())
	if err != nil {
		LogEcholn("Error: %s", err)
		os.Exit(2)
	}
	defer conn.Close()

	client := pb.NewMessagingClient(conn)
	topic, text := "PUSH/"+UID(), strings.Join(flags.Args(), " ")

	for i := 0; i < options.Count; i++ {
		n := func(id uint64, first, last bool) int {
			msg := &pb.Message{
				Id:    id,
				Topic: topic,
			}
			if first {
				msg.Text = text
			}
			Logf("Push: [%s][%d] %s", msg.Topic, msg.Id, msg.Text)
			req := &pb.Messages{Messages: []*pb.Message{msg}}
			resp, err := func() (*pb.Messages, error) {
				ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
				defer cancel()
				return client.Push(ctx, req)
			}()
			if err != nil {
				LogEcholn("Error: %s", err)
				os.Exit(2)
			}
			if len(resp.Messages) == 0 {
				return 0
			}
			for _, m := range resp.Messages {
				fmt.Println(m.Text)
				Logf("Recv: [%s][%d] %s", m.GetTopic(), m.GetId(), m.GetText())
			}
			if !last {
				Logf("...%s", options.Sleep)
				time.Sleep(options.Sleep)
			}
			return len(resp.Messages)
		}(uint64(i+1), i == 0, i == options.Count-1)
		if n == 0 {
			break
		}
	}
}

func UID() string {
	hash := fnv.New64a()
	hash.Write([]byte(strconv.Itoa(os.Getpid())))
	if ifas, err := net.Interfaces(); err != nil {
		for _, ifa := range ifas {
			if a := ifa.HardwareAddr.String(); a != "" {
				hash.Write([]byte(a))
			}
		}
	}
	return fmt.Sprintf("%016x", hash.Sum64())
}
