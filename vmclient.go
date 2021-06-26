package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Token struct {
	Token string `json:"token"`
}

var token string

func Login(url string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("username> ")
	username := read(reader)
	fmt.Print("password> ")
	password := read(reader)

	client := &http.Client{}
	var jsonStr = []byte(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", url), bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	str, code := getRespond(client, req)
	if code == http.StatusOK {
		var tokenObj Token
		err = json.Unmarshal([]byte(str), &tokenObj)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			token = tokenObj.Token
			fmt.Println("login succeeded")
		}
	} else {
		fmt.Println("response:")
		fmt.Println(str)
	}

}

func getRespond(client *http.Client, req *http.Request) (string, int) {
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println("Request id:", resp.Header.Get("requestID"))
	fmt.Println("Status code:", resp.StatusCode)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, bodyBytes, "", "\t")
	if err != nil {
		fmt.Println("JSON parse error: ", err)
	}

	return string(prettyJSON.Bytes()), resp.StatusCode
}

func DeleteVM(url string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("vm id> ")
	id := read(reader)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/vm/%s", url, id), nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+token)
	str, _ := getRespond(client, req)
	fmt.Println("response:")
	fmt.Println(str)
}

func GetVM(url string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("vm id> ")
	id := read(reader)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/vm/%s", url, id), nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+token)
	str, _ := getRespond(client, req)
	fmt.Println("response:")
	fmt.Println(str)
}

func UpdateVM(url string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("vm id> ")
	id := read(reader)

	fmt.Print("vm name> ")
	name := read(reader)

	fmt.Print("vm status> ")
	status := read(reader)

	client := &http.Client{}
	var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "status":"%s"}`, name, status))

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/vm/%s", url, id), bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	str, _ := getRespond(client, req)
	fmt.Println("response:")
	fmt.Println(str)
}

func GetAllVM(url string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/vm/all", url), nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+token)
	str, _ := getRespond(client, req)
	fmt.Println("response:")
	fmt.Println(str)
}

func RegisterVM(url string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("vm name> ")
	name := read(reader)

	fmt.Print("vm status> ")
	status := read(reader)

	client := &http.Client{}
	var jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "status":"%s"}`, name, status))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/vm/", url), bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	str, _ := getRespond(client, req)
	fmt.Println("response:")
	fmt.Println(str)
}

func read(reader *bufio.Reader) string {
	cmd, _ := reader.ReadString('\n')
	cmd = strings.Replace(cmd, "\n", "", -1)
	return cmd
}

func main() {
	var url string
	flag.StringVar(&url, "url", "http://localhost:8010", "REST webservice")
	flag.Parse()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("command: (login/list/get/register/delete/update/exit) > ")
		cmd := read(reader)
		if strings.Compare("login", cmd) == 0 {
			Login(url)
		} else if strings.Compare("list", cmd) == 0 {
			GetAllVM(url)
		} else if strings.Compare("get", cmd) == 0 {
			GetVM(url)
		} else if strings.Compare("register", cmd) == 0 {
			RegisterVM(url)
		} else if strings.Compare("delete", cmd) == 0 {
			DeleteVM(url)
		} else if strings.Compare("update", cmd) == 0 {
			UpdateVM(url)
		} else if strings.Compare("exit", cmd) == 0 {
			os.Exit(0)
		}
	}
}
