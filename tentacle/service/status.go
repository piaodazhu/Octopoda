package service

import (
	"fmt"
	"net"
	"sync"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/message"
	"tentacle/snp"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
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

func NodeStatus(conn net.Conn, raw []byte) {
	stateLock.RLock()
	state := nodeStatus
	state.LocalTime = time.Now()
	stateLock.RUnlock()
	text := statusToText(state)
	serialized_info, _ := config.Jsoner.Marshal(&text)
	err := message.SendMessageUnique(conn, message.TypeNodeStatusResponse, snp.GenSerial(), serialized_info)
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

type StatusText struct {
	Name         string `json:"name"`
	Platform     string `json:"platform"`
	CpuCores     int    `json:"cpu_cores"`
	LocalTime    string `json:"local_time"`
	CpuLoadShort string `json:"cpu_average1"`
	CpuLoadLong  string `json:"cpu_average10"`
	MemUsage     string `json:"memory_usage"`
	DiskUsage    string `json:"disk_usage"`
}

func statusToText(status Status) StatusText {
	return StatusText{
		Name:         status.Name,
		Platform:     status.Platform,
		CpuCores:     status.CpuCores,
		LocalTime:    status.LocalTime.Format("2006-01-02 15:04:05"),
		CpuLoadShort: fmt.Sprintf("%5.1f%%", status.CpuLoadShort),
		CpuLoadLong:  fmt.Sprintf("%5.1f%%", status.CpuLoadLong),
		MemUsage: fmt.Sprintf("%5.1f%%: (%.2fGB / %.2fGB)",
			float64(status.MemUsed*100)/float64(status.MemTotal),
			float64(status.MemUsed)/1073741824,
			float64(status.MemTotal)/1073741824),
		DiskUsage: fmt.Sprintf("%5.1f%%: (%.2fGB / %.2fGB)",
			float64(status.DiskUsed*100)/float64(status.DiskTotal),
			float64(status.DiskUsed)/1073741824,
			float64(status.DiskTotal)/1073741824),
	}
}
