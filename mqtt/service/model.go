package service

import (
	"fmt"
	"strconv"
	"time"
)

type DataSlice interface {
	toDataSlice() []Data
	toSlice() interface{}
}

type ProcessSlice []*Process

func (ps ProcessSlice) toDataSlice() []Data {
	dataslice := make([]Data, 0)
	for _, v := range ps {
		dataslice = append(dataslice, v)
	}

	return dataslice
}

func (ps ProcessSlice) toSlice() interface{} {
	return ps
}

type Data interface {
	ToMap() map[string]string
	ID() string
	Update(d Data) Data
	CmpUpdateIsSame(d Data) bool
}

var _ Data = (*Process)(nil)

type Process struct {
	Pid      int     `json:"pid"`
	Command  string  `json:"command"`
	Cmd      string  `json:"cmd"`
	State    string  `json:"state"`
	Started  int64   `json:"started"`
	Uid      int64   `json:"uid"`
	User     string  `json:"user"`
	Gid      int     `json:"gid"`
	Group    string  `json:"group"`
	Ppid     int     `json:"ppid"`
	CpuUsage float32 `json:"cpu_usage"`
	MemUsage float32 `json:"mem_usage"`
}

func (p *Process) ToMap() map[string]string {
	return map[string]string{
		"Pid":      strconv.Itoa(p.Pid),
		"Command":  p.Command,
		"Cmd":      p.Cmd,
		"State":    p.State,
		"Started":  time.Unix(p.Started, 0).Format("2006-01-02 15:04:05"),
		"Uid":      strconv.Itoa(int(p.Uid)),
		"User":     p.User,
		"Gid":      strconv.Itoa(p.Gid),
		"Group":    p.Group,
		"Ppid":     strconv.Itoa(p.Ppid),
		"CpuUsage": fmt.Sprintf("%f", p.CpuUsage),
		"MemUsage": fmt.Sprintf("%f", p.MemUsage),
	}
}

func (p *Process) ID() string {
	return fmt.Sprintf("%d", p.Pid)
}

func (p Process) Update(d Data) Data {
	pro, ok := d.(*Process)
	if !ok {
		return &p
	}
	if !IsZero(pro.Cmd) {
		p.Cmd = pro.Cmd
	}
	if !IsZero(pro.Command) {
		p.Command = pro.Command
	}
	if !IsZero(pro.CpuUsage) {
		p.CpuUsage = pro.CpuUsage
	}
	if !IsZero(pro.MemUsage) {
		p.MemUsage = pro.MemUsage
	}
	if !IsZero(pro.State) {
		p.State = pro.State
	}

	return &p
}

func (p Process) CmpUpdateIsSame(d Data) bool {
	pro, ok := d.(*Process)
	if !ok {
		return false
	}
	if pro.Pid == p.Pid {
		if !IsZero(pro.Cmd) && pro.Cmd != p.Cmd {
			return false
		}
		if !IsZero(pro.Command) && pro.Command != p.Command {
			return false
		}
		if !IsZero(pro.CpuUsage) && pro.CpuUsage != p.CpuUsage {
			return false
		}
		if !IsZero(pro.MemUsage) && pro.MemUsage != p.MemUsage {
			return false
		}
	} else {
		return false
	}

	return true
}

type Payload[T DataSlice] struct {
	Mid      string `json:"mid"`
	Reply    string `json:"reply"`
	DeviceId string `json:"device_id"`
	Data     struct {
		Add    T `json:"add"`
		Del    T `json:"del"`
		Update T `json:"update"`
	} `json:"data"`
	Timestamp int64 `json:"timestamp"`
}

func (p Payload[T]) IsSameMsg(msgType string) bool {
	dataCenter.lock.RLock()
	defer dataCenter.lock.RUnlock()

	switch msgType {
	case ProcessKey:
		if len(dataCenter.records[msgType].Add) != len(p.Data.Add.toSlice().(ProcessSlice)) ||
			len(dataCenter.records[msgType].Update) != len(p.Data.Update.toSlice().(ProcessSlice)) ||
			len(dataCenter.records[msgType].Del) != len(p.Data.Del.toSlice().(ProcessSlice)) {
			return false
		}
	outAdd:
		for _, d := range dataCenter.records[msgType].Add {
			for _, pData := range p.Data.Add.toSlice().(ProcessSlice) {
				if d.ID() == pData.ID() {
					continue outAdd
				}
			}
			return false
		}
	outDel:
		for _, d := range dataCenter.records[msgType].Del {
			for _, pData := range p.Data.Del.toSlice().(ProcessSlice) {
				if d.ID() == pData.ID() {
					continue outDel
				}
			}
			return false
		}
	outUpdate:
		for _, d := range dataCenter.records[msgType].Update {
			for _, pData := range p.Data.Update.toSlice().(ProcessSlice) {
				if d.ID() == pData.ID() {
					if !d.CmpUpdateIsSame(pData) {
						return false
					}
					continue outUpdate
				}
			}
			return false
		}
	}

	return true
}

type record struct {
	Mid       string `json:"mid"`
	Reply     string `json:"reply"`
	DeviceId  string `json:"device_id"`
	Timestamp int64  `json:"timestamp"`

	Add    []Data
	Del    []Data
	Update []Data
}
