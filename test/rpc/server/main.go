package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/StabbyCutyou/troupe"
	"github.com/StabbyCutyou/troupe/test/rpc/message"
)

// Submit submits
type Submit struct {
	T *troupe.Troupe
}

// ImportantEvent logs important events
func (s *Submit) ImportantEvent(w *message.StuffHappenedEvent, b *bool) error {
	work := func() error {
		fmt.Printf("%+v\n", w)
		return nil
	}
	err := s.T.Assign(work)
	if err != nil {
		err = s.T.Assign(work)
	}
	return nil
}

func main() {
	cfg := troupe.Config{Mode: troupe.Fixed, MailboxSize: 1, Min: 0, Initial: 0, Max: 10000}
	t, err := troupe.NewTroupe(cfg)
	if err != nil {
		log.Fatal(err)
	}
	s := &Submit{
		T: t,
	}
	rpc.Register(s)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":4488")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go http.Serve(l, nil)
	<-sigc
}
