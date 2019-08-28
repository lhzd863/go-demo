package ex

import (
        "log"
        "os/exec"
        "bytes"
     
        "gitlab.51y5.net/liuhui/rpt-batch/module"

)

type Exec struct {
        metaConf *module.MetaConf
}

func NewExec(mc *module.MetaConf) *Exec{
        return &Exec{
                metaConf: mc,
        }
}

func (e *Exec) Execmd(cmdstr string) (string, error) {
        cmd := exec.Command("/bin/bash", "-c", cmdstr)
        var out bytes.Buffer
        cmd.Stdout = &out
        err := cmd.Run()
        if err != nil {
                log.Printf("error: %s", err)
                return "", err
        }
        return out.String(), err
}
