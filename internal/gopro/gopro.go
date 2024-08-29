package gopro

import (
	"net/http"

	"github.com/vtpl1/vrtc3/internal/api"
	"github.com/vtpl1/vrtc3/internal/streams"
	"github.com/vtpl1/vrtc3/pkg/core"
	"github.com/vtpl1/vrtc3/pkg/gopro"
)

func Init() {
	streams.HandleFunc("gopro", func(source string) (core.Producer, error) {
		return gopro.Dial(source)
	})

	api.HandleFunc("api/gopro", apiGoPro)
}

func apiGoPro(w http.ResponseWriter, r *http.Request) {
	var items []*api.Source

	for _, host := range gopro.Discovery() {
		items = append(items, &api.Source{Name: host, URL: "gopro://" + host})
	}

	api.ResponseSources(w, items)
}
