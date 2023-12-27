package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	source := "polaris"
	token := "acc16c5dafdf295a7e396d3c99e770f8"
	var secret string
	var op string
	var url string
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("请加入dev或者test参数表示环境!")
		return
	} else if len(args) == 1 && args[0] == "dev" {
		fmt.Println("开始进行dev环境下的ServiceId对齐...")
		op = "dev"
		url = "http://api.cloud.intra.ke.com"
	} else if len(args) == 1 && args[0] == "test" {
		fmt.Println("开始进行test环境下的ServiceId对齐...")
		op = "test"
		url = "http://dev.old-cloud.intra.ke.com"
	} else {
		fmt.Println("请加入dev或者test参数表示环境!")
		return
	}
	url80 := url + "/api/serviceid"
	url_no80 := url + "/api/application/get_application_by_ip_and_port"
	RelationMap := make(map[string]string)
	var cfgFile string = "./hello.xml"
	xmlFile, err := os.Open(cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	var data Applications

	// 解析XML数据
	fmt.Println("解析XML文件...")
	err = xml.NewDecoder(xmlFile).Decode(&data)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("XML文件解析完成")
	}

	var wg sync.WaitGroup
	maxWorkers := 20
	workCh := make(chan Application, len(data.Applications))
	var RelationMapLock sync.Mutex

	fmt.Println("正在遍历XML,并请求查询服务id中...")

	// 遍历Applications中的Application
	for i := 0; i < maxWorkers; i++ {
		go func() {
			for app := range workCh {
				for _, instance := range app.Instances {
					ip := instance.IP
					port := instance.Port
					if port == "8080" {
						serviceIDCloud := ipQuery80(ip, url80, source, secret, token)
						if serviceIDCloud != "" {
							RelationMapLock.Lock()
							RelationMap[instance.Service_id_Polaris] = serviceIDCloud
							RelationMapLock.Unlock()
						}
					} else {
						serviceIDCloud := ipQuery_no80(ip, port, url_no80, source, secret, token)
						if serviceIDCloud != "" {
							RelationMapLock.Lock()
							RelationMap[instance.Service_id_Polaris] = serviceIDCloud
							RelationMapLock.Unlock()
						}
					}
				}
				wg.Done()
			}
		}()
	}
	for _, app := range data.Applications {
		workCh <- app
		wg.Add(1)
	}

	wg.Wait()

	fmt.Println("遍历完成,已查询到对应service_id,请求数据库中...")
	db, err := NewDB()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()
	delete_exec(db, op)
	insert_exec(db, RelationMap, op)
}
