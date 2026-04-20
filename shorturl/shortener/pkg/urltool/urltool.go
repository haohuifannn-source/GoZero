package urltool

import (
	"errors"
	"net/url"
	"path"
)

func GetBasePath(targeturl string) (string, error) {
	myUrl, err := url.Parse(targeturl)
	if err != nil {
		return "", err
	}
	if len(myUrl.Host) == 0 {
		return "", errors.New("Host is null") //防止传入的是相对路径
	}
	basePath := path.Base(myUrl.Path)
	return basePath, nil
}
