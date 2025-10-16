package socket

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/ggymm/gnet"
	"github.com/panjf2000/ants/v2"

	"chat-server/config"
)

type Socket struct {
	gnet.BuiltinEventEngine
	eng gnet.Engine

	// 协程池
	pool *ants.Pool

	// 服务配置
	addr      string
	reuse     bool
	multicore bool
}

func NewSocket() *Socket {
	return &Socket{
		addr:      config.AppSocket,
		reuse:     true,
		multicore: true,
	}
}

func (s *Socket) Start() error {
	return gnet.Run(
		s, s.addr,
		gnet.WithLogger(&SocketLogger{
			log: newLog(),
		}),
		gnet.WithReuseAddr(s.reuse),
		gnet.WithMulticore(s.multicore),
	)
}

func (s *Socket) OnBoot(eng gnet.Engine) (action gnet.Action) {
	s.eng = eng

	slog.Info("socket server started")
	return
}

func (s *Socket) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	c.SetContext(&SocketCodec{})

	slog.Info("socket client connected",
		slog.String("ip", c.RemoteAddr().String()),
	)
	return
}

func (s *Socket) OnTraffic(c gnet.Conn) (action gnet.Action) {
	var codec = c.Context().(*SocketCodec)
	var packets [][]byte
	for {
		p, err := codec.Decode(c)
		if errors.Is(err, errIncompletePacket) {
			break
		}
		if err != nil {
			return gnet.Close
		}
		packets = append(packets, p)
	}

	for _, p := range packets {
		err := s.pool.Submit(func() {
			fmt.Printf("%s", p)
		})
		if err != nil {
			slog.Error(err.Error())
			continue
		}
	}
	return
}
