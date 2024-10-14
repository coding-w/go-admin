package utils

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Rules map[string][]string

type RulesMap map[string]Rules

var CustomizeMap = make(map[string]Rules)

// RegisterRule
// @description: 注册自定义规则
// @return: error
func RegisterRule(key string, rule Rules) (err error) {
	if CustomizeMap[key] != nil {
		return errors.New(key + "已注册,无法重复注册")
	} else {
		CustomizeMap[key] = rule
		return nil
	}
}

// NotEmpty
// @description: 验证数据是否不为空
// @return: string
func NotEmpty() string {
	return "notEmpty"
}

// RegexpMatch
// @description: 验证数据是否匹配正则
// @param: rule string
func RegexpMatch(rule string) string {
	return "regexp=" + rule
}

// Lt
// @description: 小于入参(<) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
// @param: mark string
func Lt(mark string) string {
	return "lt=" + mark
}

// Le
// @description: 小于等于入参(<=) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
// @param: mark string
func Le(mark string) string {
	return "le=" + mark
}

// Eq
// @description: 等于入参(==) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
// @param: mark string
func Eq(mark string) string {
	return "eq=" + mark
}

// Ne
// @description: 不等于入参(!=) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
// @param: mark string
func Ne(mark string) string {
	return "ne=" + mark
}

// Ge
// @description: 大于等于入参(>=) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
// @param: mark string
func Ge(mark string) string {
	return "ge=" + mark
}

// Gt
// @description: 大于入参(>) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
// @param: mark string
func Gt(mark string) string {
	return "gt=" + mark
}

func Verify(st interface{}, roleMap Rules) (err error) {
	compareMap := map[string]bool{
		"eq": true,
		"ne": true,
		"lt": true,
		"le": true,
		"gt": true,
		"ge": true,
	}
	// 获取结构体类型
	typ := reflect.TypeOf(st)
	// 获取结构体值
	val := reflect.ValueOf(st)
	// 获取结构体类型
	kd := val.Kind()
	if kd != reflect.Struct {
		return errors.New("验证对象必须是结构体")
	}
	// 获取结构体字段数量
	num := typ.NumField()
	for i := 0; i < num; i++ {
		// 获取字段
		targetField := typ.Field(i)
		// 获取字段值
		targetVal := val.Field(i)
		// 字段值为结构体
		if targetVal.Kind() == reflect.Struct {
			if err = Verify(targetVal.Interface(), roleMap); err != nil {
				return err
			}
		}
		// 获取字段规则
		rules := roleMap[targetField.Name]
		if len(rules) > 0 {
			// 遍历规则
			for _, v := range rules {
				switch {
				// 不能为空
				case v == "notEmpty":
					if isBlank(targetVal) {
						return errors.New(targetField.Name + "值不能为空")
					}
				// 正则
				case strings.Split(v, "=")[0] == "regexp":
					if !regexpMatch(strings.Split(v, "=")[1], val.String()) {
						return errors.New(targetField.Name + "格式校验不通过")
					}
				// 比较
				case compareMap[strings.Split(v, "=")[0]]:
					if !compareVerify(val, v) {
						return errors.New(targetField.Name + "长度或值不在合法范围," + v)
					}
				}
			}
		}
	}
	return nil
}

// compareVerify
// @description: 比较验证
// @param: value reflect.Value
// @param: rule string
func compareVerify(val reflect.Value, rule string) bool {
	switch val.Kind() {
	case reflect.String:
		return compare(len([]rune(val.String())), rule)
	case reflect.Slice, reflect.Array:
		return compare(val.Len(), rule)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compare(val.Uint(), rule)
	case reflect.Float32, reflect.Float64:
		return compare(val.Float(), rule)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compare(val.Int(), rule)
	default:
		return false
	}
}

// compare
// @description: 比较
// @param: value interface{}
// @param: rule string
func compare(value interface{}, rule string) bool {
	arr := strings.Split(rule, "=")
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		VInt, VErr := strconv.ParseInt(arr[1], 10, 64)
		if VErr != nil {
			return false
		}
		switch arr[0] {
		case "eq":
			return val.Int() == VInt
		case "ne":
			return val.Int() != VInt
		case "lt":
			return val.Int() < VInt
		case "le":
			return val.Int() <= VInt
		case "gt":
			return val.Int() > VInt
		case "ge":
			return val.Int() >= VInt
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		VInt, VErr := strconv.Atoi(arr[1])
		if VErr != nil {
			return false
		}
		switch arr[0] {
		case "eq":
			return val.Uint() == uint64(VInt)
		case "ne":
			return val.Uint() != uint64(VInt)
		case "lt":
			return val.Uint() < uint64(VInt)
		case "le":
			return val.Uint() <= uint64(VInt)
		case "gt":
			return val.Uint() > uint64(VInt)
		case "ge":
			return val.Uint() >= uint64(VInt)
		}
	case reflect.Float32, reflect.Float64:
		VFloat, VErr := strconv.ParseFloat(arr[1], 64)
		if VErr != nil {
			return false
		}
		switch arr[0] {
		case "eq":
			return val.Float() == VFloat
		case "ne":
			return val.Float() != VFloat
		case "lt":
			return val.Float() < VFloat
		case "le":
			return val.Float() <= VFloat
		case "gt":
			return val.Float() > VFloat
		case "ge":
			return val.Float() >= VFloat
		}
	default:
		return false
	}
	return false
}

// isBlank
// @description: 判断是否为空
// @param: value reflect.Value
// @return: bool
func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String, reflect.Slice:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	default:
		panic("unhandled default case")
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

// regexpMatch
// @description: 正则匹配
// @param: rule string
// @param: matchStr string
// @return: bool
func regexpMatch(rule, matchStr string) bool {
	return regexp.MustCompile(rule).MatchString(matchStr)
}
