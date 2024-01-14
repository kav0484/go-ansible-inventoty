package inventory

import (
	"encoding/json"
	"fmt"
	"strings"
	"zabbix-inventory/zabbix"
)

type Inventory struct {
	Meta Host `json:"_meta"`
}

type Host map[string]interface{}

func CreateGroupWithJoinTag() {
	fmt.Println("ok")
}

// ProcessZabbixTags Возвращается список тегов и значение по которым нужно ограничить выборку хостов в zabbix
func ProcessZabbixTags(restrictTags string) []map[string]string {
	var tags []map[string]string

	listTags := strings.Split(restrictTags, ",")

	for _, t := range listTags {
		tag := make(map[string]string)
		splitTag := strings.Split(t, ":")
		tag["tag"] = splitTag[0]
		if len(splitTag) == 2 {
			tag["value"] = splitTag[1]
		}

		tags = append(tags, tag)
	}

	return tags
}

// GetHostsRestrictByTags function which get hosts from zabbix  restricting request get hosts  by  tags
func GetHostsRestrictByTags(api zabbix.Zabbix, zabbixTags []map[string]string) ([]zabbix.ZabbixHost, error) {
	params := make(map[string]interface{})
	params["output"] = "extend"
	params["selectInterfaces"] = []string{"ip"}
	params["selectTags"] = "extend"
	params["tags"] = zabbixTags
	jsonParams, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	hosts, err := api.GetHosts(string(jsonParams))

	if err != nil {
		return nil, err
	}

	return hosts, nil
}

// GetHostsFromZabbix is function getting hosts from zabbix server
func GetHostsFromZabbix(api zabbix.Zabbix, zabbixGroups string, zabbixTags string, ignoreDisabled bool) ([]zabbix.ZabbixHost, error) {

	var groups []string
	var groupids []string
	var hosts []zabbix.ZabbixHost

	if zabbixGroups != "" {
		groups = strings.Split(zabbixGroups, ",")
		for _, group := range groups {
			zabbixGroup, err := api.GetGroup(group)

			if err != nil {
				return nil, err
			}
			groupids = append(groupids, zabbixGroup.ID)
		}
	}

	params := make(map[string]interface{})
	params["output"] = []string{"hostid", "name"}
	params["selectInterfaces"] = []string{"ip"}
	params["selectGroups"] = []string{"groupid", "name"}
	//params["filter"] = map[string][]string{"host": {"4sl-analytic", "4alex-testansible", "7mmtdmz-gdzp01"}}
	if ignoreDisabled {
		params["filter"] = map[string][]string{"status": {"0"}}
	}

	if groupids != nil {
		params["groupids"] = groupids
	}
	if zabbixTags != "" {
		params["tags"] = ProcessZabbixTags(zabbixTags)
	}

	jsonParams, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	hosts, err = api.GetHosts(string(jsonParams))

	if err != nil {
		return nil, err
	}

	return hosts, nil
}

// GetHostsWithGroups is function getting hosts from zabbix server with groups
func GetHostsWithGroups(api zabbix.Zabbix) ([]zabbix.ZabbixHost, error) {
	params := make(map[string]interface{})
	params["output"] = []string{"hostid", "name"}
	params["selectInterfaces"] = []string{"ip"}
	params["selectGroups"] = []string{"groupid", "name"}
	//params["filter"] = map[string][]string{"host": {"v2202311180000247161"}}
	jsonParams, err := json.Marshal(params)
	hosts, err := api.GetHosts(string(jsonParams))
	if err != nil {
		return nil, err
	}

	return hosts, nil

}

func GetHostWithTag(api zabbix.Zabbix, tag string) ([]zabbix.ZabbixHost, error) {
	params := make(map[string]interface{})
	params["output"] = []string{"hostid", "name"}
	params["selectInterfaces"] = []string{"ip"}
	params["selectTags"] = "extend"
	params["tags"] = ProcessZabbixTags(tag)
	jsonParams, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	hosts, err := api.GetHosts(string(jsonParams))

	if err != nil {
		return nil, err
	}

	return hosts, nil
}
