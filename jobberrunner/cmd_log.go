package main

import (
	"github.com/FuTuL/jobber/common"
	"github.com/FuTuL/jobber/ipc"
)

func (self *JobManager) doLogCmd(cmd ipc.LogCmd) ipc.ICmdResp {
	// make log list
	var logDescs []ipc.LogDesc
	entries, err := self.jfile.Prefs.RunLog.GetAll()
	if err != nil {
		return ipc.NewErrorCmdResp(err)
	}
	for _, l := range entries {
		logDesc := ipc.LogDesc{
			Time:      l.Time,
			Job:       l.JobName,
			Succeeded: l.Fate == common.SubprocFateSucceeded, // deprecated
			Fate:      l.Fate.String(),
			ExecTime:  l.ExecTime,
			Result:    l.Result.String(),
		}
		logDescs = append(logDescs, logDesc)
	}

	// make response
	return ipc.LogCmdResp{Logs: logDescs}
}
