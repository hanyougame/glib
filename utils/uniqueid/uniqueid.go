package uniqueid

import (
	"fmt"
	"github.com/sony/sonyflake"
	"net"
)

var flake *sonyflake.Sonyflake

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: getMachineID,
	})
}

func GenId() (id uint64, err error) {
	id, err = flake.NextID()
	return
}

// 获取机器id
func getMachineID() (uint16, error) {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		return 0, err
	}
	// 取第一个有效的非环回 IP 地址
	for _, addr := range addrList {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			return uint16(sum(ip.IP)), nil
		}
	}
	return 0, fmt.Errorf("no valid IP address found")
}

// 计算 IP 地址的和作为机器 ID
func sum(ip net.IP) int {
	total := 0
	for _, b := range ip {
		total += int(b)
	}
	return total % 1024 // 限制在最大机器 ID 范围内
}
