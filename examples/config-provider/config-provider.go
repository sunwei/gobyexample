package main

import (
	"fmt"
	"reflect"
	"strings"
)

// Provider 定义提供方需要具备的能力
// 通过Key查询值
// 设置键值对
// 设置默认参数
type Provider interface {
	Get(key string) any
	Set(key string, value any)
	SetDefaults(params Params)
}

// Params 参数格式定义
// 关键字为字符类型
// 值为通用类型any
type Params map[string]any

// Set 根据新传入参数，对应层级进行重写
// pp为新传入参数
// p为当前参数
// 将pp的值按层级结构写入p
// 递归完成
func (p Params) Set(pp Params) {
	for k, v := range pp {
		vv, found := p[k]
		if !found {
			p[k] = v
		} else {
			switch vvv := vv.(type) {
			case Params:
				if pv, ok := v.(Params); ok {
					vvv.Set(pv)
				} else {
					p[k] = v
				}
			default:
				p[k] = v
			}
		}
	}
}

func New() Provider {
	return &defaultConfigProvider{
		root: make(Params),
	}
}

// defaultConfigProvider Provider接口实现对象
type defaultConfigProvider struct {
	root Params
}

// Get 按key获取值
// 约定""键对应的是c.root
// 嵌套获取值
func (c *defaultConfigProvider) Get(k string) any {
	if k == "" {
		return c.root
	}
	key, m := c.getNestedKeyAndMap(strings.ToLower(k))
	if m == nil {
		return nil
	}
	v := m[key]
	return v
}

// getNestedKeyAndMap 支持多级查询
// 通过分隔符"."获取查询路径
func (c *defaultConfigProvider) getNestedKeyAndMap(
	key string) (string, Params) {
	var parts []string
	parts = strings.Split(key, ".")
	current := c.root
	for i := 0; i < len(parts)-1; i++ {
		next, found := current[parts[i]]
		if !found {
			return "", nil
		}
		var ok bool
		current, ok = next.(Params)
		if !ok {
			return "", nil
		}
	}
	return parts[len(parts)-1], current
}

// Set 设置键值对
// 统一key的格式为小写字母
// 如果传入的值符合Params的要求，通过root进行设置
// 如果为非Params类型，则直接赋值
func (c *defaultConfigProvider) Set(k string, v any) {
	k = strings.ToLower(k)

	if p, ok := ToParamsAndPrepare(v); ok {
		// Set the values directly in root.
		c.root.Set(p)
	} else {
		c.root[k] = v
	}

	return
}

// SetDefaults will set values from params if not already set.
func (c *defaultConfigProvider) SetDefaults(
	params Params) {
	PrepareParams(params)
	for k, v := range params {
		if _, found := c.root[k]; !found {
			c.root[k] = v
		}
	}
}

// ToParamsAndPrepare converts in to Params and prepares it for use.
// If in is nil, an empty map is returned.
// See PrepareParams.
func ToParamsAndPrepare(in any) (Params, bool) {
	if IsNil(in) {
		return Params{}, true
	}
	m, err := ToStringMapE(in)
	if err != nil {
		return nil, false
	}
	PrepareParams(m)
	return m, true
}

// IsNil reports whether v is nil.
func IsNil(v any) bool {
	if v == nil {
		return true
	}

	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice:
		return value.IsNil()
	}

	return false
}

// ToStringMapE converts in to map[string]interface{}.
func ToStringMapE(in any) (map[string]any, error) {
	switch vv := in.(type) {
	case Params:
		return vv, nil
	case map[string]string:
		var m = map[string]any{}
		for k, v := range vv {
			m[k] = v
		}
		return m, nil

	default:
		fmt.Println("value type not supported yet")
		return nil, nil
	}
}

// PrepareParams
// * makes all the keys lower cased
// * This will modify the map given.
// * Any nested map[string]interface{}, map[string]string
// * will be converted to Params.
func PrepareParams(m Params) {
	for k, v := range m {
		var retyped bool
		lKey := strings.ToLower(k)

		switch vv := v.(type) {
		case map[string]any:
			var p Params = v.(map[string]any)
			v = p
			PrepareParams(p)
			retyped = true
		case map[string]string:
			p := make(Params)
			for k, v := range vv {
				p[k] = v
			}
			v = p
			PrepareParams(p)
			retyped = true
		}

		if retyped || k != lKey {
			delete(m, k)
			m[lKey] = v
		}
	}
}

func main() {
	// 新建Config Provider实例
	// 实例中defaultConfigProvider实现了接口
	provider := New()

	// 模拟设置用户自定义配置项
	// config.toml中关于主题的配置信息
	// 类型是map[string]string
	// 需要转换成map[string]any，也就是Params类型
	provider.Set("", map[string]string{
		"theme": "mytheme",
	})

	// 模拟默认配置项
	// 超时默认时间为30秒
	provider.SetDefaults(Params{
		"timeout": "30s",
	})

	// 输出提前设置的所有配置信息
	fmt.Printf("%#v\n", provider.Get(""))
	fmt.Printf("%#v\n", provider.Get("theme"))
	fmt.Printf("%#v\n", provider.Get("timeout"))
}
