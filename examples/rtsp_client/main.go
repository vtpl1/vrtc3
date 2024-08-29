package main

import (
	"log"
	"os"

	"github.com/vtpl1/vrtc3/pkg/core"
	"github.com/vtpl1/vrtc3/pkg/rtsp"
	"github.com/vtpl1/vrtc3/pkg/shell"
)

func main() {
	client := rtsp.NewClient(os.Args[1])
	if err := client.Dial(); err != nil {
		log.Panic(err)
	}

	client.Medias = []*core.Media{
		{
			Kind:      core.KindAudio,
			Direction: core.DirectionRecvonly,
			Codecs: []*core.Codec{
				{Name: core.CodecPCMU, ClockRate: 8000},
			},
			ID: "streamid=0",
		},
	}
	if err := client.Announce(); err != nil {
		log.Panic(err)
	}
	if _, err := client.SetupMedia(client.Medias[0]); err != nil {
		log.Panic(err)
	}
	if err := client.Record(); err != nil {
		log.Panic(err)
	}

	shell.RunUntilSignal()
}
