# go-ansible-zabbix-inventoty
#СБОРКА
go env -w CGO_ENABLED=0
go build  -o zabbix-inventory main.go

#ПЕРЕМЕННЫЕ ОКРУЖЕНИЯ

Сначала пытается загрузить переменные из .env файла текущей директории  если файла не существует то c /etc/default/zabbix-inventory

INV_ZABBIX_URL="https://zabbix.domain.local"

INV_ZABBIX_USER="admin"

INV_ZABBIX_PASSWORD="Password"

INV_IGNORE_DISABLED_HOST="true" default false # не включает в inventory сервера которые имеют статус disable

INV_LIST_IPADDRESSES_FOR_INVENTORY="192.168.0.1,192.168.0.2" - добавляет список ip адресов в инвентори

INV_SSH_COMMON_ARGS_TAG="ansible_ssh_common_args"

INV_SSH_ANSIBLE_PORT_TAG="ansible_port"

INV_SSH_ANSIBLE_HOST_TAG = "ansible_host"

#INV_HOSTS_RESTRICT_BY_TAGS = "name:value" # можно просто по имени тега. Ограничивает список серверов inventory по тегам

INV_HOSTS_RESTRICT_BY_GROUPS = "OS/Linux" # Ограничивае список серверов в инвентори по по группам
