// +build mage

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/magefile/mage/sh"
)

// Libfuzzer builds and runs a libfuzzer harness for the provided function `fuzzFunc`
func Libfuzzer(_ context.Context, fuzzFunc string) error {
	start := time.Now()
	defer func() {
		fmt.Printf("\ntime elapsed %v\n", time.Now().Sub(start))
	}()

	if err := sh.Run("go114-fuzz-build", "-go=gotip", fmt.Sprintf("-func=%s", fuzzFunc), "-o=fuzz.a", "./"); err != nil {
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

// Betafuzzer runs the golang native fuzzer for `fuzzFunc`
func Betafuzzer(_ context.Context, fuzzFunc string) error {
	start := time.Now()
	defer func() {
		fmt.Printf("\ntime elapsed %v\n", time.Now().Sub(start))
	}()
	fuzzFunc += "$"

	ran, err := sh.Exec(nil, os.Stdout, os.Stderr, "gotip", "test", fmt.Sprintf("-fuzz=%s", fuzzFunc))
	if !ran {
		return errors.New("failed to run command")
	}
	return err
}
