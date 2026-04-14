package utils

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"golang.org/x/crypto/bcrypt"
)

var (
	node   *snowflake.Node
	secret = []byte("夏天夏天悄悄过去")
)

// InitSnowflake 初始化雪花算法节点
// nodeID: 当前服务的节点 ID (0-1023)，在分布式系统中每个实例应不同
func InitSnowflake(nodeID int64) error {
	// 可以自定义起始时间（Epoch），例如 2026-01-01
	snowflake.Epoch = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

	var err error
	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		return fmt.Errorf("failed to init snowflake node: %w", err)
	}
	return nil
}

// GenerateUserID 生成一个 64 位的唯一 ID
func GenerateUserID() int64 {
	if node == nil {
		// 如果未初始化，默认使用节点 1（建议在服务启动时显式调用 Init）
		_ = InitSnowflake(1)
	}
	return node.Generate().Int64()
}

// PasswordHash 对明文密码进行加密（自动加盐）
func PasswordHash(password string) (string, error) {
	// GenerateFromPassword 会自动生成随机盐值并混入结果中
	// DefaultCost 是 10，对于目前的硬件性能来说是一个很好的平衡点
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// PasswordVerify 验证明文密码与哈希值是否匹配
func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
