package videonetics

import (
	"context"

	"github.com/vtpl1/vrtc3/pkg/core"
)

type Channel struct {
	SiteID     int64  `json:"site_id"`
	ChannelID  int64  `json:"channel_id"`
	AppID      int64  `json:"app_id"`
	LiveOrRec  int32  `json:"live_or_rec"`
	StreamType int32  `json:"stream_type"`
	StartTS    int64  `json:"start_ts"`
	SessionID  string `json:"session_id"`
}

func NewClient(uri string, ctx *context.Context) *Conn {

	host, channel, err := ParseVideoneticsUri(uri)
	if err != nil {
		return nil
	}
	return &Conn{
		Connection: core.Connection{
			ID:         core.NewID(),
			FormatName: "videonetics",
			// Medias:     getMedias(),
		},
		uri:     uri,
		ctx:     ctx,
		host:    host,
		channel: channel,
	}
}
