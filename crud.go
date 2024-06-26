package nextcloudhelper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	gonextcloud "github.com/uussoop/nxtcloudgo"
)

type CloudClient struct {
	url  string
	user string
	pass string
}

var ShareLinkNotFoundLinkErr error = errors.New("share link not found")

var cclient gonextcloud.Client

func GetClient(url, user, pass string) (*CloudClient, error) {
	client, err := gonextcloud.NewClient(url)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	err = client.Login(user, pass)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	cclient = client
	return &CloudClient{url, user, pass}, nil
}
func (Cclient *CloudClient) Logout() {
	cclient.Logout()
}

func (Cclient *CloudClient) UploadFile(data []byte, path string) error {
	return cclient.WebDav().Write(path, data, 0644)
}
func (Cclient *CloudClient) RemoveShareLink(shareid int) {

	cclient.Shares().Delete(shareid)

}

func (Cclient *CloudClient) GetOrCreateShareLink(
	path string,
) (string, *gonextcloud.Share, error) {
	name := strings.Split(path, "/")
	filename := name[len(name)-1]
	file, err := Cclient.CheckIfShared(filename)
	if err != nil && !errors.Is(err, ShareLinkNotFoundLinkErr) {
		logrus.Error(err)
		return "", nil, err
	}
	if errors.Is(err, ShareLinkNotFoundLinkErr) {

		file, err = cclient.Shares().Create(path, gonextcloud.PublicLinkShare, gonextcloud.ReadPermission, "", false, "")

		if err != nil {
			logrus.Error(err)

		}

	}
	shareUrl := fmt.Sprintf("%s/s/%s/download/%s", Cclient.url, file.Token, filename)

	return shareUrl, &file, nil

}

func (Cclient *CloudClient) CheckIfShared(name string) (gonextcloud.Share, error) {
	sharedlist, err := cclient.Shares().List()

	if err != nil {
		logrus.Error(err)
		return gonextcloud.Share{}, err
	}
	for _, share := range sharedlist {
		listname := strings.Split(share.FileTarget, "/")

		if listname[len(listname)-1] == name {
			return share, nil
		}
	}
	return gonextcloud.Share{}, ShareLinkNotFoundLinkErr

}
