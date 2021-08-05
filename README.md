# Comparing Go Fuzz Beta with Libfuzzer

The Go [native fuzzing beta](https://blog.golang.org/fuzz-beta) uses instrumentation
which is identical to the "libfuzzer" build mode. This allows direct comparison of
the mutation engines of the beta fuzzer and libfuzzer.

## Technical Details

When using libfuzzer, integer comparison feedback is [wired up](https://golang.org/src/runtime/libfuzzer.go)
which gives it a slight edge over the beta fuzzer. In the simple test case here, using
the functionality is disabled (with the -use_cmp=0 flag).

https://github.com/stevenjohnstone/go114-fuzz-build is used to build an archive with libfuzzer instrumentation and
a ```LLVMFuzzerTestOneInput``` harness. This is a branch of https://github.com/mdempsky/go114-fuzz-build with a
command line flag added to specify the go compiler. Here, the dev.fuzz branch of golang is used with gotip.

## Tests

The first test is to find an input which makes ```magic``` return true.
```golang
func magic(input []byte) bool {
	return len(input) == 4 && input[0] == 1 && input[1] == 3 && input[2] == 3 && input[3] == 7
}
```

## Usage

Build the container:

```
$ docker build -t fuzztests .
```

To run libfuzzer
```
$ docker run --rm fuzztests mage libfuzzer
INFO: Seed: 3349006478
INFO: 64 Extra Counters
INFO: -max_len is not provided; libFuzzer will not generate inputs larger than 4096 bytes
INFO: A corpus is not provided, starting from an empty corpus
#2	INITED ft: 2 corp: 1/1b lim: 4 exec/s: 0 rss: 25Mb
#16	NEW    ft: 3 corp: 2/5b lim: 4 exec/s: 0 rss: 25Mb L: 4/4 MS: 4 InsertByte-ChangeByte-ChangeBit-CopyPart-
#3207	NEW    ft: 4 corp: 3/9b lim: 6 exec/s: 0 rss: 25Mb L: 4/4 MS: 1 ChangeBinInt-
#32768	pulse  ft: 4 corp: 3/9b lim: 33 exec/s: 16384 rss: 25Mb
#35850	NEW    ft: 5 corp: 4/13b lim: 38 exec/s: 11950 rss: 25Mb L: 4/4 MS: 3 CopyPart-ShuffleBytes-ChangeBit-
#36451	NEW    ft: 6 corp: 5/17b lim: 38 exec/s: 12150 rss: 25Mb L: 4/4 MS: 1 CopyPart-
#65536	pulse  ft: 6 corp: 5/17b lim: 63 exec/s: 13107 rss: 25Mb
panic: found

goroutine 17 [running, locked to thread]:
github.com/stevenjohnstone/fuzztests.Fuzz(...)
	github.com/stevenjohnstone/fuzztests/fuzz.go:9
main.LLVMFuzzerTestOneInput(...)
	./main.143848471.go:21
==671== ERROR: libFuzzer: deadly signal
    #0 0x450ddf in __sanitizer_print_stack_trace (/fuzztests/fuzz.libfuzzer+0x450ddf)
    #1 0x430f4b in fuzzer::PrintStackTrace() (/fuzztests/fuzz.libfuzzer+0x430f4b)
    #2 0x414b7b in fuzzer::Fuzzer::CrashCallback() (/fuzztests/fuzz.libfuzzer+0x414b7b)
    #3 0x414b3f in fuzzer::Fuzzer::StaticCrashSignalCallback() (/fuzztests/fuzz.libfuzzer+0x414b3f)
    #4 0x7f2b4ecab72f  (/lib/x86_64-linux-gnu/libpthread.so.0+0x1272f)
    #5 0x4a52a0 in runtime.raise.abi0 runtime/sys_linux_amd64.s:164

NOTE: libFuzzer has rudimentary signal handlers.
      Combine libFuzzer with AddressSanitizer or similar for better crash reports.
SUMMARY: libFuzzer: deadly signal
MS: 3 ShuffleBytes-ShuffleBytes-ChangeByte-; base unit: 307fe878e1eebdafe7c56fc8f482407037a34736
0x1,0x3,0x3,0x7,
\x01\x03\x03\x07
artifact_prefix='./'; Test unit written to ./crash-f45be6129befa590730da3f100eebb7217d6b1a0
Base64: AQMDBw==
stat::number_of_executed_units: 89044
stat::average_exec_per_sec:     12720
stat::new_units_added:          4
stat::slowest_unit_time_sec:    0
stat::peak_rss_mb:              27

time elapsed 13.248694284s
```

To run the beta fuzzer
```
$ docker run --rm fuzztests mage betafuzzer
docker run --rm fuzztests mage betafuzzer
fuzzing, elapsed: 3.0s, execs: 348353 (116053/sec), workers: 8, interesting: 5
fuzzing, elapsed: 6.0s, execs: 711390 (118534/sec), workers: 8, interesting: 5
fuzzing, elapsed: 9.0s, execs: 1081601 (120164/sec), workers: 8, interesting: 5
fuzzing, elapsed: 12.0s, execs: 1432223 (119335/sec), workers: 8, interesting: 5
fuzzing, elapsed: 15.0s, execs: 1765989 (117721/sec), workers: 8, interesting: 5
fuzzing, elapsed: 18.0s, execs: 2130881 (118377/sec), workers: 8, interesting: 5
fuzzing, elapsed: 21.0s, execs: 2512201 (119622/sec), workers: 8, interesting: 5
fuzzing, elapsed: 24.0s, execs: 2878962 (119947/sec), workers: 8, interesting: 5
fuzzing, elapsed: 27.0s, execs: 3248521 (120311/sec), workers: 8, interesting: 5
fuzzing, elapsed: 30.0s, execs: 3612030 (120398/sec), workers: 8, interesting: 5
fuzzing, elapsed: 33.0s, execs: 3972515 (120374/sec), workers: 8, interesting: 5
fuzzing, elapsed: 36.0s, execs: 4333992 (120384/sec), workers: 8, interesting: 5
fuzzing, elapsed: 39.0s, execs: 4697998 (120457/sec), workers: 8, interesting: 5
fuzzing, elapsed: 42.0s, execs: 5066002 (120616/sec), workers: 8, interesting: 5
fuzzing, elapsed: 45.0s, execs: 5430306 (120669/sec), workers: 8, interesting: 5
fuzzing, elapsed: 48.0s, execs: 5792458 (120666/sec), workers: 8, interesting: 5
fuzzing, elapsed: 51.0s, execs: 6155857 (120701/sec), workers: 8, interesting: 5
fuzzing, elapsed: 54.0s, execs: 6522308 (120781/sec), workers: 8, interesting: 5
fuzzing, elapsed: 57.0s, execs: 6889993 (120875/sec), workers: 8, interesting: 5
fuzzing, elapsed: 60.0s, execs: 7255693 (120926/sec), workers: 8, interesting: 5
fuzzing, elapsed: 63.0s, execs: 7619319 (120919/sec), workers: 8, interesting: 6
fuzzing, elapsed: 66.0s, execs: 7980752 (120918/sec), workers: 8, interesting: 7
fuzzing, elapsed: 69.0s, execs: 8354791 (121077/sec), workers: 8, interesting: 7
fuzzing, elapsed: 72.0s, execs: 8715690 (121049/sec), workers: 8, interesting: 7
fuzzing, elapsed: 75.0s, execs: 9076681 (121020/sec), workers: 8, interesting: 7
fuzzing, elapsed: 78.0s, execs: 9448445 (121132/sec), workers: 8, interesting: 7
fuzzing, elapsed: 81.0s, execs: 9816473 (121189/sec), workers: 8, interesting: 7
fuzzing, elapsed: 84.0s, execs: 10178700 (121174/sec), workers: 8, interesting: 7
fuzzing, elapsed: 87.0s, execs: 10538972 (121127/sec), workers: 8, interesting: 7
fuzzing, elapsed: 90.0s, execs: 10898534 (121093/sec), workers: 8, interesting: 7
found a crash, minimizing...
fuzzing, elapsed: 90.7s, execs: 10972162 (120956/sec), workers: 8, interesting: 7
--- FAIL: FuzzBeta (90.71s)
        --- FAIL: FuzzBeta (0.00s)
            fuzz_test.go:10: magic is [1 3 3 7]
    
    Crash written to testdata/corpus/FuzzBeta/1757d23fd60223bd5a11cfd3a7978f28cdb2b98e0b81542690f8f75ba96d043d
    To re-run:
    go test github.com/stevenjohnstone/fuzztests -run=FuzzBeta/1757d23fd60223bd5a11cfd3a7978f28cdb2b98e0b81542690f8f75ba96d043d
FAIL
exit status 1
FAIL	github.com/stevenjohnstone/fuzztests	90.767s

time elapsed 1m36.380102004s
Error: running "gotip test -fuzz=FuzzBeta" failed with exit code 1
```

## Results

Libfuzzer consistently finds a crasher in a couple of seconds with ~100000 executions. The number of executions required for the beta fuzzer
to find the crasher appears to be >= 100x that of libfuzzer.


## TODO

* more comparison tests
* perform more runs to get an idea of the average executions required to complete tests
* run libfuzzer tests with [integer comparison feedback](https://llvm.org/docs/LibFuzzer.html#id32): maybe useful to add to the beta fuzzer?
