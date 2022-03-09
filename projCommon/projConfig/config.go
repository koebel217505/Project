package projConfig

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zieckey/goini"
)

type ConfigStruct struct {
	ProxyIP     string `json:"-"`
	ProxyPort   string `json:"-"`
	SQLAccount  string `json:"SQL帳號"`
	SQLPassword string `json:"SQL密碼"`
	SQLIP       string `json:"SQLIP"`
	SQLPort     string `json:"SQLPort"`
	SQLDatabase string `json:"SQLDatabase"`
}

func init() {
	loadConfig()
}

var Config ConfigStruct

// LoadIP bla-bla
func LoadIP(addr, content, subContent string) (IP string) {
	// currentDir, err := os.Getwd()
	ini := goini.New()
	err := ini.ParseFile(addr)
	if err != nil {
		// log.Println(currentDir + "/config/" + strings.ToLower(content) + "/config.ini" + "Error")
		return
	}
	IP, ok := ini.SectionGet(content, subContent)
	if !ok {
		log.Println(content, subContent, "ini.SectionGet Error:", ok)
		return
	}

	log.Println("LoadIP:", IP)
	return
}

// LoadIPMap bla-bla
func LoadIPMap(addr, content, subContent string) (IPMap map[uint16]string) {
	IPMap = make(map[uint16]string, 0)
	// currentDir, err := os.Getwd()
	ini := goini.New()
	err := ini.ParseFile(addr)
	if err != nil {
		// log.Println(currentDir + "/config/" + strings.ToLower(content) + "/config.ini" + "Error")
		return
	}
	v, ok := ini.GetKvmap(content)
	if !ok {
		log.Println("ini.GetKvmap Error:", ok)
		return
	}

	for key, value := range v {
		if strings.Contains(strings.ToLower(key), strings.ToLower(subContent)) == false {
			if strings.Contains(key, subContent) == false {
				continue
			}
			log.Println(subContent, "ToLower Error")
			continue
		}

		// projTcp, err := net.ResolveTCPAddr("projTcp", value)
		// if err != nil {
		// 	log.Printf("Coneect ResolveTCPAddr Error: %s\n", err.Error())
		// 	return
		// }
		// projTcp.IP.String(), projTcp.Port

		s := strings.Replace(key, subContent, "", -1)
		if s == "" {
			IPMap[0] = value
			continue
		}

		if ID, err := strconv.Atoi(s); err == nil {
			IPMap[uint16(ID)] = value
		}

	}

	log.Println(content, subContent, "IPMap:", IPMap)
	return
}

func loadConfig() {
	exePath, _ := os.Executable()
	path := filepath.Dir(exePath)
	filename := filepath.Join(path + "/config.ini")
	contents, _ := ioutil.ReadFile(filename)
	ini := goini.New()
	if err := ini.ParseFile(filename); err != nil {
		if err := ini.Parse(contents[3:], goini.DefaultLineSeparator, goini.DefaultKeyValueSeparator); err != nil {
			fmt.Printf("parse INI file %v failed : %v\n", filename, err.Error())
			return
		}
	}

	Config.ProxyIP, _ = ini.Get("ProxyIP")
	Config.ProxyPort, _ = ini.Get("ProxyPort")

	Config.SQLAccount, _ = ini.Get("SQLAccount")
	Config.SQLPassword, _ = ini.Get("SQLPassword")
	Config.SQLIP, _ = ini.Get("SQLIP")
	Config.SQLPort, _ = ini.Get("SQLPort")
	Config.SQLDatabase, _ = ini.Get("SQLDatabase")
}
