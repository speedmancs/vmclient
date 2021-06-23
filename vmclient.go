package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getRespond(client *http.Client, req *http.Request) {
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	str := string(bodyBytes)
	fmt.Println("response:", str)
}

func DeleteVM(url string, id string) {
	client := &http.Client{}
	fmt.Println(fmt.Sprintf("%s/vm/%s", url, id))
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/vm/%s", url, id), nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	getRespond(client, req)

}
func GetAllVM(url string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/vms", url), nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	getRespond(client, req)
}

func RegisterVM(url string, name string, status string) {
	client := &http.Client{}
	var jsonStr = []byte(
		fmt.Sprintf(`{"name":"%s", "status":"%s"}`, name, status))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/vm", url), bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	getRespond(client, req)
}

func main() {
	var url string
	flag.StringVar(&url, "url", "http://localhost:8010", "REST webservice")
	var cmd string
	flag.StringVar(&cmd, "cmd", "get", "command")
	var id string
	flag.StringVar(&id, "id", "", "vm id")
	var name string
	flag.StringVar(&name, "name", "", "vm name")
	var status string
	flag.StringVar(&status, "status", "", "vm status")

	flag.Parse()

	fmt.Println(cmd)
	if cmd == "getall" {
		GetAllVM(url)
	} else if cmd == "register" {
		RegisterVM(url, name, status)
	} else if cmd == "delete" {
		DeleteVM(url, id)
	}
}
