package main

type Line struct {
	raw  string
	from string
	to   string
	args *Args
}

func (l *Line) Parse() error {
	return nil
}