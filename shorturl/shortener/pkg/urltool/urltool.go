package urltool

import (
	"net/url"
	"path"
)

func GetBasePath(targeturl string) (string, error) {
	myUrl, err := url.Parse(targeturl)
	if err != nil {
		return "", err
	}
	basePath := path.Base(myUrl.Path)
	return basePath, nil
}
