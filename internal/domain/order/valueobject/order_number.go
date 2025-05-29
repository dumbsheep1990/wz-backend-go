package valueobject

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

// OrderNumber 订单编号值对象
type OrderNumber string

// 订单编号格式正则表达式：WZ + 年月日 + 6位随机数，例如WZ2023052912345
var orderNumberRegex = regexp.MustCompile(`^WZ\d{8}\d{6}$`)

// NewOrderNumber 创建订单编号值对象
func NewOrderNumber(number string) (OrderNumber, error) {
	if number == "" {
		return "", errors.New("订单编号不能为空")
	}

	if !orderNumberRegex.MatchString(number) {
		return "", errors.New("订单编号格式不正确，应为WZ+年月日+6位随机数")
	}

	return OrderNumber(number), nil
}

// GenerateOrderNumber 生成新的订单编号
func GenerateOrderNumber() OrderNumber {
	// 格式：WZ + 年月日 + 6位随机数
	now := time.Now()
	dateStr := now.Format("20060102")
	randomNum := fmt.Sprintf("%06d", now.UnixNano()%1000000) // 使用纳秒时间戳的后6位作为随机数
	number := fmt.Sprintf("WZ%s%s", dateStr, randomNum)
	return OrderNumber(number)
}

// Value 获取编号值
func (n OrderNumber) Value() string {
	return string(n)
}

// String 字符串表示
func (n OrderNumber) String() string {
	return string(n)
}
