package ivideon

import (
	"github.com/vtpl1/vrtc3/internal/streams"
	"github.com/vtpl1/vrtc3/pkg/core"
	"github.com/vtpl1/vrtc3/pkg/ivideon"
)

func Init() {
	streams.HandleFunc("ivideon", func(source string) (core.Producer, error) {
		return ivideon.Dial(source)
	})
}
