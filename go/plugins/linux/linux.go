package linux

import (
	"fmt"
	"motadata-lite/client/SSHclient"
	"motadata-lite/utils"
	"motadata-lite/utils/constants"
	"strings"
)

const (
	VENDOR = "system.vendor"

	SYSTEM_NAME = "system.os.name"

	SYSTEM_VERSION = "system.os.version"

	START_TIME = "started.time"

	START_TIME_SECOND = "started.time.seconds"

	SYSTEM_MODEL = "system.model"

	SYSTEM_PRODUCT = "system.product"

	INTERRUPT_PER_SECONDS = "system.interrupts.per.sec"

	SYSTEM_CPU_IO_PERCENT = "system.cpu.io.percent"

	SYSTEM_RUNNING_PROCESSES = "system.running.processes"

	SYSTEM_NETWORK_UDP_CONNECTIONS = "system.network.udp.connections"

	SYSTEM_NETWORK_TCP_CONNECTIONS = "system.network.tcp.connections"

	SYSTEM_NETWORK_TCP_RETRANSMISSIONS = "system.network.tcp.retransmissions"

	SYSTEM_NETWORK_ERROR_PACKETS = "system.network.error.packets"

	SYSTEM_NETWORK_OUT_BYTES_RATE = "system.network.out.bytes.rate"

	SYSTEM_MEMORY_TOTAL_BYTES = "system.memory.total.bytes"

	SYSTEM_MEMORY_AVAILABLE_BYTES = "system.memory.available.bytes"

	SYSTEM_CACHE_MEMORY_BYTES = "system.cache.memory.bytes"

	SYSTEM_SWAP_MEMORY_PROVISIONED = "system.swap.memory.provisioned.bytes"

	SYSTEM_SWAP_MEMORY_USED = "system.swap.memory.used.bytes"

	SYSTEM_SWAP_MEMORY_USED_PERCENT = "system.swap.memory.used.percent"

	SYSTEM_SWAP_MEMORY_FREE_PERCENT = "system.swap.memory.free.percent"

	SYSTEM_SWAP_MEMORY_FREE_BYTES = "system.swap.memory.free.bytes"

	SYSTEM_BUFFER_MEMORY_BYTES = "system.buffer.memory.bytes"

	SYSTEM_MEMORY_USED_BYTES = "system.memory.used.bytes"

	SYSTEM_MEMORY_FREE_BYTES = "system.memory.free.bytes"

	SYSTEM_MEMORY_FREE_PERCENT = "system.memory.free.percent"

	SYSTEM_MEMORY_USED_PERCENT = "system.memory.used.percent"

	SYSTEM_OVERALL_MEMORY_FREE_PERCENT = "system.overall.memory.free.percent"

	SYSTEM_OPENED_FILE_DESCRIPTORS = "system.opened.file.descriptors"

	SYSTEM_DISK_CAPACITY_BYTES = "system.disk.capacity.bytes"

	SYSTEM_DISK_FREE_BYTES = "system.disk.free.bytes"

	SYSTEM_DISK_FREE_PERCENT = "system.disk.free.percent"

	SYSTEM_DISK_USED_PERCENT = "system.disk.used.percent"

	SYSTEM_DISK_USED_BYTES = "system.disk.used.bytes"

	SYSTEM_DISK_IO_TIME_PERCENT = "system.disk.io.time.percent"

	SYSTEM_LOAD_AVG1_MIN = "system.load.avg1.min"

	SYSTEM_LOAD_AVG5_MIN = "system.load.avg5.min"

	SYSTEM_LOAD_AVG15_MIN = "system.load.avg15.min"

	SYSTEM_CPU_INTERRUPT_PERCENT = "system.cpu.interrupt.percent"

	SYSTEM_CPU_USER_PERCENT = "system.cpu.user.percent"

	SYSTEM_CPU_PERCENT = "system.cpu.percent"

	SYSTEM_CPU_KERNEL_PERCENT = "system.cpu.kernel.percent"

	SYSTEM_CPU_IDLE_PERCENT = "system.cpu.idle.percent"

	SYSTEM_CPU_TYPE = "system.cpu.type"

	SYSTEM_CPU_CORE = "system.cpu.core"

	SYSTEM_CONTEXT_SWITCHES_PER_SEC = "system.context.switches.per.sec"
)

var logger = utils.NewLogger("goEngine/plugin", "plugin")

func Discovery(jsonInput map[string]interface{}, errContext *[]map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {

			logger.Error(fmt.Sprintf("%v", err))

			*errContext = append(*errContext, map[string]interface{}{
				constants.ErrorCode:    21,
				constants.ErrorMessage: "formating problem",
				constants.Error:        err,
			})
		}
	}()

	client := SSHclient.Client{}

	for _, credential := range jsonInput[constants.Credentials].([]interface{}) {

		client.SetContext(jsonInput, credential.(map[string]interface{}))

		isValid, _ := client.Init()

		if isValid {
			jsonInput[constants.Result] = map[string]interface{}{constants.ObjectIP: jsonInput[constants.ObjectIP].(string)}

			jsonInput["credential.profile.id"] = credential.(map[string]interface{})["credential.id"].(float64)

			return
		}
	}

	jsonInput["credential.profile.id"] = constants.InvalidCredentialCode

	logger.Trace("returning to bootstrap")
}

func Collect(jsonInput map[string]interface{}, errContext *[]map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {

			logger.Error(fmt.Sprintf("%v", err))

			*errContext = append(*errContext, map[string]interface{}{
				constants.ErrorCode:    21,
				constants.ErrorMessage: "formating problem",
				constants.Error:        err,
			})
		}
	}()

	var err error

	client := SSHclient.Client{}

	client.SetContext(jsonInput, jsonInput[constants.Credential].(map[string]interface{}))

	isValid, err := client.Init()

	if !isValid {
		*errContext = append(*errContext, map[string]interface{}{
			constants.ErrorCode:    12,
			constants.ErrorMessage: "Cannot establish connection to host",
			constants.Error:        err.Error(),
		})
		return
	}

	var command = "cat /sys/class/dmi/id/sys_vendor; uname -sr | awk {'print $1,$2'};uptime -s; date -d \"$(uptime -s)\" +%s; cat /sys/class/dmi/id/product_name;hostname;ps aux | awk {'print $8'} | grep -wc 'X';ps -eLf | wc -l;ps -e --no-headers | wc -l ;netstat -un | wc -l | awk '{print $1 - 2}';netstat -tn | wc -l | awk '{print $1 - 2}';cat /proc/net/snmp | grep -i 'TCP' | tail -n 1 | awk {'print $13'};ifconfig | awk '/errors/ {sum += $3} END {print sum}';ifconfig | awk '/TX packets/ {sum += $3} END {print sum}';top -bn1 | head -n 4 | tail -n 1 | awk {'print $4,$6,$10'};top -bn1 | head -n 5 | tail -n 1 | awk {'print $3,$5,$7'};free | head -n 2 | tail -n 1 | awk {'print $6'};free -b | head -n 2 | tail -n 1 | awk {'print $3,$4'};free | awk '/Mem:/ {printf(\"%.2f\\n\", ($3/$2) * 100)}';free | awk '/Mem:/ {printf(\"%.2f\\n\", ($4/($3+$4)) * 100)}';lsof | wc -l | awk '{print $1 - 2}';df --total | tail -n 1 | awk '{print $2,$4, ($4/$2)*100 \"%\"}';df --output=pcent | head -3 | tail -1 | awk {'print $NF'}; df --output=used / | awk 'NR==2 {print $1}';iostat |head -n 4 | tail -n -1 | awk '{print $4}';top -bn1 | head -n 1 | awk {'print $10,$11,$NF'};mpstat -I ALL | head -n 4 | tail -n 1 | awk {'print $NF'};top -bn1 | head -n 3 | tail -n 1 | awk {'print $2,$4'};mpstat | tail -n -1 | awk {'print $8'};mpstat |tail -n 1 |  awk {'print $7,$NF'};lscpu | head -n 1 | awk {'print $NF'};lscpu | head -n 5 | tail -n 1 | awk {'print $NF'};vmstat | tail -n 1 | awk {'print $12'}"

	queryOutput, err := client.ExecuteCommand(command)

	if err != nil {

		*errContext = append(*errContext, map[string]interface{}{

			constants.ErrorCode: 11,

			constants.ErrorMessage: "error in the command",

			constants.Error: err.Error(),
		})

		logger.Error(fmt.Sprintf("error in the command: %s", err.Error()))

		return

	}

	defer func(client *SSHclient.Client) {
		err := client.Close()

		if err != nil {
			logger.Error(fmt.Sprintf("error in closing ssh connection: %s", err.Error()))
		}
	}(&client)

	lines := strings.Split(string(queryOutput), "\n")

	var output = make(map[string]interface{})

	defer func() {
		if r := recover(); r != nil {
			*errContext = append(*errContext, map[string]interface{}{
				constants.ErrorCode:    16,
				constants.ErrorMessage: "error in the reading output lines",
				constants.Error:        "out of index",
			})

			logger.Error(fmt.Sprintf("error in the reading output lines: %s", err.Error()))
		}
		return
	}()

	ParseMetric(output, lines)

	jsonInput[constants.Result] = output

}
