package app

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type ServerConfig struct {
	Ip string	`ini:"ip"`
	Port int	`ini:"port"`
}

type MysqlConfig struct {
	Username string	`ini:"username"`
	Passwd string	`ini:"passwd"`
	Database string	`ini:"database"`
	Host string		`ini:"host"`
	Port int		`ini:"port"`
	Timeout float32	`ini:"timeout"`
}

type Config struct {
	Server ServerConfig	`ini:"server"`
	Mysql MysqlConfig	`ini:"mysql"`
}

func (c Config) Write(configName string) {
	f, err := os.Create(configName)
	if err != nil {
		fmt.Println("create file failed, err: ", err)
	}
	defer f.Close()


	_, _ = f.WriteString("#this is comment\n;this a comment\n;[]表示一个section\n")

	typeInfo := reflect.TypeOf(c)
	valueInfo := reflect.ValueOf(c)

	for i := 0; i < typeInfo.NumField(); i++ {
		typeSub := typeInfo.Field(i)
		typeSubType := typeSub.Type
		typeSubName := typeSub.Tag.Get("ini")

		_, _ = f.WriteString("["+typeSubName+"]\n")

		for j := 0; j < typeSubType.NumField(); j++ {
			typeSubMember := typeSubType.Field(j).Tag.Get("ini")
			valueSubMember := valueInfo.Field(i).Field(j)

			_, _ = f.WriteString(fmt.Sprintf("%v = %v", typeSubMember, valueSubMember))
			if j < typeSubType.NumField()-1 {
				_, _ = f.WriteString("\n")
			}

		}
		if i < typeInfo.NumField()-1 {
			_, _ = f.WriteString("\n\n")
		}
	}
}

func (c *Config) Read(configName string) {
	f, err := os.Open(configName)
	if err != nil {
		fmt.Println("open file failed, err: ", err)
	}
	defer f.Close()

	buff := make([]byte, 8)
	lines := make([]byte, 0)
	lineList := make([]string, 0)

	for {
		n, err := f.Read(buff)
		if err == io.EOF {
			break
		}
		for i := 0; i < n; i++ {
			lines = append(lines, buff[i])
		}

	}

	before := 0
	for i := 0; i < len(lines); i++ {
		if lines[i] == '\n' {
			line := fmt.Sprintf("%s", lines[before:i])
			if line != "" {
				lineList = append(lineList, line)
			}
			before = i + 1
		}
	}
	lineList = append(lineList, fmt.Sprintf("%s", lines[before:]))


	groupName := ""
	typeInfo := reflect.TypeOf(c)
	valueInfo := reflect.ValueOf(c)

	for _, line := range lineList {
		if len(line) < 1 {
			break
		}
		switch line[:1] {
		case ";", "#", "\n", "":
			break
		case "[":
			if strings.Contains(line, "]") {
				groupName = line[strings.Index(line,"[")+1:strings.Index(line, "]")]
				groupName = strings.ToUpper(groupName[:1]) + groupName[1:]
			}
		default:
			kvs := strings.Split(line, "=")
			key := strings.Trim(kvs[0]," ")

			key = strings.ToUpper(key[:1]) + key[1:]

			value := strings.Trim(kvs[1]," ")

			s, _ := typeInfo.Elem().FieldByName(groupName)
			t, _ := s.Type.FieldByName(key)

			v := valueInfo.Elem().FieldByName(groupName)
			switch t.Type.Kind() {
			case reflect.String:
				v.FieldByName(key).SetString(value)
			case reflect.Int:
				tempValue, _ := strconv.Atoi(value)
				v.FieldByName(key).SetInt(int64(tempValue))
			case reflect.Float32:
				tempValue, _ := strconv.ParseFloat(value,10)
				v.FieldByName(key).SetFloat(tempValue)
			default:
				fmt.Println("unknown kind")
			}
		}
	}
}

func TestReflect() {
	/*config := Config{
		ServerConfig{
			Ip: "192.168.35.129",
			Port: 1212,
		},
		MysqlConfig{
			Username: "root",
			Passwd: "redhat",
			Database: "zzz",
			Host: "localhost",
			Port: 3306,
			Timeout: 1.2,
		},
	}

	config.Write("config_1.ini")*/

	config1 := new(Config)

	config1.Read("config_3.ini")
	fmt.Println(config1)
	config1.Server.Ip = "localhost"
	config1.Write("config_4.ini")
}