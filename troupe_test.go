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

var twentymillis Work = func() error {
	time.Sleep(20 * time.Microsecond)
	return nil
}

var testCases = []testCase{
	{title: "1w:1a:20ms", work: twentymillis, cfg: Config{MailboxSize: 1, Min: 0, Initial: 0, Max: 1}},
	{title: "10w:1a:20ms", work: twentymillis, cfg: Config{MailboxSize: 10, Min: 0, Initial: 0, Max: 1}},
	{title: "100w:1a:20ms", work: twentymillis, cfg: Config{MailboxSize: 100, Min: 0, Initial: 0, Max: 1}},

	{title: "1w:10a:20ms", work: twentymillis, cfg: Config{MailboxSize: 1, Min: 0, Initial: 0, Max: 10}},
	{title: "10w:10a:20ms", work: twentymillis, cfg: Config{MailboxSize: 10, Min: 0, Initial: 0, Max: 10}},
	{title: "100w:10a:20ms", work: twentymillis, cfg: Config{MailboxSize: 100, Min: 0, Initial: 0, Max: 10}},

	{title: "1w:100a:20ms", work: twentymillis, cfg: Config{MailboxSize: 1, Min: 0, Initial: 0, Max: 100}},
	{title: "10w:100a:20ms", work: twentymillis, cfg: Config{MailboxSize: 10, Min: 0, Initial: 0, Max: 100}},
	{title: "100w:100a:20ms", work: twentymillis, cfg: Config{MailboxSize: 100, Min: 0, Initial: 0, Max: 100}},

	{title: "1w:1ka:20ms", work: twentymillis, cfg: Config{MailboxSize: 1, Min: 0, Initial: 0, Max: 1000}},
	{title: "10w:1ka:20ms", work: twentymillis, cfg: Config{MailboxSize: 10, Min: 0, Initial: 0, Max: 1000}},
	{title: "100w:1ka:20ms", work: twentymillis, cfg: Config{MailboxSize: 100, Min: 0, Initial: 0, Max: 1000}},

	{title: "1w:10ka:20ms", work: twentymillis, cfg: Config{MailboxSize: 1, Min: 0, Initial: 0, Max: 10000}},
	{title: "10w:10ka:20ms", work: twentymillis, cfg: Config{MailboxSize: 10, Min: 0, Initial: 0, Max: 10000}},
	{title: "100w:10ka:20ms", work: twentymillis, cfg: Config{MailboxSize: 100, Min: 0, Initial: 0, Max: 10000}},

	{title: "1w:100ka:20ms", work: twentymillis, cfg: Config{MailboxSize: 1, Min: 0, Initial: 0, Max: 100000}},
	{title: "10w:100ka:20ms", work: twentymillis, cfg: Config{MailboxSize: 10, Min: 0, Initial: 0, Max: 100000}},
	{title: "100w:100ka:20ms", work: twentymillis, cfg: Config{MailboxSize: 100, Min: 0, Initial: 0, Max: 100000}},
}

func TestSJ(t *testing.T) {
	for _, c := range testCases {
		time.Sleep(1 * time.Second)
		t.Run(c.title, func(t *testing.T) {
			s, _ := NewTroupe(c.cfg)
			for i := 0; i < 100; i++ {
				s.Assign(c.work)
			}
			s.Shutdown()
			s.Join()
		})
	}
}

func BenchmarkT(b *testing.B) {
	for _, c := range testCases {
		c.cfg.AssignmentMode = ModePriority
		b.Run(c.title, func(b *testing.B) {
			s, _ := NewTroupe(c.cfg)
			for i := 0; i < b.N; i++ {
				err := s.Assign(c.work)
				for err != nil {
					// Attempt to assign work, and keep trying until it succeeds
					// Do not advance b.N until this message is in a queue
					err = s.Assign(c.work)
				}
			}
			s.Shutdown()
			s.Join()
		})
	}
}

func BenchmarkTR(b *testing.B) {
	for _, c := range testCases {
		c.cfg.AssignmentMode = ModeRandom
		b.Run(c.title, func(b *testing.B) {
			s, _ := NewTroupe(c.cfg)
			for i := 0; i < b.N; i++ {
				err := s.Assign(c.work)
				for err != nil {
					// Attempt to assign work, and keep trying until it succeeds
					// Do not advance b.N until this message is in a queue
					err = s.Assign(c.work)
				}
			}
			s.Shutdown()
			s.Join()
		})
	}
}
