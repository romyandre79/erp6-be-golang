package admin

import (
	"fmt"
	"runtime"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// System Stats
func GetSysInfo(c *fiber.Ctx) error {
	v, _ := mem.VirtualMemory()
	cStats, _ := cpu.Info()
	percent, _ := cpu.Percent(500*time.Millisecond, true)
	// Use current directory for disk stats (project folder)
	dStats, _ := disk.Usage("./")
	hInfo, _ := host.Info()

	var cpuModel string
	if len(cStats) > 0 {
		cpuModel = cStats[0].ModelName
	}

	return c.JSON(fiber.Map{
		"host": fiber.Map{
			"hostname":  hInfo.Hostname,
			"os":        hInfo.OS,
			"platform":  hInfo.Platform,
			"uptime":    hInfo.Uptime,
			"bootTime":  hInfo.BootTime,
			"procs":     hInfo.Procs,
			"goVersion": runtime.Version(),
			"numCpu":    runtime.NumCPU(),
		},
		"memory": fiber.Map{
			"total":       v.Total,
			"used":        v.Used,
			"free":        v.Free,
			"usedPercent": v.UsedPercent,
		},
		"cpu": fiber.Map{
			"model":   cpuModel,
			"percent": percent, // Array per core
		},
		"disk": fiber.Map{
			"path":        dStats.Path,
			"total":       dStats.Total,
			"free":        dStats.Free,
			"used":        dStats.Used,
			"usedPercent": dStats.UsedPercent,
		},
	})
}

// Process Explorer
func GetProcesses(c *fiber.Ctx) error {
	procs, _ := process.Processes()
	var list []fiber.Map

	for _, p := range procs {
		name, _ := p.Name()
		cpuP, _ := p.CPUPercent()
		memP, _ := p.MemoryPercent()
		username, _ := p.Username()
		cmd, _ := p.Cmdline()
		status, _ := p.Status()

		list = append(list, fiber.Map{
			"pid":     p.Pid,
			"name":    name,
			"cpu":     cpuP,
			"mem":     memP,
			"user":    username,
			"command": cmd,
			"status":  status, // Check structure
		})
	}

	// Sort by CPU usage desc by default
	sort.Slice(list, func(i, j int) bool {
		return list[i]["cpu"].(float64) > list[j]["cpu"].(float64)
	})

	// Limit to top 100 to avoid huge payload if not filtered
	if len(list) > 100 {
		list = list[:100]
	}

	return c.JSON(list)
}

// Kill Process
func KillProcess(c *fiber.Ctx) error {
	pid, err := c.ParamsInt("pid")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid PID"})
	}

	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Process not found"})
	}

	if err := p.Kill(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": fmt.Sprintf("Process %d killed", pid)})
}

// System File Manager (Project folder only, not system root)
func GetSystemRoot() string {
	// Return current working directory (project folder)
	// This is safer and more useful than system root
	return "./"
}

var SysFileMgr = NewMediaController(GetSystemRoot())

func SysListFiles(c *fiber.Ctx) error {
	return SysFileMgr.ListMedia(c)
}
