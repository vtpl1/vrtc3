package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/vtpl1/vrtc3/internal/api"
	"github.com/vtpl1/vrtc3/internal/api/ws"
	"github.com/vtpl1/vrtc3/internal/app"
	"github.com/vtpl1/vrtc3/internal/bubble"
	"github.com/vtpl1/vrtc3/internal/debug"
	"github.com/vtpl1/vrtc3/internal/dvrip"
	"github.com/vtpl1/vrtc3/internal/echo"
	"github.com/vtpl1/vrtc3/internal/exec"
	"github.com/vtpl1/vrtc3/internal/expr"
	"github.com/vtpl1/vrtc3/internal/ffmpeg"
	"github.com/vtpl1/vrtc3/internal/gopro"
	"github.com/vtpl1/vrtc3/internal/hass"
	"github.com/vtpl1/vrtc3/internal/hls"
	"github.com/vtpl1/vrtc3/internal/homekit"
	"github.com/vtpl1/vrtc3/internal/http"
	"github.com/vtpl1/vrtc3/internal/isapi"
	"github.com/vtpl1/vrtc3/internal/ivideon"
	"github.com/vtpl1/vrtc3/internal/mjpeg"
	"github.com/vtpl1/vrtc3/internal/mp4"
	"github.com/vtpl1/vrtc3/internal/mpegts"
	"github.com/vtpl1/vrtc3/internal/nest"
	"github.com/vtpl1/vrtc3/internal/ngrok"
	"github.com/vtpl1/vrtc3/internal/onvif"
	"github.com/vtpl1/vrtc3/internal/roborock"
	"github.com/vtpl1/vrtc3/internal/rtmp"
	"github.com/vtpl1/vrtc3/internal/rtsp"
	"github.com/vtpl1/vrtc3/internal/srtp"
	"github.com/vtpl1/vrtc3/internal/streams"
	"github.com/vtpl1/vrtc3/internal/tapo"
	"github.com/vtpl1/vrtc3/internal/videonetics"
	"github.com/vtpl1/vrtc3/internal/webrtc"
	"github.com/vtpl1/vrtc3/internal/webtorrent"
)

func main() {
	app.Version = "1.9.4"
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	app.InternalTerminationRequest = make(chan int)

	// 1. Core modules: app, api/ws, streams

	app.Init() // init config and logs

	api.Init() // init API before all others
	ws.Init()  // init WS API endpoint

	streams.Init() // streams module

	// 2. Main sources and servers

	rtsp.Init()   // rtsp source, RTSP server
	webrtc.Init() // webrtc source, WebRTC server

	// 3. Main API

	mp4.Init()   // MP4 API
	hls.Init()   // HLS API
	mjpeg.Init() // MJPEG API

	// 4. Other sources and servers

	hass.Init()       // hass source, Hass API server
	onvif.Init()      // onvif source, ONVIF API server
	webtorrent.Init() // webtorrent source, WebTorrent module

	// 5. Other sources

	rtmp.Init()     // rtmp source
	exec.Init()     // exec source
	ffmpeg.Init()   // ffmpeg source
	echo.Init()     // echo source
	ivideon.Init()  // ivideon source
	http.Init()     // http/tcp source
	dvrip.Init()    // dvrip source
	tapo.Init()     // tapo source
	isapi.Init()    // isapi source
	mpegts.Init()   // mpegts passive source
	roborock.Init() // roborock source
	homekit.Init()  // homekit source
	nest.Init()     // nest source
	bubble.Init()   // bubble source
	expr.Init()     // expr source
	gopro.Init()    // gopro source

	videonetics.Init(&ctx) // videonetics source

	// 6. Helper modules

	ngrok.Init() // ngrok module
	srtp.Init()  // SRTP server
	debug.Init() // debug API

	// 7. Go
	doShutdown := false
	for !doShutdown {
		select {
		case <-ctx.Done():
			doShutdown = true
		case <-app.InternalTerminationRequest:
			stop()
			doShutdown = true
		}
	}

	// shell.RunUntilSignal()
}
