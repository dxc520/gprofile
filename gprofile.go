package gprofile

import (
	"encoding/json"
	"fmt"
	"github.com/olebedev/config"
	"io/ioutil"
	"reflect"
	"strconv"
)

func Profile(profile interface{}, configFile string, envHigher bool) (interface{}, error) {
	// 获取yaml配置文件并解析
	fmt.Printf("Application config file name: %s\n", configFile)
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	yamlString := string(file)
	cfg, err := config.ParseYaml(yamlString)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Application config yaml:\n %s\n", yamlString)
	// 设置环境变量、参数优先级
	if envHigher {
		cfg = cfg.Flag().Env()
	} else {
		cfg = cfg.Env().Flag()
	}
	// 获取生效的profile，默认dev
	activeProfile, err := cfg.String("profiles.active")
	if err != nil {
		cfg, err = cfg.Get(activeProfile)
		if err != nil {
			return nil, err
		}
	}
	return parseYaml(profile, cfg)
}

func parseYaml(profile interface{}, cfg *config.Config) (interface{}, error) {
	err := parseProfile(reflect.TypeOf(profile).Elem(), reflect.ValueOf(profile).Elem(), cfg, "")
	return profile, err
}

//递归遍历env
//如果属性为struct类型，递归遍历；否则执行assignment赋值
func parseProfile(t reflect.Type, v reflect.Value, cfg *config.Config, prefix string) error {
	for i := 0; i < t.NumField(); i++ {
		typeField := t.Field(i)
		valueField := v.Field(i)
		pv, pvExist := typeField.Tag.Lookup("profile")
		fieldLowerName := starterLower(typeField.Name)
		if !pvExist {
			pv = fieldLowerName
		} else if pv == "_" {
			continue
		}
		pv = prefix + pv
		if valueField.Kind() == reflect.Struct {
			err := parseProfile(typeField.Type, valueField, cfg, pv+".")
			if err != nil {
				return err
			}
		} else {
			err := assignment(&typeField.Tag, &valueField, cfg, pv)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func starterLower(s string) string {
	arr := []rune(s)
	if len(arr) > 0 && arr[0] > 64 && arr[0] < 91 {
		arr[0] += 32
	}
	return string(arr)
}

//根据属性类型赋值
//如果配置文件存在，则直接赋值；不存在则设置tag中的profileDefault默认值；默认值不存在则返回error
func assignment(tag *reflect.StructTag, v *reflect.Value, cfg *config.Config, pv string) error {
	switch v.Kind() {
	case reflect.Bool:
		cv, err := cfg.Bool(pv)
		if err == nil {
			v.SetBool(cv)
			return nil
		}
		dv, dvExist := tag.Lookup("profileDefault")
		if !dvExist {
			return fmt.Errorf("assignment error, no default value when set value: %s", pv)
		}
		n, err := strconv.ParseBool(dv)
		if err != nil {
			return fmt.Errorf("assignment error, when parse default value %s: %#v", pv, dv)
		}
		v.SetBool(n)
	case reflect.String:
		cv, err := cfg.String(pv)
		if err == nil {
			v.SetString(cv)
			return nil
		}
		dv, dvExist := tag.Lookup("profileDefault")
		if !dvExist {
			return fmt.Errorf("assignment error, no default value when set value: %s", pv)
		}
		v.SetString(dv)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		cv, err := cfg.Int(pv)
		if err == nil {
			v.SetUint(uint64(cv))
			return nil
		}
		dv, dvExist := tag.Lookup("profileDefault")
		if !dvExist {
			return fmt.Errorf("assignment error, no default value when set value: %s", pv)
		}
		n, err := strconv.ParseUint(dv, 10, 64)
		if err != nil || v.OverflowUint(n) {
			return fmt.Errorf("assignment error, when parse default value %s: %#v", pv, dv)
		}
		v.SetUint(n)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		cv, err := cfg.Int(pv)
		if err == nil {
			v.SetInt(int64(cv))
			return nil
		}
		dv, dvExist := tag.Lookup("profileDefault")
		if !dvExist {
			return fmt.Errorf("assignment error, no default value when set value: %s", pv)
		}
		n, err := strconv.ParseInt(dv, 10, 64)
		if err != nil || v.OverflowInt(n) {
			return fmt.Errorf("assignment error, when parse default value %s: %#v", pv, dv)
		}
		v.SetInt(n)
	case reflect.Float32, reflect.Float64:
		cv, err := cfg.Float64(pv)
		if err == nil {
			v.SetFloat(cv)
			return nil
		}
		dv, dvExist := tag.Lookup("profileDefault")
		if !dvExist {
			return fmt.Errorf("assignment error, no default value when set value: %s", pv)
		}
		n, err := strconv.ParseFloat(dv, v.Type().Bits())
		if err != nil || v.OverflowFloat(n) {
			return fmt.Errorf("assignment error, when parse default value %s: %#v", pv, dv)
		}
		v.SetFloat(n)
	case reflect.Slice:
		cv, err := cfg.List(pv)
		if err == nil {
			v.Set(reflect.ValueOf(cv))
			return nil
		}
		dv, dvExist := tag.Lookup("profileDefault")
		if !dvExist {
			return fmt.Errorf("assignment error, no default value when set value: %s", pv)
		}
		n := reflect.New(v.Type())
		err = json.Unmarshal([]byte(dv), n.Interface())
		if err != nil {
			return fmt.Errorf("assignment error, when parse default value %s: %#v", pv, dv)
		}
		v.Set(n.Elem())
	case reflect.Map:
		cv, err := cfg.Map(pv)
		if err == nil {
			v.Set(reflect.ValueOf(cv))
			return nil
		}
		dv, dvExist := tag.Lookup("profileDefault")
		if !dvExist {
			return fmt.Errorf("assignment error, no default value when set value: %s", pv)
		}
		n := reflect.New(v.Type())
		err = json.Unmarshal([]byte(dv), n.Interface())
		if err != nil {
			return fmt.Errorf("assignment error, when parse default value %s: %#v", pv, dv)
		}
		v.Set(n.Elem())
	default:
		return fmt.Errorf("assignment error, unsupport type: %#v", v)
	}
	return nil
}
