// Package DataBase 資料庫
package DataBase

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/go-xorm/core"

	"github.com/koebel217505/Project/projCommon/proType"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/koebel217505/Project/projCommon/projConvert"
)

var Engine *xorm.Engine
var MaxPlayerid uint32

//var databaselogger *log.Logger
var SQLDatabase string

// CheckIsTaiwanPhone 檢查是否台灣電話格式
func CheckIsTaiwanPhone(s string) bool {
	if (len(s) != 10) || ((string(s[0]) != "0") && (string(s[1]) != "9")) {
		return false
	}

	for key := range s {
		if _, e := strconv.Atoi(string(s[key])); e != nil {
			return false
		}
	}

	return true
}

////產log
//func init() {
//	os.Mkdir("log", os.ModePerm)
//
//	file, err := os.OpenFile("./log/SQL.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 666)
//	if err != nil {
//		databaselogger.Fatal(err)
//	}
//	databaselogger = log.New(file, "", log.LstdFlags)
//	databaselogger.SetFlags(log.LstdFlags | log.Lshortfile)
//	file.Close()
//}

// LoadDataBase 初始資料庫
func LoadDataBase(sqlDatabase, sqlIP, sqlPort, sqlAccount, sqlPassword string) error {
	if Engine != nil {
		if Engine.Ping() == nil {
			return nil
		} else {
			Engine.Close()
			Engine = nil
		}
	}

	sqlRoot := sqlAccount + ":" + sqlPassword + "@projTcp(" + sqlIP + ":" + sqlPort + ")/?charset=utf8"
	x, err := xorm.NewEngine("mysql", sqlRoot)
	if err != nil {
		return errors.New(fmt.Sprintf("Fail to sync database: %v\n", err))
	}

	_, err = x.Query(`CREATE DATABASE IF NOT EXISTS ` + sqlDatabase)
	if err != nil {
		return errors.New(fmt.Sprintf("Fail to sync database: %v\n", err))
	}
	defer x.Close()

	database := sqlAccount + `:` + sqlPassword + `@projTcp(` + sqlIP + ":" + sqlPort + `)/` + sqlDatabase + `?charset=utf8`
	Engine, err = xorm.NewEngine("mysql", database)
	if err != nil {
		return errors.New(fmt.Sprintf("Fail to sync database: %v\n", err))
	}

	Engine.SetLogLevel(core.LOG_UNKNOWN)
	//x.ShowSQL(true)
	//設置連接池的空閒數大小
	Engine.SetMaxIdleConns(10000)
	//設置最大打開連接數
	Engine.SetMaxOpenConns(10000)

	SQLDatabase = sqlDatabase

	CreateAccountTable()

	if result, e := Engine.QueryString(`Select Max(playerid) as playerid From mAccount`); (len(result) > 0) && (e == nil) {
		MaxPlayerid = projConvert.ConvUint32(result[0]["playerid"])

		if MaxPlayerid < 10000 {
			MaxPlayerid = 10000
		}
	} else {
		return errors.New(fmt.Sprintf("Fail to sync database: %v\n", err))
	}

	log.Println("SQL End")

	return nil
}

func Close() {
	if Engine != nil {
		Engine.Close()
		Engine = nil
	}
}

// CreateAccountTable 產Money帳號表格
func CreateAccountTable() {
	if result, _ := Engine.QueryString(`SHOW TABLES FROM ` + SQLDatabase + ` LIKE 'mAccount'`); len(result) == 0 {
		Engine.Exec(`ALTER TABLE mAccount AUTO_INCREMENT = 10001`)
	} else {

	}
	r, e := Engine.Exec(`CREATE TABLE mAccount (
   				playerid BIGINT(20) NOT NULL UNIQUE,
   				account1 VARCHAR(128) NOT NULL ,
   				password1 VARCHAR(128) NOT NULL,
   				account2 VARCHAR(128) NOT NULL ,
   				password2 VARCHAR(128) NOT NULL,
				account3 VARCHAR(128) NOT NULL ,
   				password3 VARCHAR(128) NOT NULL,
				account4 VARCHAR(128) NOT NULL ,
   				password4 VARCHAR(128) NOT NULL,
				account5 VARCHAR(128) NOT NULL ,
   				password5 VARCHAR(128) NOT NULL,
				status tinyInt(1) DEFAULT 0,
   				PRIMARY KEY (playerid),
				KEY (account1),
				KEY (account2),
				KEY (account3),
				KEY (account4),
				KEY (account5)
 				) ENGINE=INNODB DEFAULT CHARSET=utf8`)
	log.Println(r, e)
}

// DrapAccountTable 砍Money帳號表格
func DrapAccountTable() {
	r, e := Engine.Exec(`DROP TABLE mAccount`)
	log.Println(r, e)
}

// CheckIsChinaPhone 檢查是否大陸電話格式
func CheckIsChinaPhone(s string) bool {
	if (len(s) != 11) || (string(s[0]) != "1") {
		return false
	}

	for key := range s {
		if _, e := strconv.Atoi(string(s[key])); e != nil {
			return false
		}
	}

	return true
}

// CheckIsHongKongPhone 檢查是否香港電話格式
func CheckIsHongKongPhone(s string) bool {
	if len(s) != 8 {
		return false
	}

	for key := range s {
		if _, e := strconv.Atoi(string(s[key])); e != nil {
			return false
		}
	}

	return true
}

// GetByAccount 撈資料用帳號
func GetByAccount(data proType.AccountData) ([]map[string]string, error) {
	if Engine == nil {
		return nil, errors.New("money DataBase 未連線")
	}

	return Engine.QueryString(fmt.Sprintf("Select * From mAccount Where (account1='%s') OR (account2='%s') OR (account3='%s' OR (account4='%s') OR (account5='%s')) Limit 1", data.Account, data.Account, data.Account, data.Account, data.Account))
}

// AddAcount 加帳號
func AddAcount(k, playerId uint32, data proType.AccountData) ([]map[string]string, error) {
	if Engine == nil {
		return nil, errors.New("money DataBase 未連線")
	}

	if k == 1 {
		return Engine.QueryString(fmt.Sprintf("INSERT IGNORE mAccount(playerid, account1, password1) VALUES('%d', '%s', '%s')", playerId, strings.ToLower(data.Account), strings.ToLower(data.Password)))
	} else {
		setAcount(k, playerId, data)
		return nil, nil
	}

	return nil, nil
}

// setAcount 設定帳號
func setAcount(k, playerId uint32, data proType.AccountData) ([]map[string]string, error) {
	if Engine == nil {
		return nil, errors.New("money DataBase 未連線")
	}

	Engine.QueryString(fmt.Sprintf("UPDATE mAccount SET account%d='%s',password%d='%s' WHERE playerid='%d'  ", k, strings.ToLower(data.Account), k, strings.ToLower(data.Password), playerId))
	return nil, nil
}

// GetByPlayerID 撈資料用PlayerID
func GetByPlayerID(playerId uint32) ([]map[string]string, error) {
	if Engine == nil {
		return nil, errors.New("money DataBase 未連線")
	}

	return Engine.QueryString(fmt.Sprintf("Select * From mAccount Where (PlayerID='%d') Limit 1", playerId))
}

// AddPhoneAccount 加手機帳號
func AddPhoneAccount(language uint32, playerId uint32, data proType.AccountData) ([]map[string]string, error) {
	if language == 1 {
		if CheckIsTaiwanPhone(data.Account) == false && CheckIsHongKongPhone(data.Account) == false {
			return nil, errors.New("格式不對")
		}
	} else {
		if CheckIsChinaPhone(data.Account) == false {
			return nil, errors.New("格式不對")
		}
	}

	if result, e := GetByAccount(data); len(result) > 0 || e != nil {
		return result, errors.New("帳號已存在")
	}

	if _, e := AddAcount(2, playerId, data); e == nil {
		log.Println("account Error ", e)
		return nil, nil
	}

	return nil, errors.New("異常")
}

// ChangePhoneAccount 改手機帳號
func ChangePhoneAccount(language uint32, playerId uint32, data proType.AccountData) ([]map[string]string, error) {
	if language == 1 {
		if CheckIsTaiwanPhone(data.Account) == false && CheckIsHongKongPhone(data.Account) == false {
			return nil, errors.New("格式不對")
		}
	} else {
		if CheckIsChinaPhone(data.Account) == false {
			return nil, errors.New("格式不對")
		}
	}

	result, e := GetByAccount(data)
	if len(result) > 0 || e != nil {
		return result, errors.New("新帳號已存在")
	}

	result, e = GetByPlayerID(playerId)
	if len(result) == 0 || e != nil {
		return result, errors.New("舊帳號不存在")
	}

	if _, e := AddAcount(2, playerId, proType.AccountData{Account: data.Account, Password: result[0]["password2"]}); e == nil {
		log.Println("account Error ", e)
		return nil, nil
	}

	return nil, errors.New("異常")
}

// ChangePassword 改密碼
func ChangePassword(data proType.AccountData) ([]map[string]string, error) {
	result, e := GetByAccount(data)
	if e != nil {
		log.Println("account Error ", e)
	}

	if len(result) == 0 {
		return nil, nil
	} else {
		for i := 1; i <= 5; i++ {
			if strings.ToLower(data.Account) == strings.ToLower(result[0]["account"+projConvert.ConvString(i)]) {
				return Engine.QueryString(fmt.Sprintf("UPDATE mAccount SET password%d='%s' WHERE account%d='%s'", i, strings.ToLower(data.Password), i, strings.ToLower(data.Account)))
			}
		}
	}

	return nil, errors.New("異常")
}

// SetStaus 改狀態
func SetStaus(playerId uint32, status uint16) ([]map[string]string, error) {
	result, e := GetByPlayerID(playerId)
	if len(result) == 0 || e != nil {
		log.Println("playerId Error ", e)
	}

	return Engine.QueryString(fmt.Sprintf("UPDATE mAccount SET status=%d WHERE playerId=%d", status, playerId))
}

// CheckPassWord 檢查密碼
func CheckPassWord(data proType.AccountData, r []map[string]string) bool {
	for i := 1; i <= 5; i++ {
		if strings.ToLower(data.Account) == strings.ToLower(r[0]["account"+projConvert.ConvString(i)]) {
			if strings.ToLower(data.Password) == strings.ToLower(r[0]["password"+projConvert.ConvString(i)]) {
				return true
			} else {
				return false
			}
		}
	}

	return false
}

// GetPassWord 撈密碼
func GetPassWord(data proType.AccountData) string {
	result, e := GetByAccount(data)
	if e != nil {
		log.Println("account Error ", e)
	}

	if len(result) != 0 {
		for i := 1; i <= 5; i++ {
			if strings.ToLower(data.Account) == strings.ToLower(result[0]["account"+projConvert.ConvString(i)]) {
				return strings.ToLower(result[0]["password"+projConvert.ConvString(i)])
			}
		}
	}

	return ""
}

// AccountInserIfNotExist 檢查帳號，不存在就新增
func AccountInserIfNotExist(language uint32, data proType.AccountData) ([]map[string]string, uint32, uint32, uint32) {
	if result, e := GetByAccount(data); e != nil {
		log.Println("account Error ", e)

		return nil, 0, 1, 0
	} else {
		if len(result) == 0 {
			if language == 1 {
				if CheckIsTaiwanPhone(data.Account) && CheckIsHongKongPhone(data.Account) == false {
					return nil, 0, 2, 0
				}
			} else {
				if CheckIsChinaPhone(data.Account) {
					return nil, 0, 2, 0
				}
			}

			//if _, e := strconv.Atoi(string(data.Account[0])); e == nil {
			//	return nil, 0, 2, 0
			//}

			//playerid, _ := gonanoid.Generate("1234567890", 7)
			randomNum, _ := rand.Int(rand.Reader, big.NewInt(2))
			maxPlayerId := uint32(MaxPlayerid + 2 + uint32(randomNum.Uint64()))
			if result, e = AddAcount(1, maxPlayerId, data); e != nil {
				return nil, 0, 1, 0
			} else {
				MaxPlayerid = maxPlayerId
				return nil, 1, 0, maxPlayerId
			}
		} else {
			if CheckPassWord(data, result) == true {
				return result, 0, projConvert.ConvUint32(projConvert.ConvString(result[0]["status"])), projConvert.ConvUint32(projConvert.ConvString(result[0]["playerid"]))
			} else {
				return nil, 0, 1, 0
			}
		}
	}
	return nil, 0, 1, 0
}
