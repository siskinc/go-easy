package ginx

import "testing"

func TestMergeURL(t *testing.T) {
	type args struct {
		first  string
		second string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test1",
			args{"/", "/"},
			"/",
		},
		{
			"test2",
			args{"/", ""},
			"/",
		},
		{
			"test3",
			args{"/", "/auth"},
			"/auth",
		},
		{
			"test4",
			args{"/////daryl//////", "////au/th/////"},
			"/daryl/au/th",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeURL(tt.args.first, tt.args.second); got != tt.want {
				t.Errorf("MergeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
