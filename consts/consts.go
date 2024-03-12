package consts

const (
	ProcfsBaseDir      string = "/proc/fs/lustre"
	SysfsBaseDir       string = "/sys/fs/lustre"
	KernelDebugBaseDir string = "/sys/kernel/debug/lustre"
)

// Level represents the operation level of a collector
type Level int

const (
	Basic Level = iota
	Extended
	Disabled
)
