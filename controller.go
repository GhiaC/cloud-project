package main

import (
	log "github.com/sirupsen/logrus"
	vm "gogs.ghiasi.me/masoud/cloud-project/helper"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {
	servers := map[string]string{
		"vm1": "http://192.168.56.103:8080",
		"vm2": "http://192.168.56.104:8080",
		"vm3": "http://192.168.56.102:8080",
	}
	cpuUsage := map[string]int{
		"vm1": 0,
		"vm2": 0,
		"vm3": 0,
	}
	memUsage := map[string]int{
		"vm1": 0,
		"vm2": 0,
		"vm3": 0,
	}
	underPressure := map[string]bool{
		"vm1": false,
		"vm2": false,
		"vm3": false,
	}
	serverStatus := map[string]bool{
		"vm1": false,
		"vm2": false,
		"vm3": false,
	}
	// Get Cpu Usage
	for serverName, address := range servers {
		go func(serverName, address string) {
			for {
				resp, err := http.Get(address + "/GetCpuUsage")
				if err != nil {
					cpuUsage[serverName] = 0
					serverStatus[serverName] = false
					time.Sleep(5 * time.Second)
					continue
				}
				body, err := ioutil.ReadAll(resp.Body)
				_ = resp.Body.Close()
				if err != nil {
					cpuUsage[serverName] = 0
					serverStatus[serverName] = false
					time.Sleep(5 * time.Second)
					continue
				}
				cpu, _ := strconv.Atoi(string(body))
				serverStatus[serverName] = true
				cpuUsage[serverName] = cpu
			}
		}(serverName, address)
	}

	// Get Mem Usage
	for serverName, address := range servers {
		go func(serverName, address string) {
			for {
				resp, err := http.Get(address + "/GetMemUsage")
				if err != nil {
					memUsage[serverName] = 0
					serverStatus[serverName] = false
					time.Sleep(5 * time.Second)
					continue
				}
				body, err := ioutil.ReadAll(resp.Body)
				_ = resp.Body.Close()
				if err != nil {
					memUsage[serverName] = 0
					serverStatus[serverName] = false
					time.Sleep(5 * time.Second)
					continue
				}
				cpu, _ := strconv.Atoi(string(body))
				serverStatus[serverName] = true
				memUsage[serverName] = cpu
			}
		}(serverName, address)
	}

	// Get UnderPressure Usage
	for serverName, address := range servers {
		go func(serverName, address string) {
			for {
				resp, err := http.Get(address + "/UnderPressure")
				if err != nil {
					underPressure[serverName] = false
					serverStatus[serverName] = false
					time.Sleep(5 * time.Second)
					continue
				}
				body, err := ioutil.ReadAll(resp.Body)
				_ = resp.Body.Close()
				if err != nil {
					underPressure[serverName] = false
					serverStatus[serverName] = false
					time.Sleep(5 * time.Second)
					continue
				}
				if string(body) == "yes" {
					underPressure[serverName] = true
				} else {
					underPressure[serverName] = false
				}
				serverStatus[serverName] = true
			}
		}(serverName, address)
	}

	// Rule Checker
	ruleStatus := map[int]bool{
		1: false,
		2: false,
		3: false,
		4: false,
	}

	go func() {
		for {
			// RULE 1
			if serverStatus["vm1"] && !serverStatus["vm2"] && !serverStatus["vm3"] && !ruleStatus[1] {
				if cpuUsage["vm1"] > 80 {
					ruleStatus[1] = true
					go func() {
						timeCounter := 0
						for {
							if timeCounter >= 12 && cpuUsage["vm1"] < 80 {
								// Start VM2
								if _, err := vm.VboxCommandHandler("startvm", "vm2"); err != nil {
									log.Error(err)
								}
								ruleStatus[1] = false
								break
							}
							if cpuUsage["vm1"] < 80 {
								ruleStatus[1] = false
								break
							}
							timeCounter++
							time.Sleep(10 * time.Second)
						}
					}()
				}
			}
			// RULE 2
			if serverStatus["vm1"] && serverStatus["vm2"] && !serverStatus["vm3"] && !ruleStatus[2] {
				if (cpuUsage["vm1"]+cpuUsage["vm2"])/2 > 80 {
					ruleStatus[2] = true
					go func() {
						timeCounter := 0
						for {
							if timeCounter >= 12 && (cpuUsage["vm1"]+cpuUsage["vm2"])/2 < 80 {
								if _, err := vm.VboxCommandHandler("startvm", "vm3"); err != nil {
									log.Error(err)
								}
								ruleStatus[2] = false
								break
							}
							if (cpuUsage["vm1"]+cpuUsage["vm2"])/2 < 80 {
								ruleStatus[2] = false
								break
							}
							timeCounter++
							time.Sleep(10 * time.Second)
						}
					}()
				}
			}
			// RULE 3
			if serverStatus["vm1"] && serverStatus["vm2"] && serverStatus["vm3"] && !ruleStatus[3] {
				if (cpuUsage["vm1"]+cpuUsage["vm2"]+cpuUsage["vm3"])/3 < 50 {
					ruleStatus[3] = true
					go func() {
						timeCounter := 0
						for {
							if timeCounter >= 12 && (cpuUsage["vm1"]+cpuUsage["vm2"]+cpuUsage["vm3"])/3 < 50 {
								if _, err := vm.VboxCommandHandler("stopvm", "vm3"); err != nil {
									log.Error(err)
								}
								ruleStatus[3] = false
								break
							}
							if (cpuUsage["vm1"]+cpuUsage["vm2"]+cpuUsage["vm3"])/3 > 50 {
								ruleStatus[3] = false
								break
							}
							timeCounter++
							time.Sleep(10 * time.Second)
						}
					}()
				}
			}
			// RULE 4
			if serverStatus["vm1"] && serverStatus["vm2"] && !serverStatus["vm3"] && !ruleStatus[4] {
				if (cpuUsage["vm1"]+cpuUsage["vm2"])/2 < 40 {
					ruleStatus[4] = true
					go func() {
						timeCounter := 0
						for {
							if timeCounter >= 12 && (cpuUsage["vm1"]+cpuUsage["vm2"])/2 < 40 {
								if _, err := vm.VboxCommandHandler("stopvm", "vm2"); err != nil {
									log.Error(err)
								}
								ruleStatus[4] = false
								break
							}
							if (cpuUsage["vm1"]+cpuUsage["vm2"])/2 > 40 {
								ruleStatus[4] = false
								break
							}
							timeCounter++
							time.Sleep(10 * time.Second)
						}
					}()
				}
			}
		}
	}()
}
