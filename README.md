# hugo-benchmark

Benchark for building a set of random, but real [Hugo](https://github.com/spf13/hugo) sites. 

## Install and run

Clone an run the main:

```
git clone --recursive https://github.com/bep/hugo-benchmark.git
cd hugo-benchmark
go get -u -v
go test -bench=".*" -test.benchmem=true -count 3 > 1.bench
// checkout another Hugo branch
go test -bench=".*" -test.benchmem=true -count 3 > 2.bench
benchcmp -best 1.bench 2.bench
```


