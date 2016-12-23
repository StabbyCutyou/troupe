package troupe

import (
	"testing"
	"time"
)

type testCase struct {
	title string
	cfg   Config
	work  Work
}

var onenano Work = func() error {
	time.Sleep(1 * time.Nanosecond)
	return nil
}

var onemilli Work = func() error {
	time.Sleep(1 * time.Millisecond)
	return nil
}

var testCases = []testCase{
	{
		title: "1 actor - 1 nano", work: onenano, cfg: Config{Min: 0, Initial: 0, Max: 1},
	},
	{
		title: "10 actors - 1 nano", work: onenano, cfg: Config{Min: 0, Initial: 0, Max: 10},
	},
	{
		title: "100 actors - 1 nano", work: onenano, cfg: Config{Min: 0, Initial: 0, Max: 100},
	},
	{
		title: "1000 actors - 1 nano", work: onenano, cfg: Config{Min: 0, Initial: 0, Max: 1000},
	},
	{
		title: "10000 actors - 1 nano", work: onenano, cfg: Config{Min: 0, Initial: 0, Max: 10000},
	},
	{
		title: "100000 actors - 1 nano", work: onenano, cfg: Config{Min: 0, Initial: 0, Max: 100000},
	},
	{
		title: "1 actor - 1 milli", work: onemilli, cfg: Config{Min: 0, Initial: 0, Max: 1},
	},
	{
		title: "10 actors - 1 milli", work: onemilli, cfg: Config{Min: 0, Initial: 0, Max: 10},
	},
	{
		title: "100 actors - 1 milli", work: onemilli, cfg: Config{Min: 0, Initial: 0, Max: 100},
	},
	{
		title: "1000 actors - 1 milli", work: onemilli, cfg: Config{Min: 0, Initial: 0, Max: 1000},
	},
	{
		title: "10000 actors - 1 milli", work: onemilli, cfg: Config{Min: 0, Initial: 0, Max: 10000},
	},
	{
		title: "100000 actors - 1 milli", work: onemilli, cfg: Config{Min: 0, Initial: 0, Max: 100000},
	},
}

func BenchmarkTroupe(b *testing.B) {
	for _, c := range testCases {
		b.Run(c.title, func(b *testing.B) {
			s, _ := NewTroupe(c.cfg)
			for i := 0; i < b.N; i++ {
				s.Assign(c.work)
			}
			s.Shutdown()
		})
	}
}
