package main

import (
	"os/exec"
	"log"
	"strconv"
	"bufio"
	"sync"
	"strings"
	"time"
)

type SpeedTesterCLI struct{}

func (speedTester *SpeedTesterCLI) MeasureBandwidth() (result *SpeedTestResult, err error) {
	result = &SpeedTestResult{
		Timestamp: time.Now(),
	}
	command := exec.Command("speedtest-cli", "--simple")
	cmdReader, err := command.StdoutPipe()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(cmdReader)
	waitGroup := &sync.WaitGroup{}
	go func() {
		var parseError error
		waitGroup.Add(1)
		defer waitGroup.Done()
		for scanner.Scan() {
			split := strings.SplitN(scanner.Text(), ` `, 2)
			switch split[0] {
			case "Ping:":
				result.Ping, parseError = speedTester.parsePing(split[1])
				if parseError != nil {
					log.Printf("[WARNING] Could not parse ping (%v): %v", strconv.Quote(split[1]), parseError)
				}
				break
			case "Download:":
				result.Download, parseError = speedTester.parseBandwidthValue(split[1])
				if parseError != nil {
					log.Printf("[WARNING] Could not parse download value (%v): %v", strconv.Quote(split[1]), parseError)
				}
				break
			case "Upload:":
				result.Upload, parseError = speedTester.parseBandwidthValue(split[1])
				if parseError != nil {
					log.Printf("[WARNING] Could not parse upload value (%v): %v", strconv.Quote(split[1]), parseError)
				}
				break
			default:
				log.Printf("[FATAL] Received unknown value from speedtest-cli command.")
			}
		}
	}()
	defer cmdReader.Close()
	if err := command.Run(); err != nil {
		panic(err)
	}
	waitGroup.Wait()
	return
}

func (speedTester *SpeedTesterCLI) parsePing(value string) (ping time.Duration, err error) {
	rawPing := strings.Replace(value, " ", "", -1)
	ping, err = time.ParseDuration(rawPing)
	return
}

func (speedTester *SpeedTesterCLI) parseBandwidthValue(value string) (bandwidth Bandwidth, err error) {
	rawPing := strings.Replace(value, " ", "", -1)
	rawPing = strings.Replace(rawPing, "Mbit/s", "", -1)
	var bandwidth64 float64
	bandwidth64, err = strconv.ParseFloat(rawPing, 32)
	bandwidth = Bandwidth(bandwidth64)
	return
}
