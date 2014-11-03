package executor

import "github.com/cloudfoundry-incubator/runtime-schema/models"

type State string

const (
	StateInvalid      State = ""
	StateReserved     State = "reserved"
	StateInitializing State = "initializing"
	StateCreated      State = "created"
	StateCompleted    State = "completed"
)

type Container struct {
	Guid string `json:"guid"`

	State State `json:"state"`

	MemoryMB  int  `json:"memory_mb"`
	DiskMB    int  `json:"disk_mb"`
	CPUWeight uint `json:"cpu_weight"`

	Tags Tags `json:"tags,omitempty"`

	AllocatedAt int64 `json:"allocated_at"`

	RootFSPath string        `json:"root_fs"`
	Ports      []PortMapping `json:"ports"`
	Log        LogConfig     `json:"log"`

	Actions []models.ExecutorAction `json:"actions"`
	Env     []EnvironmentVariable   `json:"env,omitempty"`

	RunResult ContainerRunResult `json:"run_result"`
}

type EnvironmentVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type LogConfig struct {
	Guid       string `json:"guid"`
	SourceName string `json:"source_name"`
	Index      *int   `json:"index"`
}

type PortMapping struct {
	ContainerPort uint32 `json:"container_port"`
	HostPort      uint32 `json:"host_port,omitempty"`
}

type ContainerRunResult struct {
	Guid string `json:"guid"`

	Failed        bool   `json:"failed"`
	FailureReason string `json:"failure_reason"`
}

type ExecutorResources struct {
	MemoryMB   int `json:"memory_mb"`
	DiskMB     int `json:"disk_mb"`
	Containers int `json:"containers"`
}

type Tags map[string]string
