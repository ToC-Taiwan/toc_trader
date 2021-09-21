package stockutil

import "testing"

func TestGetNewClose(t *testing.T) {
	type args struct {
		close float64
		unit  int64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Input 0",
			args: args{
				close: 0,
				unit:  1,
			},
			want: 0,
		},
		{
			name: "Input 10 and 1",
			args: args{
				close: 10,
				unit:  1,
			},
			want: 10.05,
		},
	}
	for _, tt := range tests {
		close := tt.args.close
		unit := tt.args.unit
		isError := tt.want
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNewClose(close, unit); got != isError {
				t.Errorf("GetNewClose() = %v, want %v", got, isError)
			}
		})
	}
}
