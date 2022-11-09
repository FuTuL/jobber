package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/FuTuL/jobber/ipc"
)

func doResumeCmd(args []string) int {
	// parse flags
	flagSet := flag.NewFlagSet(ResumeCmdStr, flag.ExitOnError)
	flagSet.Usage = subcmdUsage(ResumeCmdStr, "[JOBS...]", flagSet)
	var help_p *bool = flagSet.Bool("h", false, "help")
	var timeout_p = flagSet.Duration("t", 5*time.Second, "timeout")
	flagSet.Parse(args)

	if *help_p {
		flagSet.Usage()
		return 0
	}

	// get jobs
	var jobs []string = flagSet.Args()

	// get current user
	usr, err := user.Current()
	if err != nil {
		fmt.Fprintf(
			os.Stderr, "Failed to get current user: %v\n", err,
		)
		return 1
	}

	// send command
	var resp ipc.ResumeCmdResp
	err = CallDaemon(
		"IpcService.Resume",
		ipc.ResumeCmd{Jobs: jobs},
		&resp,
		usr,
		timeout_p,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	// handle response
	fmt.Printf("Resumed %v jobs.\n", resp.NumResumed)
	return 0
}
