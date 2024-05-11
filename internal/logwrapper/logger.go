package logwrapper

import "fmt"

type ILogMsg interface {
	String() string
}

type ILoggerWrapper interface {
	Log(...ILogMsg)
	AddLogger(l ILoggerWrapper)
}

type DebugMsg struct {
	ILogMsg
	msg string
}

func (m *DebugMsg) String() string {
	return fmt.Sprintf("DEBUG | %v", m.msg)
}

func NewDebugMsg(msg string) *DebugMsg {
	return &DebugMsg{msg: msg}
}

type ErrorMsg struct {
	ILogMsg
	msg string
}

func (m *ErrorMsg) String() string {
	return fmt.Sprintf("ERROR | %v", m.msg)
}

func NewErrorMsg(msg string) *ErrorMsg {
	return &ErrorMsg{msg: msg}
}

type InforMsg struct {
	ILogMsg
	msg string
}

func (m *InforMsg) String() string {
	return fmt.Sprintf("INFOR | %v", m.msg)
}

func NewInforMsg(msg string) *InforMsg {
	return &InforMsg{msg: msg}
}

type WarnMsg struct {
	ILogMsg
	msg string
}

func (m *WarnMsg) String() string {
	return m.msg
}

func NewWarnMsg(msg string) *WarnMsg {
	return &WarnMsg{msg: msg}
}
