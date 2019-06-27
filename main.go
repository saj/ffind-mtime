package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s <path> <age>\n", filepath.Base(os.Args[0]))
}

type Predicate func(os.FileInfo) (bool, error)

func NewOldPredicate(age time.Duration) Predicate {
	cutoff := time.Now().Add(-1 * age)
	return func(info os.FileInfo) (bool, error) {
		return !info.ModTime().After(cutoff), nil
	}
}

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(2)
	}

	var (
		path = os.Args[1]
		age  = os.Args[2]
	)

	var pred Predicate
	{
		d, err := ParseDuration(age)
		if err != nil {
			die(fmt.Errorf("age: %s", err))
		}
		pred = NewOldPredicate(d)
	}

	ok, err := search(path, pred)
	if err != nil {
		die(err)
	}
	if !ok {
		os.Exit(10)
	}
}

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

var predicateFailure = errors.New("predicate failure")

func search(root string, pred Predicate) (bool, error) {
	err := filepath.Walk(root, (&walker{pred: pred}).fn)
	if err != nil {
		if err == predicateFailure {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type walker struct {
	pred Predicate
}

func (w *walker) fn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	ok, err := w.pred(info)
	if err != nil {
		return err
	}
	if !ok {
		return predicateFailure
	}
	return nil
}
