package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	r "gopkg.in/dancannon/gorethink.v2"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	session, _ := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})

	// for i := 0; i < 1000; i++ {
	// 	r.DB("test").TableCreate(fmt.Sprintf("table_create_test_%04d", i)).Run(session)
	// }

	for i := 0; i < 1000; i++ {
		res, _ := r.DB("test").Table(fmt.Sprintf("table_create_test_%04d", i)).Changes().Run(session)
		go func(res *r.Cursor, i int) {
			var value interface{}
			for res.Next(&value) {
				fmt.Printf("changefeed from table_create_test_%04d\n", i)
				fmt.Println(value)
			}
		}(res, i)
	}

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
