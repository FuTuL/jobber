package main

import (
	"github.com/FuTuL/jobber/common"
	"github.com/FuTuL/jobber/ipc"
)

func (self *JobManager) doCatCmd(cmd ipc.CatCmd) ipc.ICmdResp {
	// find job
	job, ok := self.jfile.Jobs[cmd.Job]
	if !ok {
		return ipc.NewErrorCmdResp(&common.Error{What: "No such job."})
	}

	// make response
	return ipc.CatCmdResp{Result: job.Cmd}
}
