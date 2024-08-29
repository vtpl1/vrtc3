package main

import (
	"github.com/vtpl1/vrtc3/internal/api"
	"github.com/vtpl1/vrtc3/internal/app"
	"github.com/vtpl1/vrtc3/internal/hass"
	"github.com/vtpl1/vrtc3/internal/streams"
	"github.com/vtpl1/vrtc3/pkg/shell"
)

func main() {
	app.Init()
	streams.Init()

	api.Init()

	hass.Init()

	shell.RunUntilSignal()
}
