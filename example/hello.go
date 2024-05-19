package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/tinyrange/multiprogress"
)

func main() {
	prog := multiprogress.NewProgress(100, "hello")
	prog2 := multiprogress.NewProgress(500, "hello 2")

	render := multiprogress.NewTerminalRenderer(os.Stdout, multiprogress.ArrayRenderer{
		multiprogress.StringRenderer("Hello"),
		multiprogress.StringRenderer("World"),
		multiprogress.ArrayRenderer{
			multiprogress.StringRenderer("Testing 1"),
			prog,
			multiprogress.ArrayRenderer{
				multiprogress.StringRenderer("Testing 2"),
				multiprogress.ArrayRenderer{
					prog2,
				},
			},
		},
	})

	if err := render.Start(); err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer prog.Close()
		defer wg.Done()

		for i := 0; i < 100; i++ {
			prog.Add(1)

			time.Sleep(10 * time.Millisecond)
		}
	}()

	wg.Add(1)
	go func() {
		defer prog2.Close()
		defer wg.Done()

		for i := 0; i < 500; i++ {
			prog2.Add(1)

			time.Sleep(5 * time.Millisecond)
		}
	}()

	wg.Wait()

	time.Sleep(1 * time.Second)
}
