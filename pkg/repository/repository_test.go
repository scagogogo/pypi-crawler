package repository

import (
	"context"
	project_root_directory "github.com/golang-infrastructure/go-project-root-directory"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPypiRepository_GetPackage(t *testing.T) {
	pkg, err := NewRepository().GetPackage(context.Background(), "requests")
	assert.Nil(t, err)
	assert.NotNil(t, pkg)
}

func TestPypiRepository_DownloadIndex(t *testing.T) {
	indexPageBytes, err := NewRepository().DownloadIndex(context.Background())
	assert.Nil(t, err)
	assert.True(t, len(indexPageBytes) > 0)
}

func TestPypiRepository_ParseIndexPage(t *testing.T) {
	indexFilepath, err := project_root_directory.GetRootFilePath("data/sample.html")
	assert.Nil(t, err)
	indexPageBytes, err := os.ReadFile(indexFilepath)
	assert.Nil(t, err)
	packageIndexes, err := NewRepository().ParseIndexPage(string(indexPageBytes))
	assert.Nil(t, err)
	assert.NotNil(t, packageIndexes)
	assert.True(t, len(packageIndexes) > 0)
}

func RepositoryTest(t *testing.T, r *Repository) {

	// 不为空
	assert.NotNil(t, r)

	// 能够正常获取到索引
	index, err := r.DownloadIndex(context.Background())
	assert.Nil(t, err)
	if err != nil {
		t.Log(err.Error())
	}
	assert.True(t, len(index) > 0)

	// 能够获取到包的信息
	//packageName := index[rand.Intn(len(index))]
	packageName := "requests"
	packageInformation, err := r.GetPackage(context.Background(), packageName)
	assert.Nil(t, err)
	assert.NotNil(t, packageInformation)
	assert.Equal(t, packageName, packageInformation.Information.Name)

}
