package service

import (
	"encoding/json"
	"fmt"
	"net"
	"nworkerd/config"
	"nworkerd/logger"
	"nworkerd/message"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type State struct {
	Id        int
	Name      string
	Platform  string
	CpuCores  int
	Ip        string
	StartTime int64

	CpuLoadShort float64
	CpuLoadLong  float64
	MemUsed      uint64
	MemTotal     uint64
	DiskUsed     uint64
	DiskTotal    uint64
}

var nodeState State
var stateLock sync.RWMutex

func initNodeState() {
	nodeState = State{
		Id:        config.GlobalConfig.Worker.Id,
		Name:      config.GlobalConfig.Worker.Name,
		Platform:  GetCpuInfo(),
		CpuCores:  GetCpuCores(),
		Ip:        config.GlobalConfig.Worker.Ip,
		StartTime: time.Now().Unix(),

		CpuLoadShort: 0.0,
		CpuLoadLong:  0.0,
		MemUsed:      GetMemUsed(),
		MemTotal:     GetMemTotal(),
		DiskUsed:     GetDiskUsed(),
		DiskTotal:    GetDiskTotal(),
	}
	go maintainInfo()
}

func NodeState(conn net.Conn, raw []byte) {
	stateLock.RLock()
	state := nodeState
	stateLock.RUnlock()
	serialized_info, _ := json.Marshal(&state)
	err := message.SendMessage(conn, message.TypeNodeStateResponse, serialized_info)
	if err != nil {
		logger.Server.Println("NodeState service error")
	}
}

func GetCpuInfo() string {

	info, _ := cpu.Info()
	return fmt.Sprint(info[0].ModelName)
}

func GetCpuCores() int {
	Cnt, _ := cpu.Counts(true)
	return Cnt
}

// maybe not 1 second...
func GetCpuLoad() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetMemTotal() uint64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.Total
}

func GetMemUsed() uint64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.Used
}

func GetDiskTotal() uint64 {
	// should be workspace
	diskInfo, _ := disk.Usage("/")
	return diskInfo.Total
}

func GetDiskUsed() uint64 {
	// should be workspace
	diskInfo, _ := disk.Usage("/")
	return diskInfo.Used
}

func maintainInfo() {
	logger.Server.Println("maintainInfo start")
	CpuLoadLen := 10
	CpuLoadBuf := make([]float64, CpuLoadLen)
	CpuLoadPtr := 0
	CpuLoadTot := 0.0
	for {
		load := GetCpuLoad()
		mem_total := GetMemTotal()
		mem_used := GetMemUsed()
		disk_total := GetDiskTotal()
		disk_used := GetDiskUsed()

		// use a circle queue
		CpuLoadTot -= CpuLoadBuf[CpuLoadPtr]
		CpuLoadTot += load
		CpuLoadBuf[CpuLoadPtr] = load
		CpuLoadPtr++
		if CpuLoadPtr == CpuLoadLen {
			CpuLoadPtr = 0
		}
		avr_load := CpuLoadTot / float64(CpuLoadLen)

		stateLock.Lock()

		nodeState.CpuLoadShort = load
		nodeState.CpuLoadLong = avr_load
		nodeState.MemTotal = mem_total
		nodeState.MemUsed = mem_used
		nodeState.DiskTotal = disk_total
		nodeState.DiskUsed = disk_used

		stateLock.Unlock()

		// logger.Server.Println(nodeState)

		time.Sleep(time.Second)
	}
}
