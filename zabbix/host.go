package zabbix

type ZabbixHost struct {
	ID         string           `json:"hostid"`
	Host       string           `json:"host"`
	Name       string           `json:"name"`
	Status     string           `json:"status,omitempty"`
	Groups     []HostGroup      `json:"groups"`
	Interfaces []HostInterfaces `json:"interfaces"`
	Tags       []HostTag        `json:"tags"`
}

type HostGroup struct {
	ID   string `json:"groupid"`
	Name string `json:"name"`
	UUID string `json:"uuid,omitempty"`
	// Internal int
	// Flags    int
}
type HostInterfaces struct {
	ID string `json:"id"`
	IP string `json:"ip"`
}

type HostTag struct {
	Tag   string `json:"tag"`
	Value string `json:"value"`
}
