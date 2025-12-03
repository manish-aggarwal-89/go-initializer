package helper

import (
	"testing"
	"time"
	"{{MODULE_NAME}}/internal/shared"
)

func TestConvertDateTimeStringToOtherFormat(t *testing.T) {
	type args struct {
		dateTimeString string
		dateTimeFormat string
		targetFormat   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success: date only to date time",
			args: args{
				dateTimeString: "2023-10-01",
				dateTimeFormat: time.DateOnly,
				targetFormat:   shared.DateTimeFormatRFC3339Zoneless,
			},
			want: "2023-10-01T00:00:00",
		},
		{
			name: "failed: error parsing",
			args: args{
				dateTimeFormat: "2025-01-01",
			},
			wantErr: true,
		},
		{
			name: "test",
			args: args{
				dateTimeString: "2025-05-31",
				dateTimeFormat: time.DateOnly,
				targetFormat:   shared.DateTimeFormatRFC3339Zoneless,
			},
			want: "2025-05-31T00:00:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertDateTimeStringToOtherFormat(tt.args.dateTimeString, tt.args.dateTimeFormat, tt.args.targetFormat)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertDateTimeStringToOtherFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertDateTimeStringToOtherFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
