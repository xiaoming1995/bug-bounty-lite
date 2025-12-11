package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// 时间格式常量
const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
)

// DateTime 自定义时间类型，JSON 序列化时输出 "年-月-日 时:分:秒" 格式
type DateTime time.Time

// MarshalJSON 实现 json.Marshaler 接口
func (t DateTime) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", time.Time(t).Format(DateTimeFormat))
	return []byte(formatted), nil
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (t *DateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	// 去掉引号
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	parsed, err := time.ParseInLocation(DateTimeFormat, str, time.Local)
	if err != nil {
		// 尝试解析 ISO 格式
		parsed, err = time.Parse(time.RFC3339, str)
		if err != nil {
			return err
		}
	}
	*t = DateTime(parsed)
	return nil
}

// Value 实现 driver.Valuer 接口（写入数据库）
func (t DateTime) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return nil, nil
	}
	return time.Time(t), nil
}

// Scan 实现 sql.Scanner 接口（从数据库读取）
func (t *DateTime) Scan(value interface{}) error {
	if value == nil {
		*t = DateTime(time.Time{})
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*t = DateTime(v)
	case []byte:
		parsed, err := time.ParseInLocation(DateTimeFormat, string(v), time.Local)
		if err != nil {
			parsed, err = time.Parse(time.RFC3339, string(v))
			if err != nil {
				return err
			}
		}
		*t = DateTime(parsed)
	case string:
		parsed, err := time.ParseInLocation(DateTimeFormat, v, time.Local)
		if err != nil {
			parsed, err = time.Parse(time.RFC3339, v)
			if err != nil {
				return err
			}
		}
		*t = DateTime(parsed)
	default:
		return fmt.Errorf("cannot scan type %T into DateTime", value)
	}
	return nil
}

// Time 返回标准 time.Time
func (t DateTime) Time() time.Time {
	return time.Time(t)
}

// String 返回格式化的字符串
func (t DateTime) String() string {
	return time.Time(t).Format(DateTimeFormat)
}

