package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/spf13/hugo/commands"
	"github.com/spf13/hugo/helpers"
	jww "github.com/spf13/jwalterweatherman"
)

var (
	renderToMem       = true
	firstOnly         = false
	iterationsPerSite = 4
	cpuProfile        = flag.String("cpuProfile", "", "write cpu profile to file")
	heapProfile       = flag.String("heapProfile", "", "write heap profile to file")
)

type benchmark struct {
	sites []*site
}

type site struct {
	name    string
	path    string
	elapsed time.Duration
	runs    int
	errors  []error
}

func main() {
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	os.Args = os.Args[:1]

	fmt.Println("Start Hugo benchmark ...")
	b := createBench(firstOnly)
	for _, s := range b.sites {
		for i := 0; i < iterationsPerSite; i++ {
			s.build()
		}
	}

	var (
		total time.Duration
		ok    = true
	)

	fmt.Println("\n\n")

	for _, s := range b.sites {
		total += s.elapsed
		status := "OK"
		if len(s.errors) > 0 {
			ok = false
			status = fmt.Sprintf("%q", s.errors)
		}
		fmt.Printf("%s|%d|%s|%s\n", s.elapsed, s.runs, s.name, status)
	}

	if *heapProfile != "" {
		f, err := os.Create(*heapProfile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.WriteHeapProfile(f)
		f.Close()
	}

	fmt.Println("\n\n")
	totalStatus := "All OK"
	if !ok {
		totalStatus = "Failed!"
	}

	fmt.Println(runtime.Version(), helpers.HugoVersion(), total, "-", totalStatus)

	fmt.Println("\n\n")
}

func createBench(firstOnly bool) *benchmark {
	pwd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	sitesPath := filepath.Join(pwd, "sites")

	fis, err := ioutil.ReadDir(sitesPath)

	if err != nil {
		log.Fatal(err)
	}

	b := &benchmark{}

	for i, fi := range fis {
		if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") {
			subFolder := ""
			if strings.HasPrefix(fi.Name(), "hugo") {
				subFolder = "docs"
			}
			p := filepath.Join(sitesPath, fi.Name(), subFolder)
			b.sites = append(b.sites, &site{name: fi.Name(), path: p})
		}

		if i == 0 && firstOnly {
			break
		}

	}

	return b
}

func (b *benchmark) build() error {
	var elapsed time.Duration

	for _, s := range b.sites {
		s.build()
		if len(s.errors) > 0 {
			return fmt.Errorf("Build failed: %q", s.errors)
		}
		elapsed += s.elapsed
	}
	return nil
}

func (s *site) build() {

	defer s.incrementElapsed(time.Now())
	s.runs += 1
	err := buildHugoSite(s.path)
	if err != nil {
		s.errors = append(s.errors, err)
	}
}

func (s *site) incrementElapsed(start time.Time) {
	s.elapsed += time.Since(start)
}

var logError = errors.New("error(s) in log")

func buildHugoSite(path string) error {
	defer jww.ResetLogCounters()
	defer commands.Reset()
	flags := []string{"--quiet", fmt.Sprintf("--source=%s", path)}
	os.Args = []string{os.Args[0]}

	if renderToMem {
		flags = append(flags, "--renderToMemory")
	}

	if err := commands.HugoCmd.ParseFlags(flags); err != nil {
		log.Fatal(err)
	}

	if _, err := commands.HugoCmd.ExecuteC(); err != nil {
		return err

	}

	if jww.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError) > 0 {
		return logError
	}

	return nil
}

func init() {
	jww.SetStdoutThreshold(jww.LevelError)
	commands.HugoCmd.SilenceUsage = true
	commands.AddCommands()
}
