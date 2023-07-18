package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

func main(){
	// 读取目录下所有文件
	err := filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
		// 判断是否为目录
		if info.IsDir(){
			return nil
		}
		// 计算当前深度
		depth := 0
		for _, c := range path {
			if c == os.PathSeparator {
				depth++
			}
		}
		if depth > 0{
			return nil
		}
		ext := filepath.Ext(path)
		if ext == ""{ // 没有后缀名
			return nil
		}
		// 创建目录 & 移动
		dir := ext[1:]
		_, err = os.Stat(dir)
		if os.IsNotExist(err){ // 文件不存在 则创建文件
			if err = os.Mkdir(dir, os.ModeDir); err != nil{
				return err
			}
		}
		if err = os.Rename(path, ext[1:]+"/"+path); err!=nil{
			return err
		}
		return nil
	})
	// 错误❌处理
	if err != nil{
		panic(err)
	}
}
