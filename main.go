package main

import (
	"log"
	"zabbix-inventory/zabbix"
)

//'/api_jsonrpc.php'
func main() {
	api := zabbix.Zabbix{}
	err := api.Login()

	if err != nil {
		log.Fatal(err)
	}

	defer api.Logout()
	api.GetHosts()

	//fmt.Println(api)

}
