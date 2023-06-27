package sys

import (
	"brain/config"
	"brain/logger"
	"brain/model"
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

var nodeStatus model.Status
var stateLock sync.RWMutex

func InitNodeStatus() {
	nodeStatus = model.Status{
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

func LocalStatus() model.Status {
	stateLock.RLock()
	state := nodeStatus
	state.LocalTime = time.Now()
	stateLock.RUnlock()
	return state
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
