package main

import (
	"github.com/vtpl1/vrtc3/internal/app"
	"github.com/vtpl1/vrtc3/internal/rtsp"
	"github.com/vtpl1/vrtc3/internal/streams"
	"github.com/vtpl1/vrtc3/pkg/shell"
)

func main() {
	app.Init()
	streams.Init()

	rtsp.Init()

	shell.RunUntilSignal()
}
