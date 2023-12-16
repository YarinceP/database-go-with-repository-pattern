package database

import "testing"

func Test_getConnectionString(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "OK",
			want: "root:@tcp(localhost:3306)/db_lib_go?parseTime=true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConnectionString(); got != tt.want {
				t.Errorf("getConnectionString() = %v, want %v", got, tt.want)
			}
		})
	}
}
