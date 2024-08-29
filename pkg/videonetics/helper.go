package videonetics

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

func ParseVideoneticsUri(uri string) (string, Channel, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", Channel{}, nil
	}
	if u.Scheme != "videonetics" {
		err = errors.New("the scheme must be of videoentics")
		return "", Channel{}, err
	}
	host := "dns:///" + u.Host
	paths := strings.Split(u.Path, "/")
	// SiteID     int64  `json:"site_id"`
	// ChannelID  int64  `json:"channel_id"`
	// AppID      int64  `json:"app_id"`
	// LiveOrRec  int32  `json:"live_or_rec"`
	// StreamType int32  `json:"stream_type"`
	// StartTS    int64  `json:"start_ts"`
	// SessionID  string `json:"session_id"`
	var channel Channel = Channel{
		LiveOrRec: 1,
	}
	positionToScan := 7
	if len(paths) > positionToScan {
		channel.SessionID = paths[positionToScan]
	}
	positionToScan = 6
	if len(paths) > positionToScan {
		i, err := strconv.ParseInt(paths[positionToScan], 10, 64)
		if err != nil {
			return host, channel, err
		}
		channel.StartTS = i
	}
	positionToScan = 5
	if len(paths) > positionToScan {
		i, err := strconv.ParseInt(paths[positionToScan], 10, 32)
		if err != nil {
			return host, channel, err
		}
		channel.StreamType = int32(i)
	}
	positionToScan = 4
	if len(paths) > positionToScan {
		i, err := strconv.ParseInt(paths[positionToScan], 10, 32)
		if err != nil {
			return host, channel, err
		}
		channel.LiveOrRec = int32(i)
	}
	positionToScan = 3
	if len(paths) > positionToScan {
		i, err := strconv.ParseInt(paths[positionToScan], 10, 64)
		if err != nil {
			return host, channel, err
		}
		channel.AppID = i
	}
	positionToScan = 2
	if len(paths) > positionToScan {
		i, err := strconv.ParseInt(paths[positionToScan], 10, 64)
		if err != nil {
			return host, channel, err
		}
		channel.ChannelID = i
	}
	positionToScan = 1
	if len(paths) > positionToScan {
		i, err := strconv.ParseInt(paths[positionToScan], 10, 64)
		if err != nil {
			return host, channel, err
		}
		channel.SiteID = i
	}
	return host, channel, nil
}
