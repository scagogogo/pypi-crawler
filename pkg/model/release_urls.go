package model

// ReleaseUrls 表示一个版本发布的release的集合，可能会有多个平台啥的
type ReleaseUrls []*ReleaseURL

func (x ReleaseUrls) Len() int {
	return len(x)
}
