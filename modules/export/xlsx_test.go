package export

import (
	"testing"
	"time"

	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/stretchr/testify/assert"
)

func Test_ExportModule_getCellValue(t *testing.T) {
	t.Parallel()

	mockTime := time.Date(2025, time.January, 16, 12, 17, 34, 0, time.UTC).Format(time.RFC3339Nano)

	type args struct {
		data   map[string]interface{}
		column *dto.Column
	}

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "valid",
			args: args{
				data: map[string]interface{}{
					"id": 5,
				},
				column: &dto.Column{
					Name:  "id",
					Style: IntColumnStyle,
				},
			},
			want: float64(5),
		},
		{
			name: "valid_DayColumnStyle",
			args: args{
				data: map[string]interface{}{
					"date": mockTime,
				},
				column: &dto.Column{
					Name:  "date",
					Style: DayColumnStyle,
				},
			},
			want: "2025-01-16",
		},
		{
			name: "valid_booleanReplacement_true",
			args: args{
				data: map[string]interface{}{
					"bool": true,
				},
				column: &dto.Column{
					Name:  "bool",
					Style: BooleanColumnStyle,
					BooleanReplacement: &dto.BooleanReplacement{
						True:  "is_bool",
						False: "is_not_bool",
					},
				},
			},
			want: "is_bool",
		},
		{
			name: "valid_booleanReplacement_false",
			args: args{
				data: map[string]interface{}{
					"bool": false,
				},
				column: &dto.Column{
					Name:  "bool",
					Style: BooleanColumnStyle,
					BooleanReplacement: &dto.BooleanReplacement{
						True:  "is_bool",
						False: "is_not_bool",
					},
				},
			},
			want: "is_not_bool",
		},
		{
			name: "valid_intWithMultiply",
			args: args{
				data: map[string]interface{}{
					"gpp": 55,
				},
				column: &dto.Column{
					Name:     "gpp",
					Style:    PercentageColumnStyle,
					Multiply: helpers.GetPointerToFloat64(0.01),
				},
			},
			want: 0.55,
		},
		{
			name: "valid_floatWithMultiply",
			args: args{
				data: map[string]interface{}{
					"gpp": 55.5,
				},
				column: &dto.Column{
					Name:     "gpp",
					Style:    PercentageColumnStyle,
					Multiply: helpers.GetPointerToFloat64(0.01),
				},
			},
			want: 0.555,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := NewExportModule()

			got := e.getCellValue(tt.args.data, tt.args.column)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_ExportModule_getFormattedDatetimeString(t *testing.T) {
	t.Parallel()

	mockTime := time.Date(2025, time.January, 16, 12, 17, 34, 0, time.UTC).Format(time.RFC3339Nano)

	type args struct {
		value interface{}
		style string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid_HourColumnStyle",
			args: args{
				value: mockTime,
				style: HourColumnStyle,
			},
			want: "2025-01-16 12:00:00",
		},
		{
			name: "valid_DayColumnStyle",
			args: args{
				value: mockTime,
				style: DayColumnStyle,
			},
			want: "2025-01-16",
		},
		{
			name: "valid_WeekColumnStyle",
			args: args{
				value: mockTime,
				style: WeekColumnStyle,
			},
			want: "03 2025",
		},
		{
			name: "valid_MonthColumnStyle",
			args: args{
				value: mockTime,
				style: MonthColumnStyle,
			},
			want: "January, 2025",
		},
		{
			name: "valid_QuarterColumnStyle",
			args: args{
				value: mockTime,
				style: QuarterColumnStyle,
			},
			want: "Q1 2025",
		},
		{
			name: "valid_YearColumnStyle",
			args: args{
				value: mockTime,
				style: YearColumnStyle,
			},
			want: "2025",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := NewExportModule()

			got := e.getFormattedDatetimeString(tt.args.value, tt.args.style)
			assert.Equal(t, tt.want, got)
		})
	}
}
