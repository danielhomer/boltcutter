package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
	fmt.Println("Seperator:", args.sep)

	if err := process(args.input, args.output, args.sep); err != nil {
		log.Fatal(err)
	}
}

func process(input string, output string, sep string) error {
	data, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	fmt.Print(string(data))

	return nil
}
