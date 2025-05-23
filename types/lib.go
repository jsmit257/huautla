package types

import (
	"context"
	"fmt"
	"strings"
	"time"

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

func GetContextDataMetrics(ctx context.Context) *prometheus.CounterVec {
	if result, ok := ctx.Value(Metrics).(*prometheus.CounterVec); ok {
		return result
	}

	return DataMetrics.MustCurryWith(prometheus.Labels{
		"pkg":      "data",
		"function": "ERROR",
		"status":   "no data metrics set",
	})
}

type BadIntervalError error
type InvalidTimestampError error

var validIntervals = map[string]struct{}{
	"hour":  {},
	"day":   {},
	"week":  {},
	"month": {},
	"year":  {},
}

func (ts *Timestamp) Validate() error {
	if len(ts.Fields) == 0 {
		return InvalidTimestampError(fmt.Errorf("no fields specified for update"))
	} else if ts.Origin == nil {
		return InvalidTimestampError(fmt.Errorf("origin date must be specified"))
	}

	for i, fact := range ts.Factor {
		if fact.Delta == 0 {
			ts.Factor = append(ts.Factor[:i], ts.Factor[i+1:]...)
		} else if _, ok := validIntervals[fact.Interval]; !ok {
			return BadIntervalError(fmt.Errorf("invalid interval: '%s'", fact.Interval))
		}
	}

	return nil
}

func (ts *Timestamp) UpdateString() (string, error) {
	if err := ts.Validate(); err != nil {
		return "", err
	}

	temp := []string{}
	for _, fact := range ts.Factor {
		temp = append(temp, fmt.Sprintf("+ interval '%d %s'", fact.Delta, fact.Interval))
	}

	eq := fmt.Sprintf(" = timestamp '%s' %s",
		ts.Origin.Format(time.RFC3339),
		append([]byte{}, strings.Join(temp, " ")...),
	)

	return strings.Join(ts.Fields, string(append([]byte(eq), ", "...))) + eq, nil
}
