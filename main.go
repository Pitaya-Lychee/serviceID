package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	source80 := "polaris"
	token80 := "acc16c5dafdf295a7e396d3c99e770f8"
	url80_dev := "http://external-apiserver.cloud.ke.com"
	url80_test := "http://app-center-cloud-external-apiserver.ttb.test.ke.com"

	source_no80 := "cloud-registry-service"
	secret_no80 := "NmS6LAooGJS3CRfC"
	url_no80_dev := "http://api.cloud.intra.ke.com"
	url_no80_test := "http://dev.old-cloud.intra.ke.com"

	var op string
	var url80 string
	var url_no80 string
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("请加入dev或者test参数表示环境!")
		return
	} else if len(args) == 1 && args[0] == "dev" {
		fmt.Println("开始进行dev环境下的ServiceId对齐...")
		op = "dev"
		url80 = url80_dev
		url_no80 = url_no80_dev
	} else if len(args) == 1 && args[0] == "test" {
		fmt.Println("开始进行test环境下的ServiceId对齐...")
		op = "test"
		url80 = url80_test
		url_no80 = url_no80_test
	} else {
		fmt.Println("请加入dev或者test参数表示环境!")
		return
	}
	final_url80 := url80 + "/api/serviceid"
	final_url_no80 := url_no80 + "/api/application/get_application_by_ip_and_port"
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
						serviceIDCloud := ipQuery80(ip, final_url80, source80, token80)
						if serviceIDCloud != "" {
							RelationMapLock.Lock()
							RelationMap[instance.Service_id_Polaris] = serviceIDCloud
							RelationMapLock.Unlock()
						}
					} else {
						serviceIDCloud := ipQuery_no80(ip, port, final_url_no80, source_no80, secret_no80)
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
