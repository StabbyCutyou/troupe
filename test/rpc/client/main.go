package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"

	"github.com/StabbyCutyou/troupe/test/rpc/message"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "0.0.0.0:4488")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	i := 0
	for {
		i++
		args := message.StuffHappenedEvent{
			When:  time.Now(),
			Stuff: fmt.Sprintf("Stuff, yo %d", i),
		}
		var b *bool
		// We don't really care about the reply, just that we make the call
		client.Go("Submit.ImportantEvent", args, &b, nil)
	}
}
