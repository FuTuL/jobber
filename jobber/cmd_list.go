package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/FuTuL/jobber/common"
	"github.com/FuTuL/jobber/ipc"
)

type ListRespRec struct {
	usr  *user.User
	resp *ipc.ListJobsCmdResp
}

func formatTime(t *time.Time) string {
	if t == nil {
		return "none"
	} else {
		tmp := t.Local()
		return tmp.Format("Jan _2 15:04:05 -0700 MST")
	}
}

func formatResponseRecs(recs []ListRespRec, showUser bool) string {
	// make table header
	var buffer bytes.Buffer
	var writer *tabwriter.Writer = tabwriter.NewWriter(&buffer,
		5, 0, 2, ' ', 0)
	headers := []string{
		"NAME",
		"STATUS",
		"SEC/MIN/HR/MDAY/MTH/WDAY",
		"NEXT RUN TIME",
		"NOTIFY ON SUCCESS",
		"NOTIFY ON ERR",
		"NOTIFY ON FAIL",
		"ERR HANDLER",
	}
	if showUser {
		headers = append(headers, "USER")
	}
	writer.Write([]byte(strings.Join(headers, "\t")))
	writer.Write([]byte("\n"))

	// make table rows
	var rows []string
	for _, respRec := range recs {
		for _, j := range respRec.resp.Jobs {
			fields := []string{
				j.Name,
				j.Status,
				j.Schedule,
				formatTime(j.NextRunTime),
				fmt.Sprintf("%v", j.NotifyOnSuccess),
				fmt.Sprintf("%v", j.NotifyOnErr),
				fmt.Sprintf("%v", j.NotifyOnFail),
				j.ErrHandler,
			}
			if showUser {
				fields = append(fields, respRec.usr.Username)
			}
			rows = append(rows, strings.Join(fields, "\t"))
		}
	}
	writer.Write([]byte(strings.Join(rows, "\n")))

	// finish up
	writer.Flush()
	return buffer.String()
}

func doListCmd_allUsers(timeout_p *time.Duration) int {
	// get all users
	users, err := common.AllUsersWithSockets()
	if err != nil {
		fmt.Fprintf(
			os.Stderr, "Failed to get all users: %v\n", err,
		)
		return 1
	}

	// send cmd
	var responses []ListRespRec
	for _, usr := range users {
		var resp ipc.ListJobsCmdResp
		err = CallDaemon(
			"IpcService.ListJobs",
			ipc.ListJobsCmd{},
			&resp,
			usr,
			timeout_p,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"Failed to list jobs for %v: %v\n", usr.Username, err)
			continue
		}
		rec := ListRespRec{usr: usr, resp: &resp}
		responses = append(responses, rec)
	}

	// display response records
	fmt.Println(formatResponseRecs(responses, true))

	return 0
}

func doListCmd_currUser(timeout_p *time.Duration) int {
	// get current user
	usr, err := user.Current()
	if err != nil {
		fmt.Fprintf(
			os.Stderr, "Failed to get current user: %v\n", err,
		)
		return 1
	}

	// send cmd
	var resp ipc.ListJobsCmdResp
	err = CallDaemon(
		"IpcService.ListJobs",
		ipc.ListJobsCmd{},
		&resp,
		usr,
		timeout_p,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	// display response records
	rec := ListRespRec{usr: usr, resp: &resp}
	fmt.Println(formatResponseRecs([]ListRespRec{rec}, false))

	return 0
}

func doListCmd(args []string) int {
	// parse flags
	flagSet := flag.NewFlagSet(ListCmdStr, flag.ExitOnError)
	flagSet.Usage = subcmdUsage(ListCmdStr, "", flagSet)
	var help_p = flagSet.Bool("h", false, "help")
	var allUsers_p = flagSet.Bool("a", false, "all-users")
	var timeout_p = flagSet.Duration("t", 5*time.Second, "timeout")
	flagSet.Parse(args)

	if *help_p {
		flagSet.Usage()
		return 0
	}

	if *allUsers_p {
		return doListCmd_allUsers(timeout_p)
	} else {
		return doListCmd_currUser(timeout_p)
	}
}
