package prometheus

import (
	"bookstore/demo/web"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MiddlewareBuilder struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
}

func (b MiddlewareBuilder) Build() web.Middleware {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      b.Name,
		Subsystem: b.Subsystem,
		Namespace: b.Namespace,
		Help:      b.Help,
	}, []string{"pattern", "method", "status"})
	// 注册一下
	prometheus.MustRegister(vector)
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			startTime := time.Now()
			next(ctx)
			go report(time.Since(startTime), ctx, vector)
		}
	}
}

func report(dur time.Duration, ctx *web.Context, vec prometheus.ObserverVec) {
	status := ctx.RespStatusCode
	route := "unknown"
	if ctx.MatchedRoute != "" {
		route = ctx.MatchedRoute
	}

	ms := dur / time.Millisecond
	vec.WithLabelValues(route, ctx.Req.Method, strconv.Itoa(status)).Observe(float64(ms))
}
