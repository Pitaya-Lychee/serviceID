package main

import (
	"encoding/xml"
)

type Application struct {
	Name      string     `xml:"name"`
	Instances []Instance `xml:"instance"`
}

type Instance struct {
	IP                 string   `xml:"ipAddr"`
	Port               string   `xml:"port"`
	Service_id_Polaris string   `xml:"app"`
	Metadata           Metadata `xml:"metadata"`
}

type Metadata struct {
	Service_id_cloud string `xml:"__serviceId"`
}

type Applications struct {
	XMLName      xml.Name      `xml:"applications"`
	Applications []Application `xml:"application"`
}
