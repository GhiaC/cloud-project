package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {
	servers := map[string]string{
		"vm1": "http://192.168.56.4:8080",
		"vm2": "http://192.168.56.5:8080",
		"vm3": "http://192.168.56.6:8080",
	}
	cpuUsage := map[string]int{
		"vm1": 0,
		"vm2": 0,
		"vm3": 0,
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
			if serverStatus["vm1"] && !serverStatus["vm2"] && !ruleStatus[1] {
				if cpuUsage["vm1"] > 80 {
					ruleStatus[1] = true
					go func() {
						timeCounter := 0
						for {
							if timeCounter >= 12 && cpuUsage["vm1"] < 80 {
								// Start VM2
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
			if serverStatus["vm1"] && serverStatus["vm2"] && !ruleStatus[2] {
				if (cpuUsage["vm1"]+cpuUsage["vm2"])/2 > 80 {
					ruleStatus[2] = true
					go func() {
						timeCounter := 0
						for {
							if timeCounter >= 12 && (cpuUsage["vm1"]+cpuUsage["vm2"])/2 < 80 {
								// Start VM3
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
		}
	}()
}
