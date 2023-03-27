package pprof

import (
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

const (
	// defaultPrefix url prefix of pprof
	defaultPrefix = "/debug/pprof"
)

// Router the standard HandlerFuncs from the net/http/pprof package with
// the provided gin.Engine. prefixOptions is a optional. If not prefixOptions,
// the default path prefix("/debug/pprof") is used, otherwise first prefixOptions will be path prefix.
func Router(router *gin.Engine, prefixOptions ...string) {
	prefix := defaultPrefix
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}

	g := router.Group(prefix)
	{
		g.GET("/", gin.WrapF(pprof.Index))
		g.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		g.GET("/profile", gin.WrapF(pprof.Profile))
		g.POST("/symbol", gin.WrapF(pprof.Symbol))
		g.GET("/symbol", gin.WrapF(pprof.Symbol))
		g.GET("/trace", gin.WrapF(pprof.Trace))
		g.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		g.GET("/block", gin.WrapH(pprof.Handler("block")))
		g.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		g.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		g.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		g.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}
}
