package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func NewDB() (*sql.DB, error) {
	dbUser := "lj"
	dbPass := "123456"
	dbName := "find_service"
	host := "127.0.0.1:3306"
	connstr := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, host, dbName)

	database, err := sql.Open("mysql", connstr)
	if err != nil {
		fmt.Println("数据库连接失败:", err)
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("数据库连接成功...")
	return database, nil
}

func insert_exec(Db *sql.DB, RelationMap map[string]string, op string) {
	var table_name string
	if op == "dev" {
		table_name = "service_id_dev"
	} else if op == "test" {
		table_name = "service_id_test"
	}
	var a int
	for service_id_polaris, service_id_cloud := range RelationMap {
		_, err := Db.Exec("insert into "+table_name+" (Service_id_polaris, Service_id_cloud) values (?, ?)", service_id_polaris, service_id_cloud)
		if err != nil {
			a++
			fmt.Println("数据库插入失败, ", err)
			fmt.Println("当前插入失败数: ", a)
			log.Fatal(err)
		}
	}
	fmt.Println("已成功插入数据到" + table_name + "表中...")
}

func delete_exec(Db *sql.DB, op string) {
	var table_name string
	if op == "dev" {
		table_name = "service_id_dev"
	} else if op == "test" {
		table_name = "service_id_test"
	}
	fmt.Println("清空原始数据库中...")
	_, err := Db.Exec("TRUNCATE TABLE " + table_name)
	if err != nil {
		fmt.Println("清空原始数据库时发生错误:", err)
		return
	}
	fmt.Println("原始数据库已清空,开始写入数据...")
}
