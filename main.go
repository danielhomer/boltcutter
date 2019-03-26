package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	args := Args{}

	flag.StringVar(&args.input,"input", "./input.txt", "The path to the input file")
	flag.StringVar(&args.output, "output", "./output.txt", "The path the output file")
	flag.StringVar(&args.sep, "sep", ",", "The separator between the source and destination in the input file")

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
	var wg sync.WaitGroup
	res := make(chan Line)

	defer wg.Wait()

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

	wg.Add(len(lines))

	for _, line := range lines {
		go func() {
			defer wg.Done()
			res <- line
		}()
	}

	go func() {
		for line := range res {
			fmt.Println(line.raw)
		}
	}()

	return nil
}
