package video

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type PartUploader func(name string, data []byte) (url string, err error)

func MakeM3u8(prefix, input, output, partDir string, uploader PartUploader) error {
	fi, er := os.OpenFile(input, os.O_RDONLY, os.ModePerm)
	if er != nil {
		return er
	}
	fo, eo := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if eo != nil {
		return eo
	}
	defer fi.Close()
	defer fo.Close()

	br := bufio.NewReader(fi)

	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(a)
		if strings.Index(line, prefix) >= 0 {
			pp := path.Join(partDir, line)
			if b, e := ioutil.ReadFile(pp); e == nil {
				url, eu := uploader(line, b)
				if eu != nil {
					return eu
				}
				if _, e := fo.Write([]byte(url + "\r\n")); e != nil {
					return nil
				}
			} else {
				return e
			}
		} else {
			if _, e := fo.Write([]byte(line + "\r\n")); e != nil {
				return e
			}
		}
	}
	return nil
}
