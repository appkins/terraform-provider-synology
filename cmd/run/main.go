package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	client "github.com/appkins/terraform-provider-synology/synology/client"
	"github.com/appkins/terraform-provider-synology/synology/client/api"
	"github.com/appkins/terraform-provider-synology/synology/client/api/filestation"
	"github.com/appkins/terraform-provider-synology/synology/client/util/form"
	"github.com/google/go-querystring/query"
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
		log.Print(string(res))
		// multipart.File
	}

	opt := Options{"foo", true, 2, []int{1, 2, 3}}
	v, _ := query.Values(opt)
	s := v.Encode() // will output: "q=foo&all=true&page=2"
	fmt.Print(url.QueryUnescape(s))

}

func main() {

	// tst()

	// return

	host := "https://appkins.synology.me:5001" // os.Getenv("SYNOLOGY_HOST")
	user := "terraform"                        // os.Getenv("SYNOLOGY_USER")
	password := "ach2vzw*dnx5BPV9njr"          // os.Getenv("SYNOLOGY_PASSWORD")

	client, err := client.New(host, true)

	if err != nil {
		panic(err)
	}

	err = client.Login(user, password, "webui")

	if err != nil {
		panic(err)
	}

	if err := client.Upload("/data/foo/bar", &form.File{Name: "main.go", Content: "package main"}, true, true); err != nil {
		panic(err)
	}

	infoRequest := filestation.NewFileStationInfoRequest(2)
	infoResponse := filestation.FileStationInfoResponse{}

	err = client.Get(infoRequest, &infoResponse)

	if err != nil {
		panic(err)
	}

	println(infoResponse.Hostname)
	println(infoResponse.Supportsharing)

	listGuestResp, err := client.ListGuests()

	if err != nil {
		panic(err)
	}

	listGuestRespBytes, _ := json.Marshal(listGuestResp)

	println(string(listGuestRespBytes))

	for _, guest := range listGuestResp.Guests {
		println(guest.Name)
	}

	createFolder(client)

	if err := Upload(client, host, mustOpen("main.go")); err != nil {
		panic(err)
	}
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func createFolder(client client.SynologyClient) {
	resp, err := client.CreateFolder("/data/foo", "bar", true)

	if err != nil {
		panic(err)
	}

	for _, folder := range resp.Folders {
		println(folder.Path)
		println(folder.Name)
		println(folder.IsDir)
	}
}

func Upload(client client.SynologyClient, host string, r io.Reader) error {
	data := map[string]string{
		"api":            "SYNO.FileStation.Upload",
		"version":        "2",
		"method":         "upload",
		"path":           "/data/foo/bar",
		"create_parents": "true",
		"overwrite":      "true",
	}

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for key, val := range data {
		w.WriteField(key, val)
	}

	if x, ok := r.(io.Closer); ok {
		defer x.Close()
	}

	// Add an image file
	if x, ok := r.(*os.File); ok {
		if fw, err := w.CreateFormFile("file", x.Name()); err != nil {
			return err
		} else {
			if _, err := io.Copy(fw, r); err != nil {
				return err
			}
		}
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "/webapi/entry.cgi",
	}

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest(http.MethodPost, u.String(), &b)
	if err != nil {
		return err
	}

	// file, err := os.ReadFile("main.go")
	// if err != nil {
	// 	return err
	// }

	// log.Print(req.FormValue("api"))

	// Don't forget to set the content type, this will contain the boundary.
	// req.Header.Set("Content-Disposition", `form-data; name="file"; filename="main.go"`)
	// req.Header.Set("Content-Type", "application/octet-stream")
	// req.Header.Set("Content-Length", fmt.Sprintf("%d", b.Len()))
	//
	// req.Form.Set("file", string(file))

	// req.PostForm = url.Values{}

	// req.PostForm.Set("api", string("SYNO.FileStation.Upload"))
	// req.PostForm.Set("version", string("2"))
	// req.PostForm.Set("method", string("upload"))

	log.Print(req.URL.RequestURI())

	var res = api.BaseResponse{}

	// Submit the request
	err = client.Do(req, &res)
	if err != nil {
		return err
	}

	// Check the response
	if !res.Success() {
		return fmt.Errorf("bad status: %s", res.GetError())
	}

	log.Print("Upload successful")

	return nil
}
