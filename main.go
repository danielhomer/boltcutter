package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

func main() {
	args := Args{}

	flag.StringVar(&args.input,"input", "./input.txt", "The path to the input file")
	flag.StringVar(&args.output, "output", "./output.txt", "The path the output file")
	flag.StringVar(&args.sep, "sep", ",", "The separator between the source and destination in the input file")
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

	var wg sync.WaitGroup
	queue := make(chan string)
	done := make(chan string)

	// Spawn some workers
	for i := 0; i < a.threads; i++ {
		wg.Add(1)
		go worker(&wg, queue, done, i)
	}

	go func() {
		for _, line := range lines {
			queue <- line.raw
		}
		close(queue)
	}()

	go func() {
		for l := range done {
			fmt.Println(l)
		}
	}()

	wg.Wait()

	return nil
}

func worker(wg *sync.WaitGroup, queue chan string, done chan<- string, id int) {
	for p := range queue {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1)))
		//fmt.Printf("[Worker %v] processing %s\n", id, p)
		done <- p
	}

	wg.Done()
}