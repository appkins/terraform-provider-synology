package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/appkins/terraform-provider-synology/synology/client/api"
	"github.com/appkins/terraform-provider-synology/synology/client/api/filestation"
	"github.com/appkins/terraform-provider-synology/synology/client/api/virtualization"
	"github.com/appkins/terraform-provider-synology/synology/client/util/form"
	"golang.org/x/net/publicsuffix"
)

type SynologyClient interface {
	Login(user, password, sessionName string) error
	CreateFolder(folderPath string, name string, forceParent bool) (*filestation.CreateFolderResponse, error)
	ListShares() (*filestation.ListShareResponse, error)
	GetGuest(name string) (*virtualization.GetGuestResponse, error)
	ListGuests() (*virtualization.ListGuestResponse, error)
	Upload(path string, file *form.File, createParents bool, overwrite bool) error
}
type synologyClient struct {
	httpClient *http.Client

	baseURL url.URL
}

// Upload implements SynologyClient.
func (c *synologyClient) Upload(path string, file *form.File, createParents bool, overwrite bool) error {
	request := filestation.NewUploadRequest(path, file)
	request.WithCreateParents(createParents)
	request.WithOverwrite(overwrite)

	response := filestation.UploadResponse{}

	return c.Post(request, &response)
}

// GetGuest implements Client.
func (c *synologyClient) GetGuest(name string) (*virtualization.GetGuestResponse, error) {

	request := virtualization.NewGetGuestRequest(name)
	response := virtualization.GetGuestResponse{}
	if err := c.Get(request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// ListGuests implements Client.
func (c *synologyClient) ListGuests() (*virtualization.ListGuestResponse, error) {
	return GetResponse[*virtualization.ListGuestResponse](c, api.NewRequest("SYNO.Virtualization.API.Guest", "list"))
}

// ListShares implements Client.
func (c *synologyClient) ListShares() (*filestation.ListShareResponse, error) {
	panic("unimplemented")
}

func (c synologyClient) CreateFolder(folderPath string, name string, forceParent bool) (*filestation.CreateFolderResponse, error) {
	request := filestation.NewCreateFolderRequest(2)
	request.WithFolderPath(folderPath)
	request.WithName(name)
	request.WithForceParent(forceParent)

	response := filestation.CreateFolderResponse{}

	err := c.Get(request, &response)

	return &response, err
}

// New initializes "client" instance with minimal input configuration.
func New(host string, skipCertificateVerification bool) (SynologyClient, error) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipCertificateVerification,
		},
	}

	// currently, 'Cookie' is the only supported method for providing 'sid' token to DSM
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Transport: transport,
		Jar:       jar,
	}

	baseURL, err := url.Parse(host)

	baseURL.Scheme = "https"
	baseURL.Path = "/webapi/entry.cgi"

	if err != nil {
		return nil, err
	}

	return &synologyClient{
		httpClient: httpClient,
		baseURL:    *baseURL,
	}, nil
}
