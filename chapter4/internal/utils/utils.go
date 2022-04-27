package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

func ToHexForInt(input int64) []byte {
	// 创建一个字节buffer,未初始化零值
	buff := new(bytes.Buffer)
	// 按大端顺序，写入数据
	err := binary.Write(buff, binary.BigEndian, input)
	if err != nil {
		log.Panicln(err)
	}
	// 返回转换后的字节数组
	return buff.Bytes()
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// FileExists 判断文件是否存在
func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}
