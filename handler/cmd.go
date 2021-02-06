package handler

import (
	"bufio"
	"fmt"
	"github.com/mackerelio/go-osstat/memory"
	"gogs.ghiasi.me/masoud/cloud-project/config"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewHandler() *Handler {
	return &Handler{
		AvgCpuUsage:           0,
		checkSeconds:          0,
		memoryUsage:           0,
		percentageMemoryUsage: 0,
		minimumUsage:          0,
		maxSecondsOfAvg:       120,
	}
}

type Handler struct {
	AvgCpuUsage           float64
	currentCpuUsage       float64
	checkSeconds          int
	percentageMemoryUsage int
	minimumUsage          float64
	memoryUsage           uint64
	maxSecondsOfAvg       int
}

func (h *Handler) Start() {
	h.cpuChecker()
}

func (h *Handler) cpuChecker() {
	go func() {
		flag := false
		var prevIdleTime, prevTotalTime uint64
		for {
			file, err := os.Open("/proc/stat")
			if err != nil {
				log.Fatal(err)
			}
			scanner := bufio.NewScanner(file)
			scanner.Scan()
			firstLine := scanner.Text()[5:] // get rid of cpu plus 2 spaces
			file.Close()
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
			split := strings.Fields(firstLine)
			idleTime, _ := strconv.ParseUint(split[3], 10, 64)
			totalTime := uint64(0)
			for _, s := range split {
				u, _ := strconv.ParseUint(s, 10, 64)
				totalTime += u
			}
			cpuUsage := float64(0)
			if flag {
				deltaIdleTime := idleTime - prevIdleTime
				deltaTotalTime := totalTime - prevTotalTime
				cpuUsage = (1.0 - float64(deltaIdleTime)/float64(deltaTotalTime)) * 100.0
				if h.checkSeconds > h.maxSecondsOfAvg {
					h.AvgCpuUsage = ((h.AvgCpuUsage * float64(h.maxSecondsOfAvg-1)) + cpuUsage) / float64(h.maxSecondsOfAvg)
				} else {
					h.AvgCpuUsage = ((h.AvgCpuUsage * float64(h.checkSeconds)) + cpuUsage) / float64(h.checkSeconds+1)
					h.checkSeconds++
				}
				h.currentCpuUsage = cpuUsage
				//fmt.Printf("cpu usage : %6.3f \n", h.AvgCpuUsage)
			} else {
				flag = true
			}
			prevIdleTime = idleTime
			prevTotalTime = totalTime
			if cpuUsage < h.minimumUsage {
				go h.IncreaseLoad(int(h.minimumUsage-cpuUsage) * 10000)
			}
			h.PrintMemUsage()
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (h *Handler) PrintMemUsage() {
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	h.memoryUsage = bToMb(memory.Used)
	h.percentageMemoryUsage = int((float64(memory.Used) / float64(memory.Total)) * 100)
	//fmt.Printf("memory used: %d mg\n", h.memoryUsage)
	//fmt.Printf("percentageMemoryUsage used: %d %%\n", h.percentageMemoryUsage)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func (h *Handler) GetMemAvg() string {
	return strconv.Itoa(h.percentageMemoryUsage)
}

func (h *Handler) GetCpuAvg() string {
	return strconv.Itoa(int(h.AvgCpuUsage))
}

func (h *Handler) IncreaseLoad(n int) {
	go func() {
		if int(h.currentCpuUsage) > config.Conf.Core.MaxCpuUsage {
			return
		}
		//high := 50000000000
		high := 500000 * n
		for i := high; i > 0; i-- {
			r := high * high
			r--
		}
	}()
}

func (h *Handler) IncreaseMemUsage(n int) {
	go func() {
		//high := 50000000000
		if h.percentageMemoryUsage > config.Conf.Core.MaxMemUsage {
			return
		}
		for i := n; i > 0; i-- {
			go func() {
				temp := make([]string, 20000000)
				temp = temp
				time.Sleep(10 * time.Second)
			}()
		}
	}()
}
