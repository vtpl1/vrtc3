package expr

import (
	"errors"

	"github.com/vtpl1/vrtc3/internal/app"
	"github.com/vtpl1/vrtc3/internal/streams"
	"github.com/vtpl1/vrtc3/pkg/expr"
)

func Init() {
	log := app.GetLogger("expr")

	streams.RedirectFunc("expr", func(url string) (string, error) {
		v, err := expr.Run(url[5:])
		if err != nil {
			return "", err
		}

		log.Debug().Msgf("[expr] url=%s", url)

		if url = v.(string); url == "" {
			return "", errors.New("expr: result is empty")
		}

		return url, nil
	})
}
