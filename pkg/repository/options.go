package repository

type Options struct {

	// 仓库服务器的地址
	ServerURL string

	// 请求时使用的代理IP
	Proxy string
}

const DefaultServerURL = "https://pypi.org"

func NewOptions() *Options {
	return &Options{
		ServerURL: DefaultServerURL,
	}
}

func (x *Options) SetServerURL(serverURL string) *Options {
	x.ServerURL = serverURL
	return x
}

func (x *Options) SetProxy(proxy string) *Options {
	x.Proxy = proxy
	return x
}

