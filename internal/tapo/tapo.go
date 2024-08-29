package tapo

import (
	"github.com/vtpl1/vrtc3/internal/streams"
	"github.com/vtpl1/vrtc3/pkg/core"
	"github.com/vtpl1/vrtc3/pkg/kasa"
	"github.com/vtpl1/vrtc3/pkg/tapo"
)

func Init() {
	streams.HandleFunc("kasa", func(source string) (core.Producer, error) {
		return kasa.Dial(source)
	})

	streams.HandleFunc("tapo", func(source string) (core.Producer, error) {
		return tapo.Dial(source)
	})
}
