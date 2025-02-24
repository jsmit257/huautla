package types

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type ctxkey string

const (
	Cid     ctxkey = "cid"
	Metrics ctxkey = "metrics"
	Log     ctxkey = "log"
)

func GetContextCID(ctx context.Context) CID {
	if result, ok := ctx.Value(Cid).(CID); ok {
		return result
	}
	return CID(fmt.Sprintf("context has no cid attribute: %#v", ctx))
}

func GetContextLog(ctx context.Context) *logrus.Entry {
	if result, ok := ctx.Value(Log).(*logrus.Entry); ok {
		return result
	}

	l := logrus.WithFields(logrus.Fields{
		"ctx":   ctx,
		"bogus": true,
	})

	l.
		WithError(fmt.Errorf("context has no log attribute: %#v", ctx)).
		Error("getting context")

	return l
}

func GetContextMetrics(ctx context.Context) *prometheus.CounterVec {
	if result, ok := ctx.Value(Metrics).(*prometheus.CounterVec); ok {
		return result
	}

	return nil
	// return ServiceMetrics.MustCurryWith(prometheus.Labels{
	// 	"url":    "/missing/metrics/context/attribute",
	// 	"proto":  "ERROR",
	// 	"method": "metrics.GetContextMetrics",
	// })
}
