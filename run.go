// +build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/magefile/mage/sh"
)

func Libfuzzer() error {
	start := time.Now()
	defer func() {
		fmt.Printf("\ntime elapsed %v\n", time.Now().Sub(start))
	}()

	if err := sh.Run("go114-fuzz-build", "-go=gotip", "-func=Fuzz", "-o=fuzz.a", "./"); err != nil {
		return err
	}

	if err := sh.Run("clang", "-fsanitize=fuzzer", "fuzz.a", "-o", "fuzz.libfuzzer"); err != nil {
		return err
	}

	ran, err := sh.Exec(nil, os.Stdout, os.Stderr, "./fuzz.libfuzzer", "-use_cmp=0", "-print_final_stats=1", "-error_exitcode=0")
	if !ran {
		return errors.New("failed to run command")
	}
	return err
}

func Betafuzzer() error {
	// TODO: clean out testdata and cache
	start := time.Now()
	defer func() {
		fmt.Printf("\ntime elapsed %v\n", time.Now().Sub(start))
	}()

	ran, err := sh.Exec(nil, os.Stdout, os.Stderr, "gotip", "test", "-fuzz=FuzzBeta")
	if !ran {
		return errors.New("failed to run command")
	}
	return err
}
