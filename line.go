package main

import (
	"github.com/pkg/errors"
	"strings"
)

type Line struct {
	raw  string
	from string
	to   string
	args *Args
}

func (l *Line) Parse() error {
	parts := strings.Split(l.raw, l.args.sep)
	if len(parts) != 2 {
		return errors.New("Invalid line")
	}
	l.from = parts[0]
	l.to = parts[1]

	return nil
}