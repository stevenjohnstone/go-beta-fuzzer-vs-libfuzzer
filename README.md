# Comparing Go Fuzz Beta with Libfuzzer

The Go [native fuzzing beta](https://blog.golang.org/fuzz-beta) uses instrumentation
which is identical to the "libfuzzer" build mode. This allows direct comparison of
the mutation engines of the beta fuzzer and libfuzzer.

Note: applies to go version devel go1.17-7b6893d2d2 Mon Aug 23 20:58:28 2021 +0000 linux/amd64

## Technical Details

When using libfuzzer, integer comparison feedback is [wired up](https://golang.org/src/runtime/libfuzzer.go)
which gives it a slight edge over the beta fuzzer. In the simple test case here, using
the functionality is disabled (with the -use_cmp=0 flag).

https://github.com/stevenjohnstone/go114-fuzz-build is used to build an archive with libfuzzer instrumentation and
a ```LLVMFuzzerTestOneInput``` harness. This is a branch of https://github.com/mdempsky/go114-fuzz-build with a
command line flag added to specify the go compiler. Here, the dev.fuzz branch of golang is used with gotip.


## Usage

Build the container:

```
$ docker build -t fuzztests .
```

There are four tests:

Non-looping (```magic``` in [fuzz.go](/fuzz.go))
```
docker run --rm fuzztests mage libfuzzer
docker run --rm fuzztests mage betafuzzer
```
which run the libfuzzer and beta fuzzer tests, respectively.


Looping (```loopmagic``` in [fuzz.go](/fuzz.go))
```
docker run --rm -e FUZZ_FUNC=FuzzLoop fuzztests mage libfuzzer
docker run --rm -e FUZZ_FUNC=FuzzLoop fuzztests mage betafuzzer
```


## Example Results

To run libfuzzer
```
$ docker run --rm fuzztests mage libfuzzer
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
$ docker run --rm fuzztests mage betafuzzer
fuzzing, elapsed: 3.0s, execs: 330708 (110206/sec), workers: 8, interesting: 3
fuzzing, elapsed: 6.0s, execs: 679284 (113198/sec), workers: 8, interesting: 3
fuzzing, elapsed: 9.0s, execs: 1021705 (113508/sec), workers: 8, interesting: 3
fuzzing, elapsed: 12.0s, execs: 1366081 (113798/sec), workers: 8, interesting: 3
fuzzing, elapsed: 15.0s, execs: 1713104 (114173/sec), workers: 8, interesting: 3
fuzzing, elapsed: 18.0s, execs: 2050969 (113939/sec), workers: 8, interesting: 3
fuzzing, elapsed: 21.0s, execs: 2405213 (114510/sec), workers: 8, interesting: 3
fuzzing, elapsed: 24.0s, execs: 2757696 (114883/sec), workers: 8, interesting: 3
fuzzing, elapsed: 27.0s, execs: 3107932 (115103/sec), workers: 8, interesting: 3
fuzzing, elapsed: 30.0s, execs: 3456089 (115193/sec), workers: 8, interesting: 3
fuzzing, elapsed: 33.0s, execs: 3801287 (115188/sec), workers: 8, interesting: 3
fuzzing, elapsed: 36.0s, execs: 4143806 (115102/sec), workers: 8, interesting: 3
fuzzing, elapsed: 39.0s, execs: 4487768 (115068/sec), workers: 8, interesting: 3
fuzzing, elapsed: 42.0s, execs: 4833030 (115058/sec), workers: 8, interesting: 4
fuzzing, elapsed: 45.0s, execs: 5186331 (115250/sec), workers: 8, interesting: 4
fuzzing, elapsed: 48.0s, execs: 5537571 (115364/sec), workers: 8, interesting: 4
fuzzing, elapsed: 51.0s, execs: 5879385 (115280/sec), workers: 8, interesting: 4
fuzzing, elapsed: 54.0s, execs: 6222607 (115232/sec), workers: 8, interesting: 4
fuzzing, elapsed: 57.0s, execs: 6565687 (115186/sec), workers: 8, interesting: 5
fuzzing, elapsed: 60.0s, execs: 6907395 (115120/sec), workers: 8, interesting: 5
fuzzing, elapsed: 63.0s, execs: 7247749 (115042/sec), workers: 8, interesting: 5
fuzzing, elapsed: 66.0s, execs: 7588275 (114973/sec), workers: 8, interesting: 5
fuzzing, elapsed: 69.0s, execs: 7936124 (115009/sec), workers: 8, interesting: 5
fuzzing, elapsed: 72.0s, execs: 8283781 (115050/sec), workers: 8, interesting: 5
fuzzing, elapsed: 75.0s, execs: 8631168 (115081/sec), workers: 8, interesting: 5
fuzzing, elapsed: 78.0s, execs: 8975726 (115067/sec), workers: 8, interesting: 5
fuzzing, elapsed: 81.0s, execs: 9329377 (115176/sec), workers: 8, interesting: 5
found a crash, minimizing...
fuzzing, elapsed: 82.3s, execs: 9472196 (115032/sec), workers: 8, interesting: 5
--- FAIL: Fuzz (82.34s)
        --- FAIL: Fuzz (0.00s)
            fuzz_test.go:10: magic is [1 3 3 7]
    
    Crash written to testdata/corpus/Fuzz/1757d23fd60223bd5a11cfd3a7978f28cdb2b98e0b81542690f8f75ba96d043d
    To re-run:
    go test github.com/stevenjohnstone/fuzztests -run=Fuzz/1757d23fd60223bd5a11cfd3a7978f28cdb2b98e0b81542690f8f75ba96d043d
FAIL
exit status 1
FAIL	github.com/stevenjohnstone/fuzztests	82.391s
Error: running "gotip test -fuzz=Fuzz$" failed with exit code 1
```

## Looping and the Beta Fuzzer
Finding crashers with a simple loop appears to be about similar in performance between libfuzzer and beta fuzzer.

```
docker run --rm -e FUZZ_FUNC=FuzzLoop fuzztests mage betafuzzer
fuzzing, elapsed: 3.0s, execs: 334297 (111390/sec), workers: 8, interesting: 3
fuzzing, elapsed: 6.0s, execs: 688004 (114632/sec), workers: 8, interesting: 4
fuzzing, elapsed: 9.0s, execs: 1037669 (115277/sec), workers: 8, interesting: 4
fuzzing, elapsed: 12.0s, execs: 1373068 (114408/sec), workers: 8, interesting: 4
fuzzing, elapsed: 15.0s, execs: 1722606 (114834/sec), workers: 8, interesting: 5
found a crash, minimizing...
fuzzing, elapsed: 15.4s, execs: 1756923 (114089/sec), workers: 8, interesting: 5
--- FAIL: FuzzLoop(15.40s)
        --- FAIL: FuzzLoop (0.00s)
            fuzz_test.go:21: magic is [1 3 3 7]
    
    Crash written to testdata/corpus/FuzzLoop/1757d23fd60223bd5a11cfd3a7978f28cdb2b98e0b81542690f8f75ba96d043d
    To re-run:
    go test github.com/stevenjohnstone/fuzztests -run=FuzzLoop/1757d23fd60223bd5a11cfd3a7978f28cdb2b98e0b81542690f8f75ba96d043d
FAIL
exit status 1
FAIL	github.com/stevenjohnstone/fuzztests	15.415s

```
```
docker run --rm -e FUZZ_FUNC=FuzzLoop fuzztests mage libfuzzer
INFO: Seed: 4224861379
INFO: 66 Extra Counters
INFO: -max_len is not provided; libFuzzer will not generate inputs larger than 4096 bytes
INFO: A corpus is not provided, starting from an empty corpus
#2	INITED ft: 4 corp: 1/1b lim: 4 exec/s: 0 rss: 25Mb
#22	NEW    ft: 7 corp: 2/5b lim: 4 exec/s: 0 rss: 25Mb L: 4/4 MS: 5 InsertByte-ShuffleBytes-InsertByte-EraseBytes-CopyPart-
#16314	NEW    ft: 9 corp: 3/9b lim: 17 exec/s: 8157 rss: 25Mb L: 4/4 MS: 2 ChangeBinInt-ChangeBinInt-
#16384	pulse  ft: 9 corp: 3/9b lim: 17 exec/s: 8192 rss: 25Mb
#18469	NEW    ft: 11 corp: 4/13b lim: 17 exec/s: 9234 rss: 25Mb L: 4/4 MS: 5 EraseBytes-InsertByte-ChangeBit-ChangeBit-CopyPart-
#19680	NEW    ft: 13 corp: 5/17b lim: 17 exec/s: 9840 rss: 25Mb L: 4/4 MS: 1 CopyPart-
panic: ([]uint8) 0xc00000e018

goroutine 17 [running, locked to thread]:
github.com/stevenjohnstone/fuzztests.FuzzLoop(...)
	github.com/stevenjohnstone/fuzztests/fuzz.go:31
main.LLVMFuzzerTestOneInput(...)
	./main.238601699.go:21
==683== ERROR: libFuzzer: deadly signal
    #0 0x450ddf in __sanitizer_print_stack_trace (/fuzztests/fuzz.libfuzzer+0x450ddf)
    #1 0x430f4b in fuzzer::PrintStackTrace() (/fuzztests/fuzz.libfuzzer+0x430f4b)
    #2 0x414b7b in fuzzer::Fuzzer::CrashCallback() (/fuzztests/fuzz.libfuzzer+0x414b7b)
    #3 0x414b3f in fuzzer::Fuzzer::StaticCrashSignalCallback() (/fuzztests/fuzz.libfuzzer+0x414b3f)
    #4 0x7fe3f929772f  (/lib/x86_64-linux-gnu/libpthread.so.0+0x1272f)
    #5 0x4a53e0 in runtime.raise.abi0 runtime/sys_linux_amd64.s:164

NOTE: libFuzzer has rudimentary signal handlers.
      Combine libFuzzer with AddressSanitizer or similar for better crash reports.
SUMMARY: libFuzzer: deadly signal
MS: 1 ChangeBinInt-; base unit: 4a4b1eab5bdc554747e986e5c6a2152a56c6af7b
0x1,0x3,0x3,0x7,
\x01\x03\x03\x07
artifact_prefix='./'; Test unit written to ./crash-f45be6129befa590730da3f100eebb7217d6b1a0
Base64: AQMDBw==
stat::number_of_executed_units: 31406
stat::average_exec_per_sec:     7851
stat::new_units_added:          4
stat::slowest_unit_time_sec:    0
stat::peak_rss_mb:              27

time elapsed 11.139092246s
```



## Results

When trying to find inputs to ```magic``` which result in a return value ```true``` and cause a crash,  libfuzzer consistently finds a crasher in a couple of seconds with ~100000 executions. The number of executions required for the beta fuzzer
to find the crasher appears to be >= 100x that of libfuzzer.


## TODO

* more comparison tests
* perform more runs to get an idea of the average executions required to complete tests
* run libfuzzer tests with [integer comparison feedback](https://llvm.org/docs/LibFuzzer.html#id32): maybe useful to add to the beta fuzzer?
