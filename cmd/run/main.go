package main

import (
	"encoding/json"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-querystring/query"
	client "github.com/synology-community/synology-api/package"
	"github.com/synology-community/synology-api/package/util/form"
)

type Options struct {
	Query   string `url:"q"`
	ShowAll bool   `url:"all"`
	Page    int    `url:"page"`
	Pages   []int  `url:"pages" del:","`
}

type File struct {
	Name    string `form:"name" url:"name"`
	Content string `form:"content" url:"content"`
}

type FileTest struct {
	Api     string `form:"api" url:"api"`
	Version string `form:"version" url:"version"`
	Method  string `form:"method" url:"method"`

	File File `form:"file" kind:"file"`
}

func tst() {

	c := FileTest{
		Api:     "SYNO.FileStation.Upload",
		Version: "2",
		Method:  "upload",
		File: File{
			Name:    "main.go",
			Content: "package main",
		},
	}

	if res, err := form.Marshal(&c); err != nil {
		log.Fatal(err)
	} else {
		log.Info(string(res))
		// multipart.File
	}

	opt := Options{"foo", true, 2, []int{1, 2, 3}}
	v, _ := query.Values(opt)
	s := v.Encode() // will output: "q=foo&all=true&page=2"
	log.Info(url.QueryUnescape(s))

}

func main() {

	log.SetFormatter(&log.JSONFormatter{})

	log.Info("Starting")

	// tst()

	// return

	host := "https://appkins.synology.me:5001" // os.Getenv("SYNOLOGY_HOST")
	user := "terraform"                        // os.Getenv("SYNOLOGY_USER")
	password := "ach2vzw*dnx5BPV9njr"          // os.Getenv("SYNOLOGY_PASSWORD")

	client, err := client.New(host, true)

	if err != nil {
		panic(err)
	}

	_, err = client.Login(user, password, "webui")

	if err != nil {
		panic(err)
	}

	_, err = client.FileStationAPI().Upload("/data/foo/bar", &form.File{Name: "main.go", Content: "package main"}, true, true)

	if err != nil {
		panic(err)
	}

	if _, err := client.FileStationAPI().Upload("/data/foo/bar", &form.File{Name: "main.go", Content: "package main"}, true, true); err != nil {
		panic(err)
	}

	listGuestResp, err := client.VirtualizationAPI().ListGuests()

	if err != nil {
		panic(err)
	}

	listGuestRespBytes, _ := json.Marshal(listGuestResp)

	println(string(listGuestRespBytes))

	for _, guest := range listGuestResp.Guests {
		println(guest.Name)
	}

	createFolder(client)
}

func createFolder(client client.SynologyClient) {
	resp, err := client.FileStationAPI().CreateFolder([]string{"/data/foo"}, []string{"bar"}, true)

	if err != nil {
		panic(err)
	}

	for _, folder := range resp.Folders {
		println(folder.Path)
		println(folder.Name)
		println(folder.IsDir)
	}
}
