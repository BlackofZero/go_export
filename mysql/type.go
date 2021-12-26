package mysql

import (
	"database/sql"
	"fmt"
	"strings"
)



const (

	SHOW_SLOW_SQL= `SELECT 
    USER,
	HOST,
	DB,
	INFO,
    TIME
FROM
	information_schema.PROCESSLIST 
WHERE
	1 = 1 
	AND state LIKE '%Sending data%'
`
	SELECT_THEAD_ID=`SELECT
	GROUP_CONCAT( concat( 'kill ', id ) SEPARATOR ';' ) AS cmd 
FROM
	information_schema.PROCESSLIST 
WHERE
	1 = 1 
	AND state LIKE '%Sending data%'`

)
type SolwLog struct {
	USER string
	HOST string
	DB string
	INFO string
	TIME int64
}
type MysqlConnect struct {
	db * sql.DB
	connUrl string
	Msg MysqlMsg
}
type MysqlMsg struct {
	MYSQL_USER string
	MYSQL_PASSWD string
	MYSQL_HOST string
	MYSQL_PORT string
	MONITOR_USER string
}

func (mm *MysqlMsg) ConcatMsg() string{
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",mm.MYSQL_USER,mm.MYSQL_PASSWD,mm.MYSQL_HOST,mm.MYSQL_PORT,"information_schema")
}
func (mm *MysqlMsg) SpiltUser() string{
	users:="and user in ( ' "
	var strs string
	if len(mm.MONITOR_USER)>0 {

		ss := strings.Split(mm.MONITOR_USER, ",")
		for _,str:=range ss{
			strs=fmt.Sprint(strs,"','",str)
		}
	}
	return fmt.Sprint(users,strs,"')")
}


//CAST( sql_text AS CHAR ( 20000 ) CHARACTER SET utf8 )
type MError struct {
	ErrCode int
	ErrMsg string
}
