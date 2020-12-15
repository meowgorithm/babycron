package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"

	te "github.com/muesli/termenv"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var (
	Version   = "unknown; compiled from source"
	CommitSHA = ""

	runOnStart bool

	rootCmd = &cobra.Command{
		Use:     "babycron [cron expression] [task]",
		Short:   "A cron scheduler for a single task",
		Example: "  babycron '* */6 * * *' 'sh /backup.sh'\n  babycron -r '30 16 * * *' '/bin/bash /backup'",
		Args:    cobra.ExactArgs(2),
		RunE:    execute,
	}

	subtle = te.Style{}.Foreground(te.ColorProfile().Color("241")).Styled
	warn   = te.Style{}.Foreground(te.ColorProfile().Color("203")).Styled
)

func execute(cmd *cobra.Command, args []string) error {
	c := cron.New()

	if runOnStart {
		go runJob(args[1])
	}

	_, err := c.AddFunc(args[0], func() {
		go runJob(args[1])
	})
	if err != nil {
		return fmt.Errorf("could not parse cron expression: %v", err)
	}

	report("Babycron running...")
	c.Start()

	select {}
}

func runJob(argString string) {
	args := strings.Split(argString, " ")

	// Search path for executable if no path is specified
	exe, err := exec.LookPath(args[0])
	if err != nil {
		report("Error getting path of %s: %v\n", args[0], err)
		return
	}
	args[0] = exe

	cmd := &exec.Cmd{
		Path: exe,
		Args: args,
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		report("Could not get stdout pipe: %v", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		report("Could not get stderr pipe: %v", err)
		return
	}

	report("Starting job...")
	startTime := time.Now()
	if err := cmd.Start(); err != nil {
		report("Coult not start job: %v", err)
		return
	}

	go listen(bufio.NewReader(stdout), false)
	go listen(bufio.NewReader(stderr), true)
	err = cmd.Wait()

	dur := time.Since(startTime).String()
	if err != nil {
		report("Job exited with error: %v", err)
		report("Job failed %s", subtle(dur))
	} else {
		report("Job finished %s", subtle(dur))
	}
}

func listen(reader io.Reader, isErrBuf bool) {
	buf := bufio.NewReader(reader)

	prefix := subtle(">")
	if isErrBuf {
		prefix = warn("[ERR]")
	}

	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			return
		} else if err != nil {
			report("Error reading line: %v", err)
			return
		}
		report("%s %s", prefix, string(line))
	}
}

func report(fmt string, args ...interface{}) {
	log.Printf(fmt+"\n", args...)
}

func init() {
	if len(CommitSHA) >= 7 {
		vt := rootCmd.VersionTemplate()
		rootCmd.SetVersionTemplate(vt[:len(vt)-1] + " (" + CommitSHA[0:7] + ")\n")
	}

	rootCmd.Version = Version
	rootCmd.PersistentFlags().BoolVarP(&runOnStart, "run-on-start", "r", false, "also run job when starting")
}

func main() {
	rootCmd.Execute()
}
