package model

import (
	"encoding/json"
	"fmt"

	"github.com/emirpasic/gods/maps/linkedhashmap"
)

// Releases 表示包的所有历史版本信息集合
//
// 包含了版本号的顺序以及每个版本对应的发布文件信息
// 实现了自定义的JSON序列化和反序列化，以保持版本的有序性
//
// 使用示例:
//
//	pkg, err := repo.GetPackage(context.Background(), "requests")
//	if err == nil {
//	    releases := pkg.Releases
//	    fmt.Printf("该包共有 %d 个历史版本\n", len(releases.VersionOrders))
//
//	    // 按顺序访问所有版本
//	    for _, version := range releases.VersionOrders {
//	        releaseUrls := releases.VersionMap[version]
//	        fmt.Printf("版本 %s 有 %d 个发布文件\n", version, releaseUrls.Len())
//	    }
//
//	    // 获取特定版本的信息
//	    specificVersion := "2.28.1"
//	    if releaseUrls, ok := releases.VersionMap[specificVersion]; ok {
//	        fmt.Printf("版本 %s 的详细信息:\n", specificVersion)
//	        for i, url := range *releaseUrls {
//	            fmt.Printf("  文件 #%d: %s (%s)\n", i+1, url.Filename, url.Packagetype)
//	        }
//	    }
//	}
type Releases struct {
	// VersionOrders 所有版本号的有序列表
	// 按照版本发布的时间顺序排列
	VersionOrders []string

	// VersionMap 版本号到对应发布文件集合的映射
	// 键为版本号字符串，值为该版本的发布文件集合
	VersionMap map[string]*ReleaseUrls
}

var _ json.Marshaler = &Releases{}
var _ json.Unmarshaler = &Releases{}

// UnmarshalJSON 自定义JSON反序列化方法
//
// 此方法保证在反序列化过程中保留版本顺序
// 首先使用linkedhashmap反序列化以获取版本顺序，然后作为普通map反序列化内容
//
// 参数:
//   - bytes: 要反序列化的JSON字节数据
//
// 返回值:
//   - error: 反序列化过程中的错误
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

// MarshalJSON 自定义JSON序列化方法
//
// 此方法确保在序列化过程中保持版本的顺序
// 使用linkedhashmap来保证版本顺序在序列化后得以保留
//
// 返回值:
//   - []byte: 序列化后的JSON字节数据
//   - error: 序列化过程中的错误
func (x *Releases) MarshalJSON() ([]byte, error) {
	m := linkedhashmap.New()
	for _, version := range x.VersionOrders {
		m.Put(version, x.VersionMap[version])
	}
	return m.ToJSON()
}
