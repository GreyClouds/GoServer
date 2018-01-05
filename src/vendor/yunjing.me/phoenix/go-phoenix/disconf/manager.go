package disconf

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// 配置管理模块
type Manager struct {
	strPathPrefix string // 目录前缀
}

func New(prefix string) *Manager {
	return &Manager{
		strPathPrefix: prefix,
	}
}

// 加载JSON文件
func (self Manager) LoadJSONFile(filename string, payload interface{}) error {
	path := filepath.Join(self.strPathPrefix, filename)

	if filepath.Ext(path) != ".json" {
		return errors.New("文件务必是.json后缀")
	}

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("读取%s文件时出错: %v", filename, err)
		return err
	}

	if err := json.Unmarshal(raw, payload); err != nil {
		log.Printf("解析%s文件时出错: %v", filename, err)
		return err
	}

	return nil
}

// 加载JSON文件目录
func (self Manager) GetGroupJSONFile(dir string, files *[]string) error {
	return filepath.Walk(filepath.Join(self.strPathPrefix, dir), self.doTraverseJSONFile(files))
}

// 迭代每个目录下的文件
func (self Manager) doTraverseJSONFile(files *[]string) func(string, os.FileInfo, error) error {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return errors.New("配置表目录不存在")
		}

		if f.IsDir() {
			return nil
		}

		// 过滤非JSON文档
		ext := filepath.Ext(path)
		if ext != ".json" {
			return nil
		}

		*files = append(*files, path)

		return nil
	}
}
