package transport

type Task struct {
	IP        string `json:"ip"`
	HostName  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"dta"`
}

// 接收任务
var (
	tasks = make(chan Task, 10)
)
