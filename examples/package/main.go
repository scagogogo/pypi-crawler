package main

import (
	"context"
	"fmt"
	"github.com/scagogogo/pypi-crawler/pkg/repository"
)

func main() {

	// 获取名为requests的包的信息
	r := repository.NewRepository()
	packageInformation, err := r.GetPackage(context.Background(), "requests")
	if err != nil {
		panic(err)
	}
	fmt.Println(packageInformation.Information.Name)
	fmt.Println(packageInformation.Information.Version)
	fmt.Println(packageInformation.Information.Description)

	// Output:
	// requests
	//2.31.0
	//# Requests
	//
	//**Requests** is a simple, yet elegant, HTTP library.
	// ...

}
