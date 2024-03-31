package consts

const (
	ProcfsBaseDir      string = "/proc/fs/lustre"
	SysfsBaseDir       string = "/sys/fs/lustre"
	KernelDebugBaseDir string = "/sys/kernel/debug/lustre"

	OST_REG string = `(?P<filesystem>.*)-(?P<ost>OST[0-9A-Fa-f]+)`
	MDT_REG string = `(?P<filesystem>.*)-(?P<mdt>MDT[0-9A-Fa-f]+)`
	QMT_REG string = `(?P<filesystem>.*)-(?P<qmt>QMT[0-9A-Fa-f]+)`
)

// Level represents the operation level of a collector
type Level int

const (
	Disabled Level = iota
	Basic
	Extended
)
