package isapi

import (
	"github.com/vtpl1/vrtc3/internal/streams"
	"github.com/vtpl1/vrtc3/pkg/core"
	"github.com/vtpl1/vrtc3/pkg/isapi"
)

func Init() {
	streams.HandleFunc("isapi", func(source string) (core.Producer, error) {
		return isapi.Dial(source)
	})
}
