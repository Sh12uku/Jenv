package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var envPath = map[string]string{}
var exePath, _ = os.Executable()

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	readConf()
	switch os.Args[1] {
	case "list":
		fmt.Println("List of recorded path:")
		for jVersion, path := range envPath {
			if jVersion != getCurrentJEnv() {
				fmt.Println("  " + jVersion + "\t" + path)
			} else {
				fmt.Println("* " + jVersion + "\t" + path)
			}
		}
	case "add":
		envPath[os.Args[2]] = os.Args[3]
	case "del":
		if len(envPath) < 1 {
			fmt.Println("None conf")
		} else {
			delete(envPath, os.Args[2])
		}
	case "use":
		flag, err := setJEnv(os.Args[2])
		if flag {
			fmt.Println("Succeed!")
		} else {
			fmt.Println(err)
		}
	case "help":
		printHelp()
	default:
		fmt.Println("Wrong arguments!")
		printHelp()
	}
	// fmt.Println(os.Getenv("JAVA_HOME"))
	writeConf()
	// readConf()

	// writeConf()
}

// 获取当前环境变量，返回对应名称
func getCurrentJEnv() string {
	env := os.Getenv("JAVA_HOME")
	for jVersion, path := range envPath {
		if strings.Contains(path, env) {
			return jVersion
		}
	}
	return ""
}

// 设置环境变量
func setJEnv(jVersion string) (bool, error) {
	// err := os.Setenv("JAVA_HOME", envPath[jVersion])
	if len(envPath) < 1 {
		return false, errors.New("None conf")
	}
	cmd := exec.Command("setx", "JAVA_HOME", envPath[jVersion])
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	// if os.Getenv("JAVA_HOME") != envPath[jVersion] {
	// 	return false, errors.New("Fail")
	// }
	return true, nil
}

// 读取配置文件
func readConf() (bool, error) {
	_, err := os.Stat(filepath.Dir(exePath) + "\\jenv.conf") // 判断文件是否存在
	if err != nil {
		return false, err
	}
	conf, _ := os.OpenFile(filepath.Dir(exePath)+"\\jenv.conf", os.O_CREATE, 0666) // 无配置文件时创建一个
	var input = make([]byte, 1024)
	total, _ := conf.Read(input)
	err = json.Unmarshal(input[:total], &envPath) // 解析json
	if err != nil {
		fmt.Println("配置文件异常")
	}
	return true, nil
}

// 写入配置文件
func writeConf() bool {
	conf, err := os.OpenFile(filepath.Dir(exePath)+"\\jenv.conf", os.O_CREATE|os.O_TRUNC, 0666) // 打开时清空文件
	if err != nil {                                                                             // 有错误说明打开文件失败，一般为权限问题
		fmt.Println(err)
		return false
	}
	serStr, _ := json.Marshal(envPath)
	conf.Write(serStr)
	err = conf.Sync()
	defer conf.Close()
	if err != nil { // 有错误说明写入失败
		fmt.Println(err)
		return false
	}
	return true
}

func printHelp() {
	fmt.Println("Simple tool for switching JDK version")
	fmt.Println("usage: jenv.exe [action] [args...]\n")
	fmt.Println("example:")
	fmt.Println("jenv list                         List all local jdk versions")
	fmt.Println("jenv add jdk8 path/to/java        Add jdk path records")
	fmt.Println("jenv del jdk8                     Remove a record")
	fmt.Println("jenv use jdk8                     Switch to one of jdk version")
	fmt.Println("jenv help                         Print this message")
}
