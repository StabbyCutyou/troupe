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
	for i := 0; i < 1000; i++ {
		args := message.StuffHappenedEvent{
			When:  time.Now(),
			Stuff: fmt.Sprintf("Stuff, yo %d", i),
		}
		var b *bool
		client.Go("Submit.ImportantEvent", args, &b, nil)
		//if err = client.Call("Submit.ImportantEvent", args, &b); err != nil {
		//	log.Fatal("call reply and error:", b, err)
		//}
	}
}
