package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/hugo/commands"
	"github.com/spf13/hugo/helpers"
	jww "github.com/spf13/jwalterweatherman"
)

var renderToMem = true

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
	fmt.Println("Start Hugo benchmark ...")
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

	for _, fi := range fis {
		if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") {
			b.sites = append(b.sites, &site{name: fi.Name(), path: filepath.Join(sitesPath, fi.Name())})
		}
	}

	for _, s := range b.sites {
		for i := 0; i < 6; i++ {
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

	fmt.Println("\n\n")
	totalStatus := "All OK"
	if !ok {
		totalStatus = "Failed!"
	}

	fmt.Println(runtime.Version(), helpers.HugoVersion(), total, "-", totalStatus)

	fmt.Println("\n\n")
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
	flags := []string{fmt.Sprintf("--source=%s", path)}

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
