package models

import (
	"database/sql"
	"log"

	"github.com/qiucarson/blog/config"
)

var (
	db        *sql.DB
	db_prefix string
)

func init() {
	c := config.MysqlConfig()
	db, _ = sql.Open("mysql", c.Mysql)
	db_prefix = c.Prefix
}
func GetOption(name string) string {

	stmt, _ := db.Prepare("SELECT option_value FROM " + db_prefix + "options WHERE option_name=? limit 1")
	rows, err := stmt.Query(name)

	defer rows.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	var value string
	rows.Next()
	rows.Scan(&value)
	return value

}
