package ginX

import (
	"log"
	"net/http"
	"testing"
)

func TestEngine(t *testing.T) {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe("0.0.0.0:9999", engine))
}

func TestCtx(t *testing.T) {
	engine := New()
	engine.addRoute("GET", "/hello", func(ctx *Context) {
		ctx.String(http.StatusOK, "hello ginX")
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:9999", engine))
}
