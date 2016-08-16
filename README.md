# hugo-benchmark

Builds a set of random, but real [Hugo](https://github.com/spf13/hugo) sites and prints some stats. 

## Install and run

Clone an run the main:

```
git clone --recursive https://github.com/bep/hugo-benchmark.git
cd hugo-benchmark
go get -u -v
go run main.go
```


## 2016-08-16

4 random but fairly big Hugo sites, each built 6 times, rendered to memory.

The MULTILINGUAL branch is a new and [comprehensive branch](https://github.com/spf13/hugo/pull/2303) adding multiple languages support to Hugo.

- go1.6.2 0.17-DEV 24.479727984s - All OK
- go1.6.2 0.17-MULTILINGUAL 26.49965009s - All OK
- go1.7 0.17-DEV 21.830323612s - All OK 
- go1.7 0.17-MULTILINGUAL 23.164846844s - All OK
