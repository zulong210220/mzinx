/*
 * File : id_generator.go
 * CreateDate : 2019-12-15 20:48:49
 * */
package utils

import (
	"errors"
	"net"
	"sync"
	"time"
)

var (
	startTime   = int64(1474992000000)
	machineId   = uint16(0)
	elapsedTime = int64(0)
	sequence    = uint16(0)
	mutex       *sync.Mutex
)

const (
	BIT_LEN_SEQUENCE   = 12
	BIT_LEN_MACHINE_ID = 10
	BIT_LEN_TIME       = 41
)

const (
	snowFlakeTimeUnit = 1e6
	maskSequence      = uint16(1<<BIT_LEN_SEQUENCE - 1)
)

func init() {
	mutex = &sync.Mutex{}
	SetMachineId(0)
}

func getLower16BitPrivateIp() (uint16, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return 0, err
	}

	var ipv4 net.IP
	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if ip != nil && (ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168) {
			ipv4 = ip
			break
		}
	}

	if ipv4 == nil {
		return 0, errors.New("no private ip address")
	}

	return uint16(ipv4[2]>>(BIT_LEN_MACHINE_ID-8) + ipv4[3]<<(BIT_LEN_MACHINE_ID-8)), nil
}

func SetMachineId(mid uint16) bool {
	maxMachineId := uint16(1<<BIT_LEN_MACHINE_ID - 1)
	if maxMachineId < mid {
		return false
	}

	if mid == 0 {
		var err error
		machineId, err = getLower16BitPrivateIp()
		if err != nil {
			panic(err)
		}
	} else {
		machineId = mid
	}
	return true
}

func toSnowFlakeTime(t time.Time) int64 {
	return t.UTC().UnixNano() / snowFlakeTimeUnit
}

func currentElapsedTime() int64 {
	return toSnowFlakeTime(time.Now().UTC()) - startTime
}

func sleepTime(overTime int64) time.Duration {
	return time.Duration(overTime)*time.Millisecond - time.Duration(time.Now().UTC().UnixNano()%snowFlakeTimeUnit)*time.Nanosecond
}

func toId() (int64, error) {
	if elapsedTime >= 1<<BIT_LEN_TIME {
		return 0, errors.New("elapsed time over time limit")
	}
	return int64(elapsedTime<<(BIT_LEN_SEQUENCE+BIT_LEN_MACHINE_ID)) | int64(machineId<<BIT_LEN_SEQUENCE) | int64(sequence), nil
}

func NextId() (int64, error) {
	mutex.Lock()
	defer mutex.Unlock()

	current := currentElapsedTime()

	if elapsedTime < current {
		elapsedTime = current
		sequence = 0
	} else {
		sequence = (sequence + 1) & maskSequence
		if sequence == 0 {
			elapsedTime++
			overTime := elapsedTime - current
			time.Sleep(sleepTime(overTime))
		}
	}

	return toId()
}

/* vim: set tabstop=4 set shiftwidth=4 */
