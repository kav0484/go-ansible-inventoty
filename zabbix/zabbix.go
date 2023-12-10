package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Zabbix struct {
	Url    string
	Token  string `json:"result"`
	Client http.Client
}

func (z *Zabbix) Login(zabbixUrl string, zabbixUser string, zabbixPassword string) error {
	z.Url = zabbixUrl

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
	}`, zabbixUser, zabbixPassword)

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
	body, err := io.ReadAll(resp.Body)
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

func (z *Zabbix) GetHosts(params string) ([]ZabbixHost, error) {
	jsonGetHosts := fmt.Sprintf(`
	{
		"jsonrpc": "2.0",
		"method": "host.get",
		"params": %s,
		"id": 1,
		"auth": "%s" 
	}`, params, z.Token)

	data := bytes.NewBuffer([]byte(jsonGetHosts))

	req, err := http.NewRequest("POST", z.Url+"/api_jsonrpc.php", data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.Client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	//fmt.Println(string(body))

	if err != nil {
		return nil, err
	}

	type response struct {
		Jsonrpc string       `json:"jsonrpc"`
		Result  []ZabbixHost `json:"result"`
	}

	var res *response

	err = json.Unmarshal(body, &res)

	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (z *Zabbix) GetGroup(groupName string) (HostGroup, error) {

	listGroup := []string{groupName}

	convertString, _ := json.Marshal(listGroup)

	jsonGetGroup := fmt.Sprintf(`
	{
		"jsonrpc": "2.0",
		"method": "hostgroup.get",
		"params": {
			"output": ["groupid","name"],
			"filter": {
				"name": %v
			}
		},		
		"id": 1,
		"auth": "%s" 
	}`, string(convertString), z.Token)

	data := bytes.NewBuffer([]byte(jsonGetGroup))

	req, err := http.NewRequest("POST", z.Url+"/api_jsonrpc.php", data)
	if err != nil {
		return HostGroup{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.Client.Do(req)

	if err != nil {
		return HostGroup{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return HostGroup{}, err
	}

	type response struct {
		Jsonrpc string      `json:"jsonrpc"`
		Result  []HostGroup `json:"result"`
	}

	var res *response

	err = json.Unmarshal(body, &res)

	if err != nil {
		return HostGroup{}, err
	}

	return res.Result[0], nil

}
