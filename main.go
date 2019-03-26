package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	args := Args{}

	flag.StringVar(&args.input,"input", "./input.txt", "The path to the input file")
	flag.StringVar(&args.output, "output", "./output.txt", "The path the output file")
	flag.StringVar(&args.sep, "sep", ",", "The separator between the source and destination in the input file")

	flag.Parse()

	args.Parse()
	args.Validate()

	fmt.Println("Input:", args.input)
	fmt.Println("Output:", args.output)
	fmt.Println("Separator:", args.sep)

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

	s := bufio.NewScanner(f)
	for s.Scan() {
		fmt.Println(s.Text())
	}
	if err := s.Err(); err != nil {
		return err
	}

	return nil
}
