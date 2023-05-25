package model

import (
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/maps/linkedhashmap"
)

// Releases 表示同一个包的多个版本
type Releases struct {

	// 版本的顺序
	VersionOrders []string

	// 版本对应的Release的URL
	VersionMap map[string]*ReleaseUrls
}

var _ json.Marshaler = &Releases{}
var _ json.Unmarshaler = &Releases{}

func (x *Releases) UnmarshalJSON(bytes []byte) error {
	// 先反序列化key，把key的顺序保留着
	m := linkedhashmap.New()
	err := m.FromJSON(bytes)
	if err != nil {
		return err
	}
	for _, versionKey := range m.Keys() {
		version, ok := versionKey.(string)
		if !ok {
			return fmt.Errorf("version key must is string type")
		}
		x.VersionOrders = append(x.VersionOrders, version)
	}
	// 然后当做一个map反序列化
	if x.VersionMap == nil {
		x.VersionMap = make(map[string]*ReleaseUrls)
	}
	return json.Unmarshal(bytes, &x.VersionMap)
}

func (x *Releases) MarshalJSON() ([]byte, error) {
	m := linkedhashmap.New()
	for _, version := range x.VersionOrders {
		m.Put(version, x.VersionMap[version])
	}
	return m.ToJSON()
}
