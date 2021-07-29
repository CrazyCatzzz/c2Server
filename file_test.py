from b64 import *
import time
from typing import Tuple
from secrets import token_bytes
from ctypes import cdll, c_char_p


def read_file(file_path):
    chunk = 16 * 1024
    # 打开文件
    with open(file_path, "rb") as f:
        stream = f.read()
        f.close()
    # 对流进行encrypt
    print("*************** xor bytes *************")
    stream = encode(stream)
    print(stream)
    # print("YVNCMjEyMzEyMyYqJqrvv6Cmv6mjqKrp")
    # time1 = time.time()
    # print("第1次base64 ", time1)
    # stream = stream.encode()
    # stream = xor_decrypt(stream).encode()
    # time2 = time.time()
    # print("xor ", time2)
    # stream1 = encode(stream)
    # time3 = time.time()
    # print("第2次base64 ", time3)
    # print("xor 所花费时间", time2-time1, "base64所花费时间", time3-time2)
    # print("len of stream", len(stream1))
    # print("*************** xor bytes *************")
    # print("*************** xor str *************")
    # with open(file_path, "rb") as f:
    #     stream = f.read()
    #     f.close()
    # stream = encode(stream)
    # time1 = time.time()
    # print("第1次base64 ", time1)
    # stream = xor_crypt(stream).encode()
    # time2 = time.time()
    # print("xor ", time2)
    # stream2 = encode(stream)
    # time3 = time.time()
    # print("第2次base64 ", time3)
    # print("xor 所花费时间", time2 - time1, "base64所花费时间", time3 - time2)
    # print("len of stream", len(stream2))
    # print("[-100:]", stream2[-100:])
    # print("*************** xor str *************")
    # 返回加密后的流
    # if stream1 == stream2:
    #     print("yes")
    # return stream1


def write_file():
    # 对流进行解密

    #
    pass


read_file("C:\\PythonWorkSpace\\tornado_server\\qqqq.txt")

# file_xor = cdll.LoadLibrary("./file_xor.so")
# file_xor.encrypt.argtypes = [c_char_p]
# file_xor.encrypt.restype = c_char_p
# file_xor.decrypt.argtypes = [c_char_p]
# file_xor.decrypt.restype = c_char_p
#
#
# def go_encrypt(data):
#     data = data.encode("utf-8")
#     resp = file_xor.encrypt(data).decode("utf-8")
#     return resp
#
#
# def go_decrypt(data):
#     data = data.encode("utf-8")
#     resp = file_xor.decrypt(data).decode("utf-8")
#     return resp


