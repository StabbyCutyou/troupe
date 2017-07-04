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

	_ "net/http/pprof"
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
	for err != nil {
		err = s.T.Assign(work)
	}
	return nil
}

func main() {
	cfg := troupe.Config{Mode: troupe.Fixed, MailboxSize: 100, Min: 0, Initial: 0, Max: 100}
	t, err := troupe.NewTroupe(cfg)
	if err != nil {
		log.Fatal(err)
	}
	s := &Submit{
		T: t,
	}
	rpc.Register(s)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":4488")
	if err != nil {
		log.Fatal(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go http.Serve(l, nil)
	go func() {
		log.Println(http.ListenAndServe(":4489", nil))
	}()
	<-sigc
}
