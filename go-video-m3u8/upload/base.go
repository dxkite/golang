package upload

import "errors"

type FileObject struct {
	// 文件名
	Name string
	// 文件内容
	Data []byte
}

type Result struct {
	// 上传URL
	Url string
	// 原始数据
	Raw []byte
}

var list = make(map[string]Uploader)

type Uploader interface {
	Upload(object *FileObject) (*Result, error)
}

func Register(name string, uploader Uploader) {
	list[name] = uploader
}

func Upload(name string, object *FileObject) (*Result, error) {
	if uploader, ok := list[name]; ok {
		return uploader.Upload(object)
	}
	return nil, errors.New("unknown uploader:" + name)
}
