package main

import (
	"fmt"
	"github.com/ghiac/go-commons/signal"
	vm "gogs.ghiasi.me/masoud/cloud-project/helper"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Handler struct {
	ShutdownTime  int
	Servers       map[string]string
	CpuUsage      map[string]int
	memUsage      map[string]int
	underPressure map[string]bool
	serverStatus  map[string]bool
	ruleStatus    map[int]bool
	mutex         map[string]*sync.Mutex
}

func main() {
	h := &Handler{
		ShutdownTime: 20,
		mutex: map[string]*sync.Mutex{
			"serverStatus": &sync.Mutex{},
			"ruleStatus":   &sync.Mutex{},
		},
	}

	h.Servers = map[string]string{
		"vm1": "http://192.168.56.103:9090",
		"vm2": "http://192.168.56.104:9090",
		"vm3": "http://192.168.56.102:9090",
	}
	h.CpuUsage = map[string]int{
		"vm1": 0,
		"vm2": 0,
		"vm3": 0,
	}
	h.memUsage = map[string]int{
		"vm1": 0,
		"vm2": 0,
		"vm3": 0,
	}
	h.underPressure = map[string]bool{
		"vm1": false,
		"vm2": false,
		"vm3": false,
	}
	h.serverStatus = map[string]bool{
		"vm1": false,
		"vm2": false,
		"vm3": false,
	}
	// Get Cpu Usage

	go func() {
		i := 0
		for {
			i++
			fmt.Println("Seconds:", i)
			fmt.Println("Servers status:", h.serverStatus)
			fmt.Println("Cpu Usage:     ", h.CpuUsage)
			fmt.Println("Mem Usage:     ", h.memUsage)
			fmt.Println("Under pressure:", h.underPressure)
			fmt.Println("Rules status:  ", h.ruleStatus)
			fmt.Println("----------------------")
			time.Sleep(1 * time.Second)
		}
	}()

	// Rule Checker
	h.ruleStatus = map[int]bool{
		1: false,
		2: false,
		3: false,
		4: false,
	}

	h.checker()
	h.memChecker()
	h.underPressureChecker()
	h.starter()
	signal.Signal.Wait()
}

func (h *Handler) starter() {
	for {
		// TODO get variables

		// RULE 1
		if h.serverStatus["vm1"] && !h.serverStatus["vm2"] && !h.serverStatus["vm3"] && !h.ruleStatus[1] {
			if h.CpuUsage["vm1"] > 60 {
				h.mutex["ruleStatus"].Lock()
				h.ruleStatus[1] = true
				h.mutex["ruleStatus"].Unlock()
				go func() {
					timeCounter := 0
					for {
						if timeCounter >= h.ShutdownTime && h.CpuUsage["vm1"] > 60 {
							// Start VM2
							if _, err := vm.VboxCommandHandler("startvm", "vm2"); err != nil {

							}
							h.mutex["ruleStatus"].Lock()
							h.ruleStatus[1] = false
							h.mutex["ruleStatus"].Unlock()
							break
						}
						if h.CpuUsage["vm1"] < 60 {
							h.mutex["ruleStatus"].Lock()
							h.ruleStatus[1] = false
							h.mutex["ruleStatus"].Unlock()
							break
						}
						timeCounter++
						time.Sleep(1 * time.Second)
					}
				}()
			}
		}
		// RULE 2
		if h.serverStatus["vm1"] && h.serverStatus["vm2"] && !h.serverStatus["vm3"] && !h.ruleStatus[2] {
			if (h.CpuUsage["vm1"]+h.CpuUsage["vm2"])/2 > 60 {
				h.mutex["ruleStatus"].Lock()
				h.ruleStatus[2] = true
				h.mutex["ruleStatus"].Unlock()
				go func() {
					timeCounter := 0
					for {
						if timeCounter >= h.ShutdownTime && (h.CpuUsage["vm1"]+h.CpuUsage["vm2"])/2 > 60 {
							if _, err := vm.VboxCommandHandler("startvm", "vm3"); err != nil {

							}
							h.mutex["ruleStatus"].Lock()
							h.ruleStatus[2] = false
							h.mutex["ruleStatus"].Unlock()
							break
						}
						if ((h.CpuUsage["vm1"] + h.CpuUsage["vm2"]) / 2) < 60 {
							h.mutex["ruleStatus"].Lock()
							h.ruleStatus[2] = false
							h.mutex["ruleStatus"].Unlock()
							break
						}
						timeCounter++
						time.Sleep(1 * time.Second)
					}
				}()
			}
		}
		// RULE 3
		if h.serverStatus["vm1"] && h.serverStatus["vm2"] && h.serverStatus["vm3"] && !h.ruleStatus[3] {
			if ((h.CpuUsage["vm1"] + h.CpuUsage["vm2"] + h.CpuUsage["vm3"]) / 3) < 50 {
				h.mutex["ruleStatus"].Lock()
				h.ruleStatus[3] = true
				h.mutex["ruleStatus"].Unlock()
				go func() {
					timeCounter := 0
					for {
						if timeCounter >= h.ShutdownTime && ((h.CpuUsage["vm1"]+h.CpuUsage["vm2"]+h.CpuUsage["vm3"])/3) < 50 {
							if _, err := vm.VboxCommandHandler("controlvm", "vm3", "poweroff"); err != nil {
							}
							h.mutex["ruleStatus"].Lock()
							h.ruleStatus[3] = false
							h.mutex["ruleStatus"].Unlock()
							break
						}
						if ((h.CpuUsage["vm1"] + h.CpuUsage["vm2"] + h.CpuUsage["vm3"]) / 3) > 50 {
							h.mutex["ruleStatus"].Lock()
							h.ruleStatus[3] = false
							h.mutex["ruleStatus"].Unlock()
							break
						}
						timeCounter++
						time.Sleep(1 * time.Second)
					}
				}()
			}
		}
		// RULE 4
		if h.serverStatus["vm1"] && h.serverStatus["vm2"] && !h.serverStatus["vm3"] && !h.ruleStatus[4] {
			if ((h.CpuUsage["vm1"] + h.CpuUsage["vm2"]) / 2) < 40 {
				h.mutex["ruleStatus"].Lock()
				h.ruleStatus[4] = true
				h.mutex["ruleStatus"].Unlock()
				go func() {
					timeCounter := 0
					for {
						if timeCounter >= h.ShutdownTime && ((h.CpuUsage["vm1"]+h.CpuUsage["vm2"])/2) < 40 {
							if _, err := vm.VboxCommandHandler("controlvm", "vm2", "poweroff"); err != nil {
							}
							h.mutex["ruleStatus"].Lock()
							h.ruleStatus[4] = false
							h.mutex["ruleStatus"].Unlock()
							break
						}
						if ((h.CpuUsage["vm1"] + h.CpuUsage["vm2"]) / 2) > 40 {
							h.mutex["ruleStatus"].Lock()
							h.ruleStatus[4] = false
							h.mutex["ruleStatus"].Unlock()
							break
						}
						timeCounter++
						time.Sleep(1 * time.Second)
					}
				}()
			}
		}
		// RULE 5
		if h.serverStatus["vm1"] && !h.serverStatus["vm2"] && !h.serverStatus["vm3"] {
			if h.underPressure["vm1"] {
				if _, err := vm.VboxCommandHandler("startvm", "vm2"); err != nil {
				}
			}
		}
		// RULE 6
		if h.serverStatus["vm1"] && h.serverStatus["vm2"] && !h.serverStatus["vm3"] {
			if h.underPressure["vm2"] {
				if _, err := vm.VboxCommandHandler("startvm", "vm3"); err != nil {
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (h *Handler) checker() {
	for serverName, address := range h.Servers {
		address := address
		serverName := serverName
		go func() {
			client := &http.Client{}
			client.Timeout = time.Second
			client.CloseIdleConnections()
			for {
				url := address + "/GetCpuUsage"
				req, err := http.NewRequest("GET", url, nil)
				resp, err := client.Do(req)
				if err != nil {
					h.mutex["serverStatus"].Lock()
					h.CpuUsage[serverName] = 0
					h.serverStatus[serverName] = false
					h.mutex["serverStatus"].Unlock()
					time.Sleep(5 * time.Second)
					continue
				}
				body, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				if err != nil {
					h.mutex["serverStatus"].Lock()
					h.CpuUsage[serverName] = 0
					h.serverStatus[serverName] = false
					h.mutex["serverStatus"].Unlock()
					time.Sleep(5 * time.Second)
					continue
				}
				h.mutex["serverStatus"].Lock()
				cpu, _ := strconv.Atoi(string(body))
				h.serverStatus[serverName] = true
				h.CpuUsage[serverName] = cpu
				h.mutex["serverStatus"].Unlock()
				time.Sleep(1 * time.Second)
			}
		}()
	}
}

func (h *Handler) memChecker() {
	for serverName, address := range h.Servers {
		address := address
		serverName := serverName
		go func() {
			client := &http.Client{}
			client.Timeout = time.Second
			client.CloseIdleConnections()
			for {
				url := address + "/GetMemUsage"
				req, err := http.NewRequest("GET", url, nil)
				resp, err := client.Do(req)
				if err != nil {
					h.memUsage[serverName] = 0
					time.Sleep(5 * time.Second)
					continue
				}
				body, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				if err != nil {
					h.memUsage[serverName] = 0
					continue
				}
				cpu, _ := strconv.Atoi(string(body))
				h.memUsage[serverName] = cpu
				time.Sleep(1 * time.Second)
			}
		}()
	}
}

func (h *Handler) underPressureChecker() {
	for serverName, address := range h.Servers {
		address := address
		serverName := serverName
		go func() {
			client := &http.Client{}
			client.Timeout = time.Second
			client.CloseIdleConnections()
			for {
				url := address + "/UnderPressure"
				req, err := http.NewRequest("GET", url, nil)
				resp, err := client.Do(req)
				if err != nil {
					h.underPressure[serverName] = false
					time.Sleep(5 * time.Second)
					continue
				}
				body, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				if err != nil {
					h.underPressure[serverName] = false
					continue
				}
				if string(body) == "yes" {
					h.underPressure[serverName] = true
				} else {
					h.underPressure[serverName] = false
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}
}
