package application

import (
	"github.com/orangebees/go-oneutils/Random"
	"github.com/savsgio/atreugo/v11"
	"os"
	"strconv"
)

var atreugoServer *atreugo.Atreugo

func GetAtreugo() *atreugo.Atreugo {
	if atreugoServer != nil {
		return atreugoServer
	}
	return createAtreugo()
}

func createAtreugo() *atreugo.Atreugo {
	atg := atreugo.Config{
		Logger:                GetLogger(),
		NoDefaultServerHeader: true,
		NoDefaultDate:         true,
		GracefulShutdown:      true,
		Addr:                  "0.0.0.0:9000",
		LogAllErrors:          false,
		MaxConnsPerIP:         0,
		Prefork:               false,
		ReduceMemoryUsage:     false,
		Compress:              false,
	}
	if v := os.Getenv("ATREUGO_ADDR"); v != "" {
		atg.Addr = v
	}
	if v := os.Getenv("ATREUGO_MAXCONNSPERIP"); v != "" {
		parseInt, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			atg.MaxConnsPerIP = int(parseInt)
		}
	}
	if v := os.Getenv("ATREUGO_PREFORK"); v == "true" {
		atg.Prefork = true
	}
	if v := os.Getenv("ATREUGO_REDUCEMEMORYUSAGE"); v == "true" {
		atg.ReduceMemoryUsage = true
	}
	if v := os.Getenv("ATREUGO_COMPRESS"); v == "true" {
		atg.Compress = true
	}
	if v := os.Getenv("ATREUGO_LOGALLERRORS"); v == "true" {
		atg.LogAllErrors = true
	}
	server := atreugo.New(atg)
	//pc := prometheus.Config{}
	//if v := os.Getenv("ATREUGO_PROMETHEUS_METHOD"); v != "" {
	//	pc.Method = v
	//}
	//if v := os.Getenv("ATREUGO_PROMETHEUS_URL"); v != "" {
	//	pc.URL = v
	//}
	//if v := os.Getenv("ATREUGO_PROMETHEUS"); v == "true" {
	//	prometheus.Register(server, pc)
	//}
	return server
}

func jsonContentType(ctx *atreugo.RequestCtx) error {
	ctx.Response.Header.SetContentType("application/json;charset=utf-8")
	return ctx.Next()
}
func SetJsonContentType(a *atreugo.Router) {
	a.UseAfter(jsonContentType)
}
func SetHealthResponse(a *atreugo.Atreugo) {
	a.GET("/health", func(ctx *atreugo.RequestCtx) error {
		ctx.SetBodyString("ok")
		return nil
	})
}
func SetRequestId(a *atreugo.Router) {
	a.UseBefore(func(ctx *atreugo.RequestCtx) error {
		requestId := ctx.RequestID()
		if requestId == nil {
			requestId = Random.RandBytes32()
		}
		ctx.SetUserValue("RequestId", requestId)
		ctx.Response.Header.SetBytesV(atreugo.XRequestIDHeader, requestId)
		return ctx.Next()
	})
}
func SetAccessLog(a *atreugo.Router) {
	a.UseBefore(func(ctx *atreugo.RequestCtx) error {
		event := GetLogger().Info()
		requestId, ok := ctx.UserValue("RequestId").([]byte)
		if ok {
			event = event.Bytes("RequestId",
				requestId)
		}
		if ip := ctx.Request.Header.Peek("X-Forwarded-For"); ip != nil {
			event = event.Bytes("RemoteIP", ip)
		} else {
			event = event.Str("RemoteIP", ctx.RemoteIP().String())
		}
		event.Bytes("Method", ctx.Method()).Str("URI", ctx.URI().String())
		event.Msg("starting...")
		return ctx.Next()
	}).UseAfter(func(ctx *atreugo.RequestCtx) error {
		event := GetLogger().Info()
		requestId, ok := ctx.UserValue("RequestId").([]byte)
		if ok {
			event = event.Bytes("RequestId", requestId)
		}
		event.Msg("ending...")
		return ctx.Next()
	})
}
