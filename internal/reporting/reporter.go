package reporting

import (
	"encoding/json"
	"github.com/cucumber/messages-go/v16"
	"sync/atomic"
	"testing"
	"time"
)

var globalReporter *reporter

type Reporter interface {
	Report(envelope *messages.Envelope)
}

type reporter struct {
	ch     chan *messages.Envelope
	failed int32
}

func (r *reporter) Report(envelope *messages.Envelope) {
	r.ch <- envelope
}

func (r *reporter) registerTest(t *testing.T) {
	t.Cleanup(func() {
		if t.Failed() {
			atomic.StoreInt32(&r.failed, 1)
		}
	})
}

var _ Reporter = &reporter{}

func GetReporter(t *testing.T) Reporter {
	t.Helper()

	if globalReporter != nil {
		globalReporter.registerTest(t)
		return globalReporter
	}

	return initReporter(t)
}

func initReporter(t *testing.T) Reporter {
	t.Helper()

	writer := getWriter()
	if writer == nil {
		return nil
	}

	ch := make(chan *messages.Envelope, 1024)
	globalReporter = &reporter{ch: ch}
	go func() {
		enc := json.NewEncoder(writer)
		for {
			env, more := <-ch
			if !more {
				break
			}

			err := enc.Encode(env)
			if err != nil {
				panic(err)
			}
		}

		timestamp := messages.GoTimeToTimestamp(time.Now())
		failed := atomic.LoadInt32(&globalReporter.failed) == 1
		err := enc.Encode(&messages.Envelope{TestRunFinished: &messages.TestRunFinished{
			Timestamp: &timestamp,
			Success:   !failed,
		}})
		if err != nil {
			panic(err)
		}

		err = writer.Close()
		if err != nil {
			panic(err)
		}
	}()

	globalReporter.Report(&messages.Envelope{Meta: &messages.Meta{
		ProtocolVersion: "",
		Implementation:  nil,
		Runtime:         nil,
		Os:              nil,
		Cpu:             nil,
		Ci:              nil,
	}})
	timestamp := messages.GoTimeToTimestamp(time.Now())
	globalReporter.Report(&messages.Envelope{TestRunStarted: &messages.TestRunStarted{Timestamp: &timestamp}})
	globalReporter.registerTest(t)

	return globalReporter
}
