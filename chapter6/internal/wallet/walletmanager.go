package wallet

import (
	"bytes"
	"encoding/gob"
	"errors"
	"gother/chapter6/internal/constant"
	"gother/chapter6/internal/utils"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type RefList map[string]string

func (r *RefList) Save() {
	filename := constant.WalletsRefList + "ref_list.data"
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(r)
	utils.Handle(err)
	err = ioutil.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}

func (r *RefList) Update() {
	// 扫描钱包目录下的文件
	err := filepath.Walk(constant.Wallets, func(path string, f fs.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		fileName := f.Name()
		// 判断文件后缀名是否是钱包文件后缀
		if strings.Compare(fileName[len(fileName)-4:], ".wlt") == 0 {
			address := fileName[:len(fileName)-4]
			// 赋值
			_, ok := (*r)[address]
			if !ok {
				(*r)[address] = ""
			}
		}
		return nil
	})
	utils.Handle(err)
}

func (r *RefList) BindRef(address string, refName string) {
	(*r)[address] = refName
}

func (r *RefList) FindRef(refName string) (string, error) {
	temp := ""
	for key, val := range *r {
		if val == refName {
			temp = key
			break
		}
	}

	// 没找到
	if temp == "" {
		err := errors.New("the refName is not found")
		return temp, err
	}
	return temp, nil
}

func LoadRefList() *RefList {
	filename := constant.WalletsRefList + "ref_list.data"
	var refList RefList
	// 如果文件存在则加载文件
	if utils.FileExists(filename) {
		fileContent, err := ioutil.ReadFile(filename)
		utils.Handle(err)
		decoder := gob.NewDecoder(bytes.NewReader(fileContent))
		err = decoder.Decode(&refList)
		utils.Handle(err)
	} else {
		// 如果不存在在则扫描钱包文件并加载
		refList = make(RefList)
		refList.Update()
	}
	return &refList
}
