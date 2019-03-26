package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	args := Args{}

	flag.StringVar(&args.host,"host", "", "The hostname to prefix relative from paths with")
	flag.StringVar(&args.input,"input", "./input.txt", "The path to the input file")
	flag.StringVar(&args.output, "output", "./output.txt", "The path the output file")
	flag.StringVar(&args.sep, "sep", " ", "The separator between the source and destination in the input file")
	flag.IntVar(&args.threads, "threads", 4, "Number of links to trace at once")

	flag.Parse()
	args.Parse()
	args.Validate()

	if err := process(&args); err != nil {
		log.Fatal(err)
	}
}

func process(a *Args) (err error) {
	f, err := os.Open(a.input)
	defer func () {
		fileErr := f.Close()
		if err != nil {
			err = fileErr
		}
	}()
	if err != nil {
		return err
	}

	var lines []Line
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		l := Line{
			raw: scanner.Text(),
			args: a,
		}
		if err := l.Parse(); err != nil {
			continue // Maybe debug log here for parse errors?
		}

		lines = append(lines, l)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	queue := make(chan Line)
	out := make(chan Line)

	// Spawn some workers
	for i := 0; i < a.threads; i++ {
		wg.Add(1)
		go worker(&wg, queue, out, i)
	}

	go func() {
		for _, line := range lines {
			queue <- line
		}
		close(queue)
	}()

	go func() {
		for l := range out {
			fmt.Println(l.from, l.to)
		}
	}()

	wg.Wait()

	return nil
}

func worker(wg *sync.WaitGroup, queue <-chan Line, out chan<- Line, id int) {
	for l := range queue {
		myURL := "https://" + l.args.host + "/" + l.from
		nextURL := myURL
		var i int
		for i < 100 {
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				} }

			resp, err := client.Get(nextURL)

			if err != nil {
				fmt.Println(err)
			}

			if resp.StatusCode == 200 {
				break
			} else {
				out <- Line{
					from: myURL,
					to: resp.Header.Get("Location"),
				}
				i += 1
			}
		}
	}

	wg.Done()
}