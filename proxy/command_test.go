package proxy

import "testing"

func TestCommand_FireCmdBind(t *testing.T) {
	cmd := &Command{
		rqst: &Rqst{
			cmd: cmdBind,
		},
	}
	rsp, err := cmd.Fire()
	if err != nil {
		t.Errorf("Excpected err nil, but got %v", err)
	}
	if rsp != rspCommandNotSupport {
		t.Errorf("Excpected rsp=%v, but got %v", rspCommandNotSupport, rsp)
	}
}

func TestCommand_FireCmdUDP(t *testing.T) {
	cmd := &Command{
		rqst: &Rqst{
			cmd: cmdUDP,
		},
	}
	rsp, err := cmd.Fire()
	if err != nil {
		t.Errorf("Excpected err nil, but got %v", err)
	}
	if rsp != rspCommandNotSupport {
		t.Errorf("Excpected rsp=%v, but got %v", rspCommandNotSupport, rsp)
	}
}

func TestCommand_FireCmdDefault(t *testing.T) {
	cmd := &Command{
		rqst: &Rqst{
			cmd: 0xff,
		},
	}
	rsp, err := cmd.Fire()
	if err != nil {
		t.Errorf("Excpected err nil, but got %v", err)
	}
	if rsp != rspCommandNotSupport {
		t.Errorf("Excpected rsp=%v, but got %v", rspCommandNotSupport, rsp)
	}
}
