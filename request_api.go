package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	// "time"
)

func ipQuery80(ip string, url string) string {
	iport := iportQuery80{
		PodIp: ip,
	}
	iportJson, err := json.Marshal(iport)
	if err != nil {
		fmt.Println("JSON 编码失败:", err)
		return ""
	}
	reqdata, err := http.NewRequest("GET", url, bytes.NewBuffer(iportJson))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return ""
	}
	reqdata.Header.Set("Content-Type", "application/json")

	// time.Sleep(1 * time.Second) // 控制请求速率为 1 秒钟发送一次
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(reqdata)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return ""
	}
	var resdata ServiceIdResponse80
	err = json.NewDecoder(resp.Body).Decode(&resdata)

	if err != nil {
		fmt.Println("解析 JSON 响应时发生错误:", err)
		return ""
	}
	if resdata.Ret {
		return resdata.Data.ServiceID
	}
	return ""

}

func ipQuery_no80(ip string, port string, url string) string {
	iport := iportQuery_no80{
		Ip:   ip,
		Port: port,
	}
	iportJson, err := json.Marshal(iport)
	if err != nil {
		fmt.Println("JSON 编码失败:", err)
		return ""
	}
	reqdata, err := http.NewRequest("GET", url, bytes.NewBuffer(iportJson))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return ""
	}

	reqdata.Header.Set("Content-Type", "application/json")

	// time.Sleep(1 * time.Second) // 控制请求速率为 1 秒钟发送一次
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(reqdata)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return ""
	}
	var resdata ServiceIdResponse
	err = json.NewDecoder(resp.Body).Decode(&resdata)
	if err != nil {
		fmt.Println("解析 JSON 响应时发生错误:", err)
		return ""
	}
	if resdata.Message == "ok" {
		return resdata.Data.AppID
	}
	return ""
}
