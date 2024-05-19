package main

import (
	glog "log"
	"os"
	"time"

	"github.com/tinyrange/multiprogress"
)

func main() {
	log := multiprogress.NewLogger(multiprogress.NewRenderGroup(20))
	prog := multiprogress.NewProgress(1000, "working")

	render := multiprogress.NewTerminalRenderer(os.Stdout, 30, multiprogress.ArrayRenderer{
		log,
		prog,
	})

	if err := render.Start(); err != nil {
		glog.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		log.Info("Hello", "i", i)
		prog.Add(1)

		time.Sleep(10 * time.Millisecond)
	}
}
