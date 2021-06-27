package product

import (
	"reflect"
	"testing"
)

func Test_generateSellerInfo(t *testing.T) {
	type args struct {
		sellerUUID string
	}
	tests := []struct {
		name string
		args args
		want *SellerInfo
	}{
		// TODO: Add test cases.
		{
			name: "test generate seller info #1",
			args: args{
				sellerUUID: "123",
			},
			want: &SellerInfo{
				UUID: "123",
				Links: &SellerLinks{
					Self: &SelfSellerLink{
						Href: "http://localhost:8080/api/v1/sellers/123",
					},
				},
			},
		},
		{
			name: "test generate seller info #2",
			args: args{
				sellerUUID: "124",
			},
			want: &SellerInfo{
				UUID: "124",
				Links: &SellerLinks{
					Self: &SelfSellerLink{
						Href: "http://localhost:8080/api/v1/sellers/124",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateSellerInfo(tt.args.sellerUUID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateSellerInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
