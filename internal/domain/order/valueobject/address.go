package valueobject

import (
	"errors"
	"fmt"
	"strings"
)

// Address 地址值对象
type Address struct {
	province      string // 省
	city          string // 市
	district      string // 区/县
	street        string // 街道
	detailAddress string // 详细地址
	postalCode    string // 邮政编码
	receiverName  string // 收货人姓名
	receiverPhone string // 收货人电话
}

// NewAddress 创建地址值对象
func NewAddress(
	province, city, district, street, detailAddress, postalCode, receiverName, receiverPhone string,
) (Address, error) {
	// 验证必填字段
	if province == "" {
		return Address{}, errors.New("省份不能为空")
	}
	if city == "" {
		return Address{}, errors.New("城市不能为空")
	}
	if detailAddress == "" {
		return Address{}, errors.New("详细地址不能为空")
	}
	if receiverName == "" {
		return Address{}, errors.New("收货人姓名不能为空")
	}
	if receiverPhone == "" {
		return Address{}, errors.New("收货人电话不能为空")
	}

	// 电话号码简单验证
	if len(strings.TrimSpace(receiverPhone)) < 11 {
		return Address{}, errors.New("收货人电话格式不正确")
	}

	return Address{
		province:      province,
		city:          city,
		district:      district,
		street:        street,
		detailAddress: detailAddress,
		postalCode:    postalCode,
		receiverName:  receiverName,
		receiverPhone: receiverPhone,
	}, nil
}

// Province 获取省份
func (a Address) Province() string {
	return a.province
}

// City 获取城市
func (a Address) City() string {
	return a.city
}

// District 获取区/县
func (a Address) District() string {
	return a.district
}

// Street 获取街道
func (a Address) Street() string {
	return a.street
}

// DetailAddress 获取详细地址
func (a Address) DetailAddress() string {
	return a.detailAddress
}

// PostalCode 获取邮政编码
func (a Address) PostalCode() string {
	return a.postalCode
}

// ReceiverName 获取收货人姓名
func (a Address) ReceiverName() string {
	return a.receiverName
}

// ReceiverPhone 获取收货人电话
func (a Address) ReceiverPhone() string {
	return a.receiverPhone
}

// FullAddress 获取完整地址
func (a Address) FullAddress() string {
	parts := []string{a.province, a.city}

	if a.district != "" {
		parts = append(parts, a.district)
	}

	if a.street != "" {
		parts = append(parts, a.street)
	}

	parts = append(parts, a.detailAddress)

	return strings.Join(parts, " ")
}

// String 地址的字符串表示
func (a Address) String() string {
	return fmt.Sprintf("%s, %s, %s", a.FullAddress(), a.receiverName, a.receiverPhone)
}

// Equals 判断两个地址是否相等
func (a Address) Equals(other Address) bool {
	return a.province == other.province &&
		a.city == other.city &&
		a.district == other.district &&
		a.street == other.street &&
		a.detailAddress == other.detailAddress &&
		a.postalCode == other.postalCode &&
		a.receiverName == other.receiverName &&
		a.receiverPhone == other.receiverPhone
}
