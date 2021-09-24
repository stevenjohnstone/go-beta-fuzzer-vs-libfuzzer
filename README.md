# Comparing Go Fuzz Beta with Libfuzzer

The Go [native fuzzing beta](https://blog.golang.org/fuzz-beta) uses instrumentation
which is identical to the "libfuzzer" build mode. This allows direct comparison of
the mutation engines of the beta fuzzer and libfuzzer.
 

## Technical Details

When using libfuzzer, integer comparison feedback is [wired up](https://golang.org/src/runtime/libfuzzer.go)
which gives it a slight edge over the beta fuzzer. In the simple test case here, using
this functionality is disabled (with the -use_cmp=0 flag) to make a level playing field.

https://github.com/stevenjohnstone/go114-fuzz-build is used to build an archive with libfuzzer instrumentation and
a ```LLVMFuzzerTestOneInput``` harness. This is a branch of https://github.com/mdempsky/go114-fuzz-build with a
command line flag added to specify the go compiler.

# Results

Note: applies to go version devel go1.18-d413908 Fri Sep 24 07:22:13 2021 +0000 linux/amd64

* Libfuzzer works well.
* Beta-fuzzer does not: coverage guidance appears to be broken


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


## Example Runs

To run libfuzzer
```
$ docker run --rm fuzztests run libfuzzer Fuzz
INFO: Seed: 3906051938
INFO: 64 Extra Counters
INFO: -max_len is not provided; libFuzzer will not generate inputs larger than 4096 bytes
INFO: A corpus is not provided, starting from an empty corpus
#2	INITED ft: 2 corp: 1/1b lim: 4 exec/s: 0 rss: 25Mb
#16	NEW    ft: 3 corp: 2/5b lim: 4 exec/s: 0 rss: 25Mb L: 4/4 MS: 4 ShuffleBytes-ShuffleBytes-CopyPart-CopyPart-
#1646	NEW    ft: 4 corp: 3/9b lim: 4 exec/s: 0 rss: 25Mb L: 4/4 MS: 5 ChangeBinInt-CopyPart-InsertByte-EraseBytes-CrossOver-
#2207	NEW    ft: 5 corp: 4/13b lim: 4 exec/s: 0 rss: 25Mb L: 4/4 MS: 1 ChangeBit-
#5074	NEW    ft: 6 corp: 5/17b lim: 6 exec/s: 0 rss: 25Mb L: 4/4 MS: 2 ChangeBit-ChangeBit-
panic: ([]uint8) 0xc00000e018

goroutine 17 [running, locked to thread]:
github.com/stevenjohnstone/fuzztests.Fuzz(...)
	github.com/stevenjohnstone/fuzztests/fuzz.go:7
main.LLVMFuzzerTestOneInput(...)
	./main.072806297.go:21
==666== ERROR: libFuzzer: deadly signal
    #0 0x450ddf in __sanitizer_print_stack_trace (/fuzztests/fuzz.libfuzzer+0x450ddf)
    #1 0x430f4b in fuzzer::PrintStackTrace() (/fuzztests/fuzz.libfuzzer+0x430f4b)
    #2 0x414b7b in fuzzer::Fuzzer::CrashCallback() (/fuzztests/fuzz.libfuzzer+0x414b7b)
    #3 0x414b3f in fuzzer::Fuzzer::StaticCrashSignalCallback() (/fuzztests/fuzz.libfuzzer+0x414b3f)
    #4 0x7f4ed7ba972f  (/lib/x86_64-linux-gnu/libpthread.so.0+0x1272f)
    #5 0x4a53e0 in runtime.raise.abi0 runtime/sys_linux_amd64.s:164

NOTE: libFuzzer has rudimentary signal handlers.
      Combine libFuzzer with AddressSanitizer or similar for better crash reports.
SUMMARY: libFuzzer: deadly signal
MS: 1 ChangeBinInt-; base unit: 5d8fd21c04eea521d1e9348e9a6a49783365d971
0x1,0x3,0x3,0x7,
\x01\x03\x03\x07
Base64: AQMDBw==
stat::number_of_executed_units: 11365
stat::average_exec_per_sec:     11365
stat::new_units_added:          4
stat::slowest_unit_time_sec:    0
stat::peak_rss_mb:              27
```
```
$ docker run --rm fuzztests run betafuzzer Fuzz
...
fuzz: elapsed: 8m15s, execs: 50573812 (102169/sec), interesting: 1
fuzz: elapsed: 8m18s, execs: 50879304 (102167/sec), interesting: 1
fuzz: elapsed: 8m21s, execs: 51147152 (102089/sec), interesting: 1
fuzz: elapsed: 8m24s, execs: 51401940 (101988/sec), interesting: 1
fuzz: elapsed: 8m27s, execs: 51669641 (101912/sec), interesting: 1
fuzz: elapsed: 8m30s, execs: 51900502 (101766/sec), interesting: 1
fuzz: elapsed: 8m33s, execs: 52206437 (101767/sec), interesting: 1
```
and I gave up waiting there. This should be quick.

## Looping and the Beta Fuzzer

```
docker run --rm fuzztests run betafuzzer FuzzLoop
warning: starting with empty corpus
fuzz: elapsed: 0s, execs: 0 (0/sec), interesting: 0
fuzz: elapsed: 3s, execs: 243023 (80976/sec), interesting: 1
fuzz: elapsed: 6s, execs: 489758 (81557/sec), interesting: 1
fuzz: elapsed: 9s, execs: 756432 (84030/sec), interesting: 1
fuzz: elapsed: 12s, execs: 1100585 (91699/sec), interesting: 1
fuzz: elapsed: 15s, execs: 1446396 (96418/sec), interesting: 1
fuzz: elapsed: 18s, execs: 1787110 (99276/sec), interesting: 1
fuzz: elapsed: 21s, execs: 2124481 (101130/sec), interesting: 1
fuzz: elapsed: 24s, execs: 2466955 (102785/sec), interesting: 1

...

fuzz: elapsed: 5m51s, execs: 36472340 (103909/sec), interesting: 1
fuzz: elapsed: 5m54s, execs: 36794235 (103938/sec), interesting: 1
fuzz: elapsed: 5m57s, execs: 37107921 (103941/sec), interesting: 1
fuzz: elapsed: 6m0s, execs: 37446922 (104019/sec), interesting: 1
fuzz: elapsed: 6m3s, execs: 37768433 (104045/sec), interesting: 1
fuzz: elapsed: 6m6s, execs: 38097014 (104090/sec), interesting: 1
fuzz: elapsed: 6m9s, execs: 38426799 (104137/sec), interesting: 1
fuzz: elapsed: 6m12s, execs: 38745656 (104154/sec), interesting: 1
fuzz: elapsed: 6m15s, execs: 39059822 (104158/sec), interesting: 1
fuzz: elapsed: 6m18s, execs: 39362111 (104132/sec), interesting: 1
fuzz: elapsed: 6m21s, execs: 39662073 (104099/sec), interesting: 1
fuzz: elapsed: 6m24s, execs: 39965336 (104076/sec), interesting: 1
fuzz: elapsed: 6m27s, execs: 40274302 (104068/sec), interesting: 1
fuzz: elapsed: 6m30s, execs: 40574478 (104036/sec), interesting: 1
fuzz: elapsed: 6m33s, execs: 40868799 (103991/sec), interesting: 1
fuzz: elapsed: 6m36s, execs: 41182369 (103995/sec), interesting: 1
fuzz: elapsed: 6m39s, execs: 41476751 (103950/sec), interesting: 1
fuzz: elapsed: 6m42s, execs: 41788795 (103951/sec), interesting: 1
fuzz: elapsed: 6m45s, execs: 42113736 (103984/sec), interesting: 1
fuzz: elapsed: 6m48s, execs: 42426033 (103985/sec), interesting: 1
fuzz: elapsed: 6m51s, execs: 42740842 (103991/sec), interesting: 1
fuzz: elapsed: 6m54s, execs: 43033277 (103945/sec), interesting: 1
fuzz: elapsed: 6m57s, execs: 43329482 (103907/sec), interesting: 1
fuzz: elapsed: 7m0s, execs: 43608474 (103829/sec), interesting: 1
```

I gave up after 7 minutes.

Libfuzzer finds a crasher quickly
```
docker run --rm fuzztests run libfuzzer FuzzLoop 
INFO: Seed: 40387679
INFO: 66 Extra Counters
INFO: -max_len is not provided; libFuzzer will not generate inputs larger than 4096 bytes
INFO: A corpus is not provided, starting from an empty corpus
#2	INITED ft: 4 corp: 1/1b lim: 4 exec/s: 0 rss: 25Mb
#22	NEW    ft: 7 corp: 2/5b lim: 4 exec/s: 0 rss: 25Mb L: 4/4 MS: 5 InsertByte-ChangeBinInt-InsertByte-ChangeByte-CopyPart-
#1789	NEW    ft: 9 corp: 3/9b lim: 4 exec/s: 0 rss: 25Mb L: 4/4 MS: 2 ChangeBit-CrossOver-
#5226	NEW    ft: 11 corp: 4/13b lim: 6 exec/s: 0 rss: 25Mb L: 4/4 MS: 2 ChangeBit-ChangeBit-
#8412	NEW    ft: 13 corp: 5/17b lim: 8 exec/s: 8412 rss: 25Mb L: 4/4 MS: 1 CopyPart-
panic: ([]uint8) 0xc00000e018

goroutine 17 [running, locked to thread]:
github.com/stevenjohnstone/fuzztests.FuzzLoopLibFuzzer(...)
	github.com/stevenjohnstone/fuzztests/fuzz.go:12
main.LLVMFuzzerTestOneInput(...)
	./main.351689957.go:21
==557== ERROR: libFuzzer: deadly signal
    #0 0x450ddf in __sanitizer_print_stack_trace (/fuzztests/fuzz.libfuzzer+0x450ddf)
    #1 0x430f4b in fuzzer::PrintStackTrace() (/fuzztests/fuzz.libfuzzer+0x430f4b)
    #2 0x414b7b in fuzzer::Fuzzer::CrashCallback() (/fuzztests/fuzz.libfuzzer+0x414b7b)
    #3 0x414b3f in fuzzer::Fuzzer::StaticCrashSignalCallback() (/fuzztests/fuzz.libfuzzer+0x414b3f)
    #4 0x7f28163c072f  (/lib/x86_64-linux-gnu/libpthread.so.0+0x1272f)
    #5 0x4a3fc0 in runtime.raise.abi0 runtime/sys_linux_amd64.s:164

NOTE: libFuzzer has rudimentary signal handlers.
      Combine libFuzzer with AddressSanitizer or similar for better crash reports.
SUMMARY: libFuzzer: deadly signal
MS: 1 ChangeBinInt-; base unit: 625977c37f48266f883271647d56469e6067b944
0x1,0x3,0x3,0x7,
\x01\x03\x03\x07
artifact_prefix='./'; Test unit written to ./crash-f45be6129befa590730da3f100eebb7217d6b1a0
Base64: AQMDBw==
stat::number_of_executed_units: 11238
stat::average_exec_per_sec:     11238
stat::new_units_added:          4
stat::slowest_unit_time_sec:    0
stat::peak_rss_mb:              27

time elapsed 9.928340415s

```




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
