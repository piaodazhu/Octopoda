package service

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type Status struct {
	Name      string
	Platform  string
	CpuCores  int
	LocalTime time.Time

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
		Name:      config.GlobalConfig.Name,
		Platform:  GetCpuInfo(),
		CpuCores:  GetCpuCores(),
		LocalTime: time.Now(),

		CpuLoadShort: 0.0,
		CpuLoadLong:  0.0,
		MemUsed:      GetMemUsed(),
		MemTotal:     GetMemTotal(),
		DiskUsed:     GetDiskUsed(),
		DiskTotal:    GetDiskTotal(),
	}
	go maintainStatus()
}

func NodeStatus(conn net.Conn, serialNum uint32, raw []byte) {
	stateLock.RLock()
	state := nodeStatus
	state.LocalTime = time.Now()
	stateLock.RUnlock()
	serialized_info, _ := config.Jsoner.Marshal(&state)
	err := protocols.SendMessageUnique(conn, protocols.TypeNodeStatusResponse, serialNum, serialized_info)
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
