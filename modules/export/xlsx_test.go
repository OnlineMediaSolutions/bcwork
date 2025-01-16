package export

import (
	"testing"
	"time"

	"github.com/m6yf/bcwork/dto"
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
			want: 5,
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

func Test_getQuarter(t *testing.T) {
	t.Parallel()

	type args struct {
		month time.Month
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Q1_January",
			args: args{
				month: time.January,
			},
			want: "Q1",
		},
		{
			name: "Q1_February",
			args: args{
				month: time.February,
			},
			want: "Q1",
		},
		{
			name: "Q1_March",
			args: args{
				month: time.March,
			},
			want: "Q1",
		},
		{
			name: "Q2_April",
			args: args{
				month: time.April,
			},
			want: "Q2",
		},
		{
			name: "Q2_May",
			args: args{
				month: time.May,
			},
			want: "Q2",
		},
		{
			name: "Q2_June",
			args: args{
				month: time.June,
			},
			want: "Q2",
		},
		{
			name: "Q3_July",
			args: args{
				month: time.July,
			},
			want: "Q3",
		},
		{
			name: "Q3_August",
			args: args{
				month: time.August,
			},
			want: "Q3",
		},
		{
			name: "Q3_September",
			args: args{
				month: time.September,
			},
			want: "Q3",
		},
		{
			name: "Q4_October",
			args: args{
				month: time.October,
			},
			want: "Q4",
		},
		{
			name: "Q4_November",
			args: args{
				month: time.November,
			},
			want: "Q4",
		},
		{
			name: "Q4_December",
			args: args{
				month: time.December,
			},
			want: "Q4",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := getQuarter(tt.args.month)
			assert.Equal(t, tt.want, got)
		})
	}
}
