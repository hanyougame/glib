package uniqueid

import (
	"errors"
	"fmt"
	"github.com/sony/sonyflake"
	"net"
	"strings"
)

var flake *sonyflake.Sonyflake

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: getMachineID,
	})
}

// GenId 生成一个唯一的雪花ID
func GenId() (id uint64, err error) {
	id, err = flake.NextID()
	return
}

// 获取机器 ID 基于 MAC 地址
func getMachineID() (uint16, error) {
	// 获取所有网络接口的地址
	interfaces, err := net.Interfaces()
	if err != nil {
		return 0, fmt.Errorf("failed to get network interfaces: %v", err)
	}

	// 查找第一个有效的网卡，并获取其 MAC 地址
	for _, iFace := range interfaces {
		// 忽略环回地址和没有 MAC 地址的接口
		if iFace.Flags&net.FlagUp == 0 || iFace.HardwareAddr == nil || len(iFace.HardwareAddr) == 0 {
			continue
		}

		// 首选外网网卡的 MAC 地址（例如非本地网卡）
		if isValidPublicInterface(iFace) {
			return uint16(sum(iFace.HardwareAddr) % 1024), nil
		}
	}

	// 如果没有找到合适的外网网卡，则使用第一个有效的网卡
	for _, iFace := range interfaces {
		if iFace.Flags&net.FlagUp == 0 || iFace.HardwareAddr == nil || len(iFace.HardwareAddr) == 0 {
			continue
		}
		return uint16(sum(iFace.HardwareAddr) % 1024), nil
	}

	return 0, errors.New("no valid network interface with MAC address found")
}

// 判断接口是否有效且为外网网卡
func isValidPublicInterface(iFace net.Interface) bool {
	// 排除环回接口（lo）
	if strings.HasPrefix(iFace.Name, "lo") {
		return false
	}
	return true
}

// 计算 MAC 地址的字节和作为机器 ID
func sum(mac net.HardwareAddr) int {
	total := 0
	for _, b := range mac {
		total += int(b)
	}
	return total
}
