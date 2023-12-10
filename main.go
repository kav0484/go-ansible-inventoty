package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"zabbix-inventory/inventory"
	"zabbix-inventory/zabbix"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

}

func main() {
	var hosts []zabbix.ZabbixHost

	//proxy_ip_tag := os.Getenv("PROXY_IP_TAG")
	//groupsWithJoinTag := os.Getenv("GROUPS_WITH_JOIN_TAG")

	zabbixUrl := os.Getenv("INV_ZABBIX_URL")
	if zabbixUrl == "" {
		log.Fatal("environment variable INV_ZABBIX_URL not defined")
	}

	zabbixUser := os.Getenv("INV_ZABBIX_USER")
	if zabbixUser == "" {
		log.Fatal("environment variable INV_ZABBIX_USER not defined")
	}

	zabbixPassword := os.Getenv("INV_ZABBIX_PASSWORD")
	if zabbixPassword == "" {
		log.Fatal(" variable INV_ZABBIX_PASSWORD not defined")
	}

	api := zabbix.Zabbix{}
	err := api.Login(zabbixUrl, zabbixUser, zabbixPassword)

	if err != nil {
		log.Fatal(err)
	}
	defer api.Logout()

	restrictGroups := os.Getenv("INV_HOSTS_RESTRICT_BY_GROUPS")
	restrictTags := os.Getenv("INV_HOSTS_RESTRICT_BY_TAGS")
	envIgnoreDisabledHost := os.Getenv("INV_IGNORE_DISABLED_HOST")

	ignoreDisabled, err := strconv.ParseBool(envIgnoreDisabledHost)

	if err != nil {
		ignoreDisabled = false
	}

	hosts, err = inventory.GetHostsFromZabbix(api, restrictGroups, restrictTags, ignoreDisabled)

	if err != nil {
		log.Fatal(err)
	}

	ansibleMeta := inventory.Inventory{}
	ansibleHostvars := map[string]interface{}{}
	hostvars := map[string]interface{}{}
	ansibleGroupVar := map[string][]string{}
	//groupvarTag := map[string][]string{}

	for _, host := range hosts {
		ansibleHost := map[string]string{}

		ansibleHost["visible_name"] = host.Name

		if len(host.Interfaces) == 0 {
			continue
		} else {
			ansibleHost["ansible_host"] = host.Interfaces[0].IP
		}

		hostvars[host.Name] = ansibleHost

		for _, g := range host.Groups {
			group := strings.Join(strings.Split(g.Name, " "), "_")
			ansibleGroupVar[group] = append(ansibleGroupVar[group], host.Name)
		}
	}

	//Hosts  search  with special tag for rewrite ip
	rewriteIPTag := os.Getenv("INV_HOST_REWRITE_IP_TAG")

	if rewriteIPTag != "" {
		rewriteIPHosts, err := inventory.GetHostWithTag(api, rewriteIPTag)

		if err != nil {
			log.Fatal(err)
		}

		if len(rewriteIPHosts) > 0 {
			for _, host := range rewriteIPHosts {
				ansibleHostValue, ok := hostvars[host.Name].(map[string]string)
				if ok {
					for _, tag := range host.Tags {
						if tag.Tag == rewriteIPTag {
							ansibleHostValue["ansible_host"] = tag.Value
						}
					}
				}
			}
		}
	}

	//Hosts search with special tag for add in inventory parameter ansible_ssh_common_args
	sshCommonArgsTag := os.Getenv("INV_SSH_COMMON_ARGS_TAG")

	if sshCommonArgsTag != "" {
		hostsWithCommonArgsTag, err := inventory.GetHostWithTag(api, sshCommonArgsTag)

		if err != nil {
			log.Fatal(err)
		}

		if len(hostsWithCommonArgsTag) > 0 {
			for _, host := range hostsWithCommonArgsTag {
				ansibleHostValue, ok := hostvars[host.Name].(map[string]string)
				if ok {
					for _, tag := range host.Tags {
						if tag.Tag == sshCommonArgsTag {
							ansibleHostValue["ansible_ssh_common_args"] = tag.Value
						}
					}
				}
			}
		}
	}

	//list ip addreses add to inventory for new servers. Inventory group is DEPLOY_NEW_SERVERS
	envListIP := os.Getenv("INV_LIST_IPADDRESSES_FOR_INVENTORY")

	if envListIP != "" {
		listIP := strings.Split(strings.TrimSpace(envListIP), ",")
		for _, ip := range listIP {
			ip = strings.TrimSpace(ip)
			hostvars[ip] = map[string]string{"ansible_host": ip}
			ansibleGroupVar["DEPLOY_NEW_SERVERS"] = append(ansibleGroupVar["DEPLOY_NEW_SERVERS"], ip)
		}
	}

	ansibleHostvars["hostvars"] = hostvars
	ansibleMeta.Meta = ansibleHostvars

	jsHostVars, _ := json.Marshal(ansibleMeta)
	jsGroup, _ := json.Marshal(ansibleGroupVar)
	//jsTag, _ := json.Marshal(groupvarTag)

	var mergeJSON map[string]interface{}

	json.Unmarshal(jsHostVars, &mergeJSON)
	json.Unmarshal(jsGroup, &mergeJSON)
	//json.Unmarshal(jsTag, &mergeJSON)

	resultMergeJson, _ := json.MarshalIndent(mergeJSON, "", "  ")

	fmt.Println(string(resultMergeJson))

}
