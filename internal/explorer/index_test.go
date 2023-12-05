package explorer

import (
	"fmt"
	"strings"
	"testing"
)

//func TestExplorer_GetSourceCode(t *testing.T) {
//
//	defaultConfig := &Explorer{
//		ApiGetSourceCode: "https://api.etherscan.io/api",
//		ApiKey:           "862Y3WJ4JB4B34PZQRFEV3IK6SZ8GNR9N5",
//	}
//
//	tests := []struct {
//		name   string
//		config *Explorer
//		target string
//		want   string
//	}{
//		// TODO: Add test cases.
//		{
//			name: "Fail. Incorrect url 1",
//			config: &Explorer{
//				ApiGetSourceCode: defaultConfig.ApiGetSourceCode,
//				ApiKey:           defaultConfig.ApiKey,
//			},
//			target: "address/0x000000000",
//			want:   "error while get source code or source code is empty",
//		},
//		{
//			name: "Success. Contract is clear",
//			config: &Explorer{
//				ApiGetSourceCode: defaultConfig.ApiGetSourceCode,
//				ApiKey:           defaultConfig.ApiKey,
//			},
//			target: "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
//			want:   "",
//		},
//		{
//			name: "Success. Contract is flattened",
//			config: &Explorer{
//				ApiGetSourceCode: defaultConfig.ApiGetSourceCode,
//				ApiKey:           defaultConfig.ApiKey,
//			},
//			target: "0x9abF23f4e439d695A7FD341a1b25873C50CFa52e",
//			want:   "",
//		},
//		{
//			name: "Success. Contract is a proxy",
//			config: &Explorer{
//				ApiGetSourceCode: defaultConfig.ApiGetSourceCode,
//				ApiKey:           defaultConfig.ApiKey,
//			},
//			target: "0xDEF171Fe48CF0115B1d80b88dc8eAB59176FEe57",
//			want:   "",
//		},
//	}
//	//for _, tt := range tests {
//	//	t.Run(tt.name, func(t *testing.T) {
//	//		expl := New(tt.config)
//	//		_, err := expl.GetSourceCode(tt.target, 3)
//	//		if !ErrorContains(err, tt.want) {
//	//			t.Errorf("Run() error = %v, wantErr %v", err, tt.want)
//	//		}
//	//	})
//	//}
//}

func ErrorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}

func TestValidChecksum(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Success. Checksum address",
			addr: "0x87cc04d6fe59859cb7eb6d970ebc22dcdcbc9368",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ChecksumAddress(tt.addr)
			fmt.Printf("addr: %s\n", tt.addr)
			fmt.Printf("got: %s\n", got)
		})
	}
}
