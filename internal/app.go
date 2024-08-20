//go:build !js || !wasm

package internal

import (
	"github.com/comame/note.comame.xyz/internal/server"
)

func RunApp() {
	server.Start()
}
