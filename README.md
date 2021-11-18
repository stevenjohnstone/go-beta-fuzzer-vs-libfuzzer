# Comparing Go Native Fuzz with Libfuzzer

The Go [native fuzzing beta](https://blog.golang.org/fuzz-beta) uses instrumentation
which is identical to the "libfuzzer" build mode. This allows direct comparison of
the mutation engines of the native fuzzer and libfuzzer.
 

## Technical Details

When using libfuzzer, integer comparison feedback is [wired up](https://golang.org/src/runtime/libfuzzer.go)
which gives it a slight edge over the native fuzzer. In the simple test case here, using
this functionality is disabled (with the -use_cmp=0 flag) to make a level playing field.

https://github.com/stevenjohnstone/go114-fuzz-build is used to build an archive with libfuzzer instrumentation and
a ```LLVMFuzzerTestOneInput``` harness. This is a branch of https://github.com/mdempsky/go114-fuzz-build with a
command line flag added to specify the go compiler.

## Usage

Build the container:

```
$ docker build -t fuzztests .
```

There are four tests:

Non-looping (```magic``` in [fuzz.go](/fuzz.go))
```
docker run --rm fuzztests run libfuzzer Fuzz
docker run --rm fuzztests run betafuzzer Fuzz
```
which run the libfuzzer and beta fuzzer tests, respectively.


Looping (```loopmagic``` in [fuzz.go](/fuzz.go))
```
docker run --rm fuzztests run libfuzzer FuzzLoop
docker run --rm fuzztests run betafuzzer FuzzLoop
```

# Results

```
$ docker run --rm fuzztests run betafuzzer Fuzz
warning: starting with empty corpus
fuzz: elapsed: 0s, execs: 0 (0/sec), new interesting: 0 (total: 0)
fuzz: minimizing 41-byte failing input file
fuzz: elapsed: 1s, minimizing
--- FAIL: Fuzz (0.95s)
    --- FAIL: Fuzz (0.00s)
        fuzz_test.go:8: magic is [1 3 3 7]
    
    Failing input written to testdata/fuzz/Fuzz/43be6ff849f66da4e248cca23eeccf1df4980769c416144a1fae2d06b8764b3c
    To re-run:
    go test -run=Fuzz/43be6ff849f66da4e248cca23eeccf1df4980769c416144a1fae2d06b8764b3c
FAIL
exit status 1
FAIL	github.com/stevenjohnstone/fuzztests	0.954s

time elapsed 6.605192741s
Error: running "gotip test -fuzz=Fuzz$" failed with exit code 1

$ docker run --rm fuzztests run libfuzzer Fuzz
INFO: Seed: 2088058185
INFO: 64 Extra Counters
INFO: -max_len is not provided; libFuzzer will not generate inputs larger than 4096 bytes
INFO: A corpus is not provided, starting from an empty corpus
#2	INITED ft: 2 corp: 1/1b lim: 4 exec/s: 0 rss: 25Mb
#47	NEW    ft: 3 corp: 2/5b lim: 4 exec/s: 0 rss: 25Mb L: 4/4 MS: 5 CopyPart-ChangeByte-CopyPart-InsertByte-InsertByte-
#5957	NEW    ft: 4 corp: 3/9b lim: 8 exec/s: 0 rss: 25Mb L: 4/4 MS: 5 CrossOver-EraseBytes-ChangeBit-InsertRepeatedBytes-ChangeByte-
#10843	NEW    ft: 5 corp: 4/13b lim: 11 exec/s: 0 rss: 25Mb L: 4/4 MS: 1 ChangeBinInt-
#10989	NEW    ft: 6 corp: 5/17b lim: 11 exec/s: 0 rss: 25Mb L: 4/4 MS: 1 CopyPart-
#32768	pulse  ft: 6 corp: 5/17b lim: 29 exec/s: 16384 rss: 25Mb
panic: ([]uint8) 0xc000108000

goroutine 17 [running, locked to thread]:
github.com/stevenjohnstone/fuzztests.FuzzLibFuzzer(...)
	github.com/stevenjohnstone/fuzztests/fuzz.go:5
main.LLVMFuzzerTestOneInput(...)
	./main.877434895.go:21
==695== ERROR: libFuzzer: deadly signal
    #0 0x450ddf in __sanitizer_print_stack_trace (/fuzztests/fuzz.libfuzzer+0x450ddf)
    #1 0x430f4b in fuzzer::PrintStackTrace() (/fuzztests/fuzz.libfuzzer+0x430f4b)
    #2 0x414b7b in fuzzer::Fuzzer::CrashCallback() (/fuzztests/fuzz.libfuzzer+0x414b7b)
    #3 0x414b3f in fuzzer::Fuzzer::StaticCrashSignalCallback() (/fuzztests/fuzz.libfuzzer+0x414b3f)
    #4 0x7fd13f01172f  (/lib/x86_64-linux-gnu/libpthread.so.0+0x1272f)
    #5 0x4a4840 in runtime.raise.abi0 runtime/sys_linux_amd64.s:167

NOTE: libFuzzer has rudimentary signal handlers.
      Combine libFuzzer with AddressSanitizer or similar for better crash reports.
SUMMARY: libFuzzer: deadly signal
MS: 1 ChangeByte-; base unit: 148bde6fe8a5fda4e69bbded6d94f5341b411e6f
0x1,0x3,0x3,0x7,
\x01\x03\x03\x07
artifact_prefix='./'; Test unit written to ./crash-f45be6129befa590730da3f100eebb7217d6b1a0
Base64: AQMDBw==
stat::number_of_executed_units: 48965
stat::average_exec_per_sec:     12241
stat::new_units_added:          4
stat::slowest_unit_time_sec:    0
stat::peak_rss_mb:              27

time elapsed 10.000823393s

$ docker run --rm fuzztests run betafuzzer FuzzLoop
warning: starting with empty corpus
fuzz: elapsed: 0s, execs: 0 (0/sec), new interesting: 0 (total: 0)
fuzz: minimizing 41-byte failing input file
fuzz: elapsed: 1s, minimizing
--- FAIL: FuzzLoop (1.17s)
    --- FAIL: FuzzLoop (0.00s)
        fuzz_test.go:16: magic is [1 3 3 7]
    
    Failing input written to testdata/fuzz/FuzzLoop/43be6ff849f66da4e248cca23eeccf1df4980769c416144a1fae2d06b8764b3c
    To re-run:
    go test -run=FuzzLoop/43be6ff849f66da4e248cca23eeccf1df4980769c416144a1fae2d06b8764b3c
FAIL
exit status 1
FAIL	github.com/stevenjohnstone/fuzztests	1.171s

time elapsed 5.950228517s
Error: running "gotip test -fuzz=FuzzLoop$" failed with exit code 1

```

This shows that for these very simple, classic demostrations of coverage guided fuzzing, the native golang
fuzzer finds crashing inputs faster than libfuzzer.



## Other Targets

Other golang packages can be tested by volume mounting the code directory to `fuzztests` e.g.


```
docker run --rm -v /path/to/code/under/test:/fuzztests fuzztests betafuzzer Fuzz
```

will run the beta fuzzer for function `Fuzz` defined in the code at `/path/to/code/under/test`.


## Cross-Compilation for Raspberry Pi

To cross-compile a beta fuzzer binary for use on a Raspberry Pi

```
docker run --rm -v "$(pwd):/output" fuzztests run crosscompile output/fuzzer
```
There will no be an executable `fuzzer` in the current directory which can be copied to the target Pi system. To run the fuzzer `Fuzz`, execute
```
./fuzzer -test.fuzz=Fuzz$ -test.run=^$ -test.fuzzcachedir=./cache

```
on the Pi. Note that all fuzzer functions are available by specifiying the `-test.fuzz` parameter.


## TODO

* more comparison tests
* perform more runs to get an idea of the average executions required to complete tests
* run libfuzzer tests with [integer comparison feedback](https://llvm.org/docs/LibFuzzer.html#id32): maybe useful to add to the beta fuzzer?
* cross-compilation of libfuzzer to rpi...stuck on getting a toolchain with the ubsan/fuzzer interceptors for aarch64 (don't come packaged with clang on ubuntu)
