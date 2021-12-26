package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"prometheus-exporter/utils"
	"strconv"
)

func Init() *MysqlConnect {
	mc := &MysqlConnect{}
	mc.Msg.MONITOR_USER = utils.GetEnvWithDefault("MONITOR_USER","parking_saas")
	mc.Msg.MYSQL_PASSWD = utils.GetEnvWithDefault("MYSQL_PASSWD","123456")
	mc.Msg.MYSQL_HOST = utils.GetEnvWithDefault("MYSQL_HOST","127.0.0.1")
	mc.Msg.MYSQL_PORT = utils.GetEnvWithDefault("MYSQL_PORT","3306")
	mc.Msg.MYSQL_USER = utils.GetEnvWithDefault("MYSQL_USER","bd69673")
	mc.connUrl = mc.Msg.ConcatMsg()
	return mc
}

func (r *MysqlConnect) Connect(){
	var err error
	r.db, err = sql.Open("mysql", r.connUrl)
	defer func() {
		if err := recover(); err != nil {
			log.Println("[error]: ",err)
		}
	}()

	if err != nil{
		panic("mysql 连接异常")
	}
	r.db.SetMaxOpenConns(10)
	r.db.SetMaxIdleConns(5)
}

func (r *MysqlConnect) Close(){
	r.db.Close()
}
//func (r *MysqlConnect) GetPushUrl() string{
//	return r.Content.PushUrl
//}

func (r *MysqlConnect) RunDelSql(sqlStr string){
	_,error:=r.db.Exec(sqlStr)
	defer func() {
		if err := recover(); err != nil {
			log.Println("[error]: ",err)

		}
	}()
	if error != nil {
		panic(error)
	}

}
func (r *MysqlConnect) RunSql(sqlStr string)([] SolwLog, error){
	var Time ,user , info ,host,db string
	//now := time.Now()
	sqlStr = fmt.Sprint(sqlStr,r.Msg.SpiltUser() )

	rows,err :=r.db.Query(sqlStr)
	defer func() {
		if err := recover(); err != nil {
			log.Println("[error]: ",err)

		}
	}()
	if err != nil {
		panic(err)
	}
	contents := [] SolwLog{}
	for rows.Next() {
		err := rows.Scan( &user,&host,&db,&info,&Time)
		content:= SolwLog{}
		content.USER =user
		content.INFO =info
		content.DB =db
		content.TIME, _ = strconv.ParseInt(Time,10,64)
		content.HOST = host
		contents = append(contents, content)
		if err != nil {
			panic(err)
		}

	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return contents  ,err
}