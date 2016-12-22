package troupe_test

import (
	"testing"
	"time"

	"github.com/StabbyCutyou/troupe"
)

func BenchmarkTroupe1(b *testing.B) {
	s, _ := troupe.NewTroupe(troupe.Config{
		Min:     0,
		Initial: 0,
		Max:     1,
	})

	for i := 0; i < b.N; i++ {
		s.Assign(func() error {
			time.Sleep(1 * time.Nanosecond)
			return nil
		})
	}
	s.Shutdown()
	s = nil
}

func BenchmarkTroupe10(b *testing.B) {
	s, _ := troupe.NewTroupe(troupe.Config{
		Min:     0,
		Initial: 0,
		Max:     10,
	})

	for i := 0; i < b.N; i++ {
		s.Assign(func() error {
			time.Sleep(1 * time.Nanosecond)
			return nil
		})
	}
	s.Shutdown()
}

func BenchmarkTroupe100(b *testing.B) {
	s, _ := troupe.NewTroupe(troupe.Config{
		Min:     0,
		Initial: 0,
		Max:     100,
	})

	for i := 0; i < b.N; i++ {
		s.Assign(func() error {
			time.Sleep(1 * time.Nanosecond)
			return nil
		})
	}
	s.Shutdown()
}

func BenchmarkTroupe1000(b *testing.B) {
	s, _ := troupe.NewTroupe(troupe.Config{
		Min:     0,
		Initial: 0,
		Max:     1000,
	})

	for i := 0; i < b.N; i++ {
		s.Assign(func() error {
			time.Sleep(1 * time.Nanosecond)
			return nil
		})
	}
	s.Shutdown()
}

func BenchmarkTroupe10000(b *testing.B) {
	s, _ := troupe.NewTroupe(troupe.Config{
		Min:     0,
		Initial: 0,
		Max:     10000,
	})

	for i := 0; i < b.N; i++ {
		s.Assign(func() error {
			time.Sleep(1 * time.Nanosecond)
			return nil
		})
	}
	s.Shutdown()
}

func BenchmarkTroupe100000(b *testing.B) {
	s, _ := troupe.NewTroupe(troupe.Config{
		Min:     0,
		Initial: 0,
		Max:     100000,
	})

	for i := 0; i < b.N; i++ {
		s.Assign(func() error {
			time.Sleep(1 * time.Nanosecond)
			return nil
		})
	}
	s.Shutdown()
}
