package gocuke_test

import (
	"testing"

	"github.com/regen-network/gocuke"
)

// TestStrict is a test that should pass (as skipped) when not in strict mode
// and should fail in strict mode (go test ./... -gocuke.strict).
func TestStrict(t *testing.T) {
	gocuke.NewRunner(t, &strictSuite{}).Path("features/simple.feature").Run()
}

type strictSuite struct{}

func (s *strictSuite) IHaveCukes(int64) {
	panic("PENDING")
}

func (s *strictSuite) IEat(int64) {
	panic("PENDING")
}

func (s *strictSuite) IHaveLeft(int64) {
	panic("PENDING")
}
