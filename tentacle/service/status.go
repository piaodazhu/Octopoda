package service

import (
	"fmt"
	"net"
	"sync"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/message"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type Status struct {
	// Id        int
	Name      string
	Platform  string
	CpuCores  int
	Ip        string
	LocalTime int64

	CpuLoadShort float64
	CpuLoadLong  float64
	MemUsed      uint64
	MemTotal     uint64
	DiskUsed     uint64
	DiskTotal    uint64
}

var nodeStatus Status
var stateLock sync.RWMutex

func initNodeStatus() {
	nodeStatus = Status{
		// Id:        config.GlobalConfig.Id,
		Name:      config.GlobalConfig.Name,
		Platform:  GetCpuInfo(),
		CpuCores:  GetCpuCores(),
		// Ip:        config.GlobalConfig.Ip,
		LocalTime: time.Now().UnixNano(),

		CpuLoadShort: 0.0,
		CpuLoadLong:  0.0,
		MemUsed:      GetMemUsed(),
		MemTotal:     GetMemTotal(),
		DiskUsed:     GetDiskUsed(),
		DiskTotal:    GetDiskTotal(),
	}
	go maintainStatus()
}

func NodeStatus(conn net.Conn, raw []byte) {
	stateLock.RLock()
	state := nodeStatus
	state.LocalTime = time.Now().UnixNano()
	stateLock.RUnlock()
	serialized_info, _ := config.Jsoner.Marshal(&state)
	err := message.SendMessage(conn, message.TypeNodeStatusResponse, serialized_info)
	if err != nil {
		logger.Comm.Println("NodeStatus service error")
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
	if len(percent) == 0 {
		return 0.0
	}
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

func maintainStatus() {
	logger.SysInfo.Println("maintainStatus start")
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

		nodeStatus.CpuLoadShort = load
		nodeStatus.CpuLoadLong = avr_load
		nodeStatus.MemTotal = mem_total
		nodeStatus.MemUsed = mem_used
		nodeStatus.DiskTotal = disk_total
		nodeStatus.DiskUsed = disk_used

		stateLock.Unlock()

		// logger.Server.Println(nodeStatus)

		time.Sleep(time.Second)
	}
}
