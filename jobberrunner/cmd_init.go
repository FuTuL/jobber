package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/FuTuL/jobber/common"
	"github.com/FuTuL/jobber/ipc"
)

const gDefaultJobfile = `## This is your jobfile: use it to tell Jobber what jobs you want it to
## run on your behalf.  For details of what you can specify here,
## please see https://dshearer.github.io/jobber/doc/.
##
## It consists of two sections: "prefs" and "jobs".  In "prefs" you can
## set various general settings.  In "jobs", you define your jobs.

version: 1.4

prefs:
  ## You can have the Jobber daemon keep a log of various activities
  ## with the "logPath" setting; the log will be written to the given
  ## path (if the path is relative, it will be interpreted relative to
  ## your home directory).  Your user account must be able to write to
  ## the given path.  NOTE: This is NOT where logs about job runs
  ## are stored --- for that, see the "runLog" setting below.  WARNING:
  ## Jobber will NOT rotate this file.
  #logPath: jobber-log

  ## You can specify how info about past runs is stored.  For
  ## "type: memory" (the default), they are stored in memory and
  ## are lost when the Jobber service stops.
  #runLog:
  #    type: memory
  #    maxLen: 100  # the max number of entries to remember

  ## For "type: file", past run logs are stored on disk.  The log file is
  ## rotated when it reaches a size of 'maxFileLen' MB.  Up to
  ## 'maxHistories' historical run logs (that is, not including the
  ## current one) are kept.
  #runLog:
  #    type: file
  #    path: /tmp/claudius
  #    maxFileLen: 50m  # in MB
  #    maxHistories: 5

resultSinks:
  #- &programSink
  #  type: program
  #  path: /home/handleError.sh

  #- &systemEmailSink
  #  type: system-email

  #- &filesystemSink
  #  type: filesystem
  #  path: /path/to/dir
  #  data: [stdout, stderr]
  #  maxAgeDays: 10

jobs:
  ## This section must contain a YAML sequence of maps like the following:
  #DailyBackup:
  #    cmd: backup daily  # shell command to execute
  #    time: '* * * * * *'  # SEC MIN HOUR MONTH_DAY MONTH WEEK_DAY.
  #    onError: Continue  # what to do when the job has an error: Stop, Backoff, or Continue
  #    notifyOnError: [*programSink]  # what to do with result when job has an error
  #    notifyOnFailure: [*systemEmailSink, *programSink]  # what to do with result when the job stops due to errors
  #    notifyOnSuccess: [*filesystemSink]  # what to do with result when the job succeeds
`

func (self *JobManager) doInitCmd(cmd ipc.InitCmd) ipc.ICmdResp {
	var resp ipc.InitCmdResp
	resp.JobfilePath = self.jobfilePath

	// open file for writing
	f, err := os.OpenFile(self.jobfilePath,
		os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		if strings.Contains(err.Error(), "file exists") {
			msg := fmt.Sprintf("Jobfile already exists at %v\n",
				self.jobfilePath)
			err = &common.Error{What: msg}
		}
		return ipc.NewErrorCmdResp(err)
	}
	defer f.Close()

	// write default jobfile
	_, err = f.WriteString(gDefaultJobfile)
	if err != nil {
		return ipc.NewErrorCmdResp(err)
	}

	return resp
}
