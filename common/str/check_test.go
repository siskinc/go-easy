package str

import "testing"

func TestCheckPhoneNumber(t *testing.T) {
	type args struct {
		phoneNumber string
	}
	tests := []struct {
		name        string
		args        args
		wantMatched bool
		wantErr     bool
	}{
		{
			"17892020997",
			args{phoneNumber: "17892020997"},
			true,
			false,
		},
		{
			"18882020998",
			args{phoneNumber: "18882020998"},
			true,
			false,
		},
		{
			"1888202099",
			args{phoneNumber: "1888202099"},
			false,
			false,
		},
		{
			"188820209999",
			args{phoneNumber: "188820209999"},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatched, err := CheckPhoneNumber(tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMatched != tt.wantMatched {
				t.Errorf("CheckPhoneNumber() gotMatched = %v, want %v", gotMatched, tt.wantMatched)
			}
		})
	}
}
