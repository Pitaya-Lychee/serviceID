package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func ipQuery80(ip string, url string, source string, token string) string {
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
	reqdata.Header.Set("source", source)
	reqdata.Header.Set("token", token)

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
		fmt.Println("解析 JSON 响应时发生错误1:", err)
		return ""
	}
	if resdata.Ret {
		return resdata.Data.ServiceID
	}
	return ""
}

func ipQuery_no80(ip string, port string, url string, source string, secret string) string {
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
	currentTime := time.Now().Unix()
	time_str := strconv.FormatInt(currentTime, 10)
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%s%d%s", source, currentTime, secret)))
	signature := hex.EncodeToString(hash.Sum(nil))
	reqdata.Header.Set("Content-Type", "application/json")
	reqdata.Header.Set("timestamp", time_str)
	reqdata.Header.Set("source", source)
	reqdata.Header.Set("signature", signature)
	// time.Sleep(1 * time.Second) // 控制请求速率为 1 秒钟发送一次
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(reqdata)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return ""
	}

	var resdata map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&resdata)
	if err != nil {
		fmt.Println("解析 JSON 响应时发生错误2:", err)
		return ""
	}
	fmt.Println("返回的json数据:", resdata)
	dataValue, ok := resdata["data"]
	if !ok || dataValue == nil {
		return ""
	}
	dataMap, ok := dataValue.(map[string]interface{})
	if !ok {
		return ""
	}
	appID, ok := dataMap["app_id"]
	if !ok {
		return ""
	}
	if appIDString, ok := appID.(string); ok {
		return appIDString
	}
	return ""
	// var resdata AllResponse_no80
	// err = json.NewDecoder(resp.Body).Decode(&resdata)
	// if err != nil {
	// 	fmt.Println("解析 JSON 响应时发生错误2:", err)
	// 	return ""
	// }
	// fmt.Println(resdata.Data.AppID)
	// fmt.Println(resdata.Code)
	// fmt.Println(resdata.Message)
}
