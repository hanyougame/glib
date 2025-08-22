package uniqueid

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sony/sonyflake"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

var (
	flake      *sonyflake.Sonyflake
	randSource rand.Source
)

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: getMachineID,
	})
	if flake == nil {
		// 如果初始化失败，使用默认设置
		flake = sonyflake.NewSonyflake(sonyflake.Settings{})
		fmt.Println("Warning: Sonyflake initialization failed, using default settings")
	}
	// 使用当前时间戳初始化随机源
	randSource = rand.NewSource(time.Now().UnixNano())
}

// GenId 生成一个唯一的雪花ID
func GenId() (id uint64, err error) {
	id, err = flake.NextID()
	return
}

func GenUserID(ctx context.Context, rdb redis.UniversalClient, key string) (id int64, err error) {
	// 使用 Redis 自增 ID
	id, err = rdb.Incr(ctx, key).Result()
	return
}

// 获取机器 ID 基于 Docker 环境
func getMachineID() (uint16, error) {
	// 判断是否在 Docker 环境中运行
	if isRunningInDocker() {
		// 尝试通过容器 ID 生成机器 ID
		containerID, err := getContainerID()
		if err != nil {
			fmt.Println("getContainerID err: ", err)
			return 0, fmt.Errorf("failed to get container ID: %v", err)
		}
		return uint16(sum([]byte(containerID)) % 1024), nil
	}

	// 如果不在 Docker 环境中，继续使用 MAC 地址方式
	return getMachineIDFromMac()
}

// 判断是否在 Docker 容器中运行
func isRunningInDocker() bool {
	// 检查容器的特征文件
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	// 检查 cgroup 信息
	data, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "docker")
}

// 获取容器 ID（适用于 Docker 环境）
func getContainerID() (string, error) {
	// 获取容器 ID
	// 一般情况下，可以通过读取 `/proc/self/cgroup` 获取容器 ID
	data, err := os.ReadFile("/proc/self/cgroup")
	if err != nil {
		return "", fmt.Errorf("failed to read /proc/self/cgroup: %v", err)
	}

	// 从文件内容中提取容器 ID
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, "docker") {
			parts := strings.Split(line, "/")
			if len(parts) > 2 {
				return parts[len(parts)-1], nil
			}
		}
	}
	return "", errors.New("container ID not found")
}

// 获取机器 ID 基于 MAC 地址（不在 Docker 环境下）
func getMachineIDFromMac() (uint16, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return 0, fmt.Errorf("failed to get network interfaces: %v", err)
	}

	// 查找第一个有效的网卡并获取其 MAC 地址
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.HardwareAddr == nil || len(iface.HardwareAddr) == 0 {
			continue
		}
		return uint16(sum(iface.HardwareAddr) % 1024), nil
	}

	return 0, errors.New("no valid network interface with MAC address found")
}

// 计算字节数组的和作为机器 ID
func sum(data []byte) int {
	total := 0
	for _, b := range data {
		total += int(b)
	}
	return total
}
