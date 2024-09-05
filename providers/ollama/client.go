package ollama

import (
	"github.com/ollama/ollama/api"
	"github.com/ollama/ollama/envconfig"
	"go.qingyu31.com/llmbridge/llm"
	"net/http"
	"net/url"
)

type Client struct {
	client *api.Client
}

func New(opts ...llm.ClientOption[*options]) (llm.Client, error) {
	c := new(Client)
	os := new(options)
	os.InitWithDefault()
	for _, o := range opts {
		o.Apply(os)
	}
	c.client = api.NewClient(os.baseUrl, os.httpClient)
	return c, nil
}

type options struct {
	baseUrl    *url.URL
	httpClient *http.Client
}

func (o *options) InitWithDefault() {
	o.baseUrl = envconfig.Host()
	o.httpClient = http.DefaultClient
}

func WithHTTPClient(client *http.Client) llm.ClientOption[*options] {
	return llm.ClientOptionFunc[*options](func(o *options) {
		o.httpClient = client
	})
}

func WithBaseUrl(u *url.URL) llm.ClientOption[*options] {
	return llm.ClientOptionFunc[*options](func(o *options) {
		o.baseUrl = u
	})
}
