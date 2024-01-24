package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/uussoop/nextcloudhelper"
	gonextcloud "github.com/uussoop/nxtcloudgo"
)

var Cclient gonextcloud.Client

const Cclienturl = "https://storagebox.dentai.org"

var ShareLinkNotFoundLinkErr error = errors.New("share link not found")

func main() {

	url := ""
	username := ""
	password := ""
	c, _ := nextcloudhelper.GetClient(url, username, password)

	defer c.Logout()

	//example get or create
	getorcreate, extras, _ := c.GetOrCreateShareLink("./Storage-Share.png")
	fmt.Println(getorcreate)

	//example remove
	shareid, _ := strconv.ParseInt(extras.ID, 10, 64)
	c.RemoveShareLink(int(shareid))

	//example upload file
	data, _ := os.ReadFile("./Storage-Share2.png")

	uped := c.UploadFile(data, "Storage-Share2.png")
	fmt.Println(uped)

}
