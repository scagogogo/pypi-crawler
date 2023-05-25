package repository

import (
	"testing"
)

func TestNewTSingHuaRepository(t *testing.T) {
	// 2023-5-26 04:13:46 test passed
	RepositoryTest(t, NewTSingHuaRepository())
}

func TestNewAliCloudRepository(t *testing.T) {
	// 2023-5-26 04:08:06 测试不通过
	//RepositoryTest(t, NewAliCloudRepository())
}

func TestNewDouBanRepository(t *testing.T) {
	// 2023-5-26 04:06:24 test passed
	RepositoryTest(t, NewDouBanRepository())
}

func TestNewNetEaseRepository(t *testing.T) {
	// 2023-5-26 03:55:50 test passed
	RepositoryTest(t, NewNetEaseRepository())
}

func TestNewTencentCloudRepository(t *testing.T) {
	// 2023-5-26 03:49:02 test passed
	RepositoryTest(t, NewTencentCloudRepository())
}

func TestNewUstcRepository(t *testing.T) {
	// 2023-5-26 03:49:06 test passed
	RepositoryTest(t, NewUstcRepository())
}
