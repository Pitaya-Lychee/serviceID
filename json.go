package main

type iportQuery80 struct {
	PodIp string `json:"podIp"`
}

// type ServiceIdResponse80 struct {
// 	Data    *ResponseData80 `json:"data,omitempty"`
// 	Message *string         `json:"message,omitempty"` // 错误信息
// 	Ret     bool            `json:"ret"`               // 正确与否
// }

// type ResponseData80 struct {
// 	ServiceID string `json:"serviceId"`
// }

type ServiceIdResponse80 struct {
	Ret     bool            `json:"ret"`     // 正确与否
	Message string          `json:"message"` // 错误信息
	Data    *ResponseData80 `json:"data"`
}

type ResponseData80 struct {
	ServiceID string `json:"serviceId"`
}

type iportQuery_no80 struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

type ResponseData_no80 struct {
	AppID      string   `json:"app_id"`
	Owner      string   `json:"owner"`
	Remark     string   `json:"remark"`
	DomainList []string `json:"domain_list"`
	ServerList []string `json:"server_list"`
	PortList   []int    `json:"port_list"`
}

type ServiceIdResponse struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Data    *ResponseData_no80 `json:"data"`
}
