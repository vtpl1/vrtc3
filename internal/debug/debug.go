package debug

import (
	"github.com/vtpl1/vrtc3/internal/api"
)

func Init() {
	api.HandleFunc("api/stack", stackHandler)
}
