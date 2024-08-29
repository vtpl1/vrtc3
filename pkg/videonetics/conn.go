package videonetics

import (
	"context"
	"strconv"
	"sync"

	"github.com/vtpl1/vrtc3/pkg/core"
	pb "github.com/vtpl1/vrtc3/pkg/videonetics/service"
	"google.golang.org/grpc"
)

type Conn struct {
	core.Connection

	// internal
	uri     string
	ctx     *context.Context
	conn    *grpc.ClientConn
	host    string
	channel Channel

	state   State
	stateMu sync.Mutex
	stream  pb.StreamService_ReadFramePVAClient

	handler core.HandlerFunc
}

type State byte

const (
	StateNone State = iota
	StateConn
	StatePlay
)

const (
	MethodSetup = "SETUP"
	MethodPlay  = "PLAY"
)

func (s State) String() string {
	switch s {
	case StateNone:
		return "NONE"
	case StateConn:
		return "CONN"
	case StatePlay:
		return MethodPlay
	}
	return strconv.Itoa(int(s))
}
