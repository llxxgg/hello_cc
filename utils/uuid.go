package utils

import "github.com/google/uuid"

// GenerateUUID 生成一个 UUID v4 (随机版本)
// 返回带连字符的标准格式 UUID 字符串
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateUUIDShort 生成一个不带连字符的 UUID 字符串
func GenerateUUIDShort() string {
	u := uuid.NewMD5(uuid.NameSpaceOID, []byte(uuid.New().String()))
	// 移除连字符
	uStr := u.String()
	result := make([]byte, 0, 32)
	for _, c := range uStr {
		if c != '-' {
			result = append(result, byte(c))
		}
	}
	return string(result)
}
