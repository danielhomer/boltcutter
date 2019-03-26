package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
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
	queue := make(chan string)
	done := make(chan string)
	kill := make(chan bool)

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


	for i := 0; i < a.threads; i++ {
		go worker(queue, i, done, kill)
	}

	for _, line := range lines {
		go func() {
			queue <- line.from
		}()
	}
	//close(queue)

	for l := range done {
		fmt.Println(l)
	}

	return nil
}

func worker(queue chan string, no int, done chan string, kill chan bool) {
	for l := range queue {
		fmt.Println("worker", no, "started  job", l)
		time.Sleep(time.Second)
		fmt.Println("worker", no, "finished job", l)
		done <- l
	}
}