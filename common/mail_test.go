package common

import (
	"testing"
)

func TestSendMail(t *testing.T) {
	SendMail("测试主题", " 这是一个测试内容")
}
