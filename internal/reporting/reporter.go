package reporting

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	messages "github.com/cucumber/messages/go/v21"
)

var globalReporter *topLevelReporter

type Reporter interface {
	Report(envelope *messages.Envelope)
}

type topLevelReporter struct {
	mutex sync.Mutex
	f     *os.File
}

type reporter struct {
	*topLevelReporter
	ch     chan *messages.Envelope
	writer io.Writer
}

func (r *reporter) Report(envelope *messages.Envelope) {
	r.ch <- envelope
}

var _ Reporter = &reporter{}

func GetReporter(t *testing.T) Reporter {
	t.Helper()

	if globalReporter == nil {
		if !initReporter(t) {
			return nil
		}
	}

	ch := make(chan *messages.Envelope, 1024)
	writer := &bytes.Buffer{}
	r := &reporter{
		topLevelReporter: globalReporter,
		ch:               ch,
		writer:           writer,
	}

	done := make(chan bool)

	go func() {
		enc := json.NewEncoder(writer)
		for env := range ch {
			err := enc.Encode(env)
			if err != nil {
				panic(err)
			}
		}
		done <- true
	}()

	t.Cleanup(func() {
		ch := r.ch
		r.ch = nil
		close(ch)
		<-done

		globalReporter.mutex.Lock()
		defer globalReporter.mutex.Unlock()
		_, err := globalReporter.f.Write(writer.Bytes())
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		// we don't know if this is the last test run, but we write TestRunFinished
		// anyway, then we move the position of the writer to right before the
		// TestRunFinished message using Seek so that a future test will overwrite
		// this message, and the last test will have the correct value for all
		// test suites
		timestamp := messages.GoTimeToTimestamp(time.Now())
		bz, err := json.Marshal(&messages.Envelope{TestRunFinished: &messages.TestRunFinished{
			Success:   !t.Failed(),
			Timestamp: &timestamp,
		}})
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		_, err = globalReporter.f.Write(bz)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		const nl = "\n"
		_, err = globalReporter.f.Write([]byte(nl))
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		offset := int64(len(bz) + len(nl))
		_, err = globalReporter.f.Seek(-offset, 1)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
	})

	return r
}

func initReporter(t *testing.T) bool {
	t.Helper()

	writer := getWriter()
	if writer == nil {
		return false
	}

	globalReporter = &topLevelReporter{
		mutex: sync.Mutex{},
		f:     writer,
	}
	enc := json.NewEncoder(writer)
	err := enc.Encode(&messages.Envelope{Meta: &messages.Meta{
		ProtocolVersion: "",
		Implementation:  nil,
		Runtime:         nil,
		Os:              nil,
		Cpu:             nil,
		Ci:              nil,
	}})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	timestamp := messages.GoTimeToTimestamp(time.Now())
	err = enc.Encode(&messages.Envelope{TestRunStarted: &messages.TestRunStarted{Timestamp: &timestamp}})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	return true
}
