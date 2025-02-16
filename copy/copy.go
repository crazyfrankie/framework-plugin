package copy

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Copier struct {
	src any
	dst any

	ignoreFields map[string]struct{}
}

func NewCopier(src any, dst any, ignoreFields ...string) Copier {
	ignore := make(map[string]struct{}, len(ignoreFields))
	for _, fd := range ignoreFields {
		ignore[fd] = struct{}{}
	}
	return Copier{
		src:          src,
		dst:          dst,
		ignoreFields: ignore,
	}
}

func (c Copier) Builder() error {
	// 获取 src dst 的反射值
	srcVal := reflect.ValueOf(c.src)
	dstVal := reflect.ValueOf(c.dst)

	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}
	if dstVal.Kind() == reflect.Ptr {
		dstVal = dstVal.Elem()
	}

	// 确保 src 和 dst 都是结构体类型
	if srcVal.Kind() != reflect.Struct || dstVal.Kind() != reflect.Struct {
		return errors.New("src and dst must be structs")
	}

	// 遍历 src 的字段
	numField := srcVal.NumField()
	for i := 0; i < numField; i++ {
		// 获取字段名称和对应值
		srcField := srcVal.Type().Field(i)
		srcFieldVal := srcVal.Field(i)

		if _, exists := c.ignoreFields[srcField.Name]; exists {
			continue
		}

		dstField := dstVal.FieldByName(srcField.Name)
		if dstField.IsValid() && dstField.CanSet() {
			// 判断字段类型是否匹配，如果不匹配则尝试转换
			if srcFieldVal.Type().AssignableTo(dstField.Type()) {
				dstField.Set(srcFieldVal)
			} else {
				// 如果类型不匹配但可以转换，尝试进行类型转换
				if err := c.convertAndSet(srcFieldVal, dstField); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// convertAndSet 尝试进行类型转换并设置
func (c Copier) convertAndSet(srcField, dstField reflect.Value) error {
	if srcField.Kind() == reflect.Int || srcField.Kind() == reflect.Int8 || srcField.Kind() == reflect.Int16 || srcField.Kind() == reflect.Int32 || srcField.Kind() == reflect.Int64 {
		if dstField.Kind() == reflect.Uint || dstField.Kind() == reflect.Uint8 || dstField.Kind() == reflect.Uint16 || dstField.Kind() == reflect.Uint32 || dstField.Kind() == reflect.Uint64 {
			dstField.SetUint(uint64(srcField.Int()))
			return nil
		}
		if dstField.Kind() == reflect.Float32 || dstField.Kind() == reflect.Float64 {
			dstField.SetFloat(float64(srcField.Int()))
			return nil
		}
		if dstField.Kind() == reflect.Int || dstField.Kind() == reflect.Int8 || dstField.Kind() == reflect.Int16 || dstField.Kind() == reflect.Int32 || dstField.Kind() == reflect.Int64 {
			dstField.SetInt(srcField.Int())
			return nil
		}
	}

	if srcField.Kind() == reflect.Int || srcField.Kind() == reflect.Int8 || srcField.Kind() == reflect.Int16 || srcField.Kind() == reflect.Int32 || srcField.Kind() == reflect.Int64 {
		if dstField.Kind() == reflect.Float32 || dstField.Kind() == reflect.Float64 {
			dstField.SetFloat(float64(srcField.Int()))
			return nil
		}
	}

	if dstField.Kind() == reflect.String {
		if srcField.Kind() == reflect.Int || srcField.Kind() == reflect.Int8 || srcField.Kind() == reflect.Int16 || srcField.Kind() == reflect.Int32 || srcField.Kind() == reflect.Int64 {
			dstField.SetString(strconv.FormatInt(srcField.Int(), 10))
			return nil
		}
		if srcField.Kind() == reflect.Float32 || srcField.Kind() == reflect.Float64 {
			dstField.SetString(fmt.Sprintf("%f", srcField.Float()))
			return nil
		}
	}

	// 如果无法转换，返回类型不匹配的错误
	return fmt.Errorf("cannot convert %s to %s", srcField.Type(), dstField.Type())
}

type A struct {
	ID       int64
	Name     string
	Password string
	Avatar   string
	Phone    string
	Ctime    int64
	Utime    int64
}

type B struct {
	ID     int64
	Name   string
	Avatar string
	Phone  string
	Ctime  int64
	Utime  int64
}
