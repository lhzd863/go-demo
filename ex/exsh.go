func (ex *ShExec) Executesh1(logname string, execcmd string) (string,error) {
        ex.slog.Info(aty.INFO,execcmd)
        cmd := exec.Command("/bin/bash", "-c", execcmd)
        stdout, err := cmd.StdoutPipe()
        if err != nil {
                msg := fmt.Sprintf("failed to cmd.StdoutPipe: %v", err)
                ex.slog.Info(aty.ERR, msg)
                fmt.Println(msg)
                return "1", fmt.Errorf("%v %v",msg,err)
        }
        stderr, err := cmd.StderrPipe()
        if err != nil {
                msg := fmt.Sprintf("failed to cmd.StdoutPipe: %v", err)
                ex.slog.Info(aty.ERR, msg)
                fmt.Println(msg)
                return "1",fmt.Errorf("%v %v",msg,err)
        }
        cmd.Start()
        reader := bufio.NewReader(stdout)
        go func() {
                for {
                        line, err2 := reader.ReadString('\n')
                        if err2 != nil || io.EOF == err2 {
                                break
                        }
                        ex.log(logname, line)
                }
        }()
        readererr := bufio.NewReader(stderr)
        go func() {
                for {
                        line, err2 := readererr.ReadString('\n')
                        if err2 != nil || io.EOF == err2 {
                                break
                        }
                        ex.log(logname, line)
                }
        }()

        cmd.Wait()
        retcd:=string(fmt.Sprintf("%v", cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()))
        retcd=strings.Replace(retcd, " ", "", -1)
        retcd=strings.Replace(retcd, "\n", "", -1)
        return retcd,nil
}
