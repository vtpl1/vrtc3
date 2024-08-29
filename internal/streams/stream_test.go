package streams

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vtpl1/vrtc3/pkg/core"
)

func TestRecursion(t *testing.T) {
	// create stream with some source
	stream1 := New("from_yaml", "does_not_matter")
	require.Len(t, streams, 1)

	// ask another unnamed stream that links vrtc3
	query, err := url.ParseQuery("src=rtsp://localhost:8554/from_yaml?video")
	require.Nil(t, err)
	stream2 := GetOrPatch(query)

	// check stream is same
	require.Equal(t, stream1, stream2)
	// check stream urls is same
	require.Equal(t, stream1.producers[0].url, stream2.producers[0].url)
	require.Len(t, streams, 2)
}

func TestTempate(t *testing.T) {
	HandleFunc("rtsp", func(url string) (core.Producer, error) { return nil, nil }) // bypass HasProducer

	// config from yaml
	stream1 := New("camera.from_hass", "ffmpeg:{input}#video=copy")
	// request from hass
	stream2 := Patch("camera.from_hass", "rtsp://example.com")

	require.Equal(t, stream1, stream2)
	require.Equal(t, "ffmpeg:rtsp://example.com#video=copy", stream1.producers[0].url)
}
