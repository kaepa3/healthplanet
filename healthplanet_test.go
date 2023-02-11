package healthplanet

import "testing"

func Test_getStatusStr(t *testing.T) {
	type args struct {
		st Status
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "general", args: args{st: Innerscan}, want: "innerscan", wantErr: false},
		{name: "general`", args: args{st: Sphygmomanometer}, want: "sphygmomanometer", wantErr: false},
		{name: "error", args: args{st: 33}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getStatusStr(tt.args.st)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStatusStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getStatusStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
