package constantx

const (
	HttpService    = "http"
	GrpcService    = "grpc"
	GethService    = "geth"
	ApiService     = "api"
	ProcessService = "process"
	MachineService = "machine"

	MachineCpu     = "cpu"
	MachineMemory  = "memory"
	MachineDisk    = "disk"
	MachineProcess = "process"
	MachineNetwork = "network"
)

var (
	MachineAll       = []string{MachineCpu, MachineMemory, MachineDisk, MachineProcess, MachineNetwork}
	MachineNoNetwork = []string{MachineCpu, MachineMemory, MachineDisk, MachineProcess}
)
