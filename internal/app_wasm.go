//go:build js && wasm

package internal

import (
	"fmt"
	"log"
	"syscall/js"

	"github.com/comame/note.comame.xyz/internal/md"
)

func RunApp() {
	js.Global().Set("go_parseMarkdown", js.FuncOf(parseMarkdown))
	log.Println("ready")

	<-make(chan struct{})
}

func parseMarkdown(_ js.Value, args []js.Value) interface{} {
	// notify the error to browser then re-panic
	defer func() {
		if v := recover(); v != nil {
			js.Global().Call("alert", js.ValueOf(fmt.Sprintf("%v", v)))
			panic(v)
		}
	}()

	markdown := args[0].String()
	html := md.ToHTML(markdown)

	return js.ValueOf(html)
}
