package zabbix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Zabbix struct {
	Url    string
	Token  string `json:"result"`
	Client http.Client
}

func (z *Zabbix) Login() error {
	godotenv.Load()

	z.Url = os.Getenv("ZABBIXURL")
	if z.Url == "" {
		return errors.New("environment variable ZABBIXURL not defined")
	}
	user := os.Getenv("ZABBIXUSER")
	if user == "" {
		return errors.New("environment variable ZABBIXUSER not defined")
	}
	passwd := os.Getenv("ZABBIXPASSWD")

	if passwd == "" {
		return errors.New("environment variable ZABBIXPASSWD not defined")
	}

	jsonAuth := fmt.Sprintf(`
	{
		"jsonrpc": "2.0",
		"method": "user.login",
		"params": {
			"user": "%s",
			"password": "%s"
		},
		"id": 1,
		"auth": null
	}`, user, passwd)

	data := bytes.NewBuffer([]byte(jsonAuth))

	req, err := http.NewRequest("POST", z.Url+"/api_jsonrpc.php", data)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := z.Client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &z)

	if err != nil {
		return err
	}

	return nil
}

func (z *Zabbix) Logout() error {
	jsonLogout := fmt.Sprintf(`
	{
		"jsonrpc": "2.0",
		"method": "user.logout",
		"params": [],
		"id": 1,
		"auth": "%s" 
	}`, z.Token)

	data := bytes.NewBuffer([]byte(jsonLogout))

	req, err := http.NewRequest("POST", z.Url+"/api_jsonrpc.php", data)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.Client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (z *Zabbix) GetHosts() error {
	jsonGetHosts := fmt.Sprintf(`
	{
		"jsonrpc": "2.0",
		"method": "host.get",
		"params": {
			"selectGroups": "extend",
			"filter": {
				"host": ["Zabbix server"]
				}
			},		
		"id": 1,
		"auth": "%s" 
	}`, z.Token)

	data := bytes.NewBuffer([]byte(jsonGetHosts))

	req, err := http.NewRequest("POST", z.Url+"/api_jsonrpc.php", data)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.Client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil

}
