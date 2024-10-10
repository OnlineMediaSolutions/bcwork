package constant

import (
	"testing"
	"time"

	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func Test_Targeting_PrepareData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		targeting Targeting
		want      Targeting
	}{
		{
			name: "valid",
			targeting: Targeting{
				Country:    []string{"ru", "uk", "ar", "us", "il"},
				DeviceType: []string{"mobile", "web", "desktop"},
				Browser:    []string{"opera", "chrome", "firefox", "edge"},
				OS:         []string{"windows", "macos", "linux"},
			},
			want: Targeting{
				Country:    []string{"ar", "il", "ru", "uk", "us"},
				DeviceType: []string{"desktop", "mobile", "web"},
				Browser:    []string{"chrome", "edge", "firefox", "opera"},
				OS:         []string{"linux", "macos", "windows"},
				Status:     TargetingStatusActive,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.targeting.PrepareData()
			// object mutates
			assert.Equal(t, tt.want, tt.targeting)
		})
	}
}

func Test_GetTargetingRegExp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		targeting *models.Targeting
		want      string
		wantErr   bool
	}{
		{
			name: "valid",
			targeting: &models.Targeting{
				PublisherID:   "1",
				Domain:        "1.com",
				UnitSize:      "300x200",
				PlacementType: null.StringFrom("placement"),
				Country:       []string{"ru", "uk", "us", "il"},
				DeviceType:    []string{"mobile", "desktop"},
				Browser:       []string{"chrome", "firefox", "edge"},
				Os:            []string{"windows", "macos", "linux"},
				KV:            null.JSONFrom([]byte(`{"key1": "value1", "key2": "value2", "key3": "value3"}`)),
			},
			want: "p=1__d=1.com__s=300x200__c=(ru|uk|us|il)__os=(windows|macos|linux)__dt=(mobile|desktop)__pt=placement__b=(chrome|firefox|edge)__key1=value1__key2=value2__key3=value3",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := GetTargetingRegExp(tt.targeting)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_GetTargetingKey(t *testing.T) {
	t.Parallel()

	type args struct {
		publisher string
		domain    string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid",
			args: args{
				publisher: "publisher",
				domain:    "1.com",
			},
			want: "jstag:publisher:1.com",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GetTargetingKey(tt.args.publisher, tt.args.domain)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_GetJSTagString(t *testing.T) {
	t.Parallel()

	now := time.Now().Format(time.DateOnly)

	type args struct {
		mod     *models.Targeting
		addGDPR bool
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				mod: &models.Targeting{
					PublisherID: "1",
					Domain:      "1.com",
					UnitSize:    "300x250",
				},
			},
			want: "<!-- HTML Tag for publisher='', domain='1.com', size='300x250', exported='" + now + "' -->\n" +
				"<script src=\"https://rt.marphezis.com/js?pid=1&size=300x250&dom=1.com\"></script>",
		},
		{
			name: "valid_WithKV",
			args: args{
				mod: &models.Targeting{
					PublisherID: "1",
					Domain:      "1.com",
					UnitSize:    "300x250",
					KV:          null.JSONFrom([]byte(`{"key_1":"value_1","key_2":"value_2"}`)),
				},
			},
			want: "<!-- HTML Tag for publisher='', domain='1.com', size='300x250', key_1='value_1', key_2='value_2', exported='" + now + "' -->\n" +
				"<script src=\"https://rt.marphezis.com/js?pid=1&size=300x250&dom=1.com&key_1=value_1&key_2=value_2\"></script>",
		},
		{
			name: "valid_WithGDPR",
			args: args{
				mod: &models.Targeting{
					PublisherID: "1",
					Domain:      "1.com",
					UnitSize:    "300x250",
				},
				addGDPR: true,
			},
			want: "<!-- HTML Tag for publisher='', domain='1.com', size='300x250', exported='" + now + "' -->\n" +
				"<script src=\"https://rt.marphezis.com/js?pid=1&size=300x250&dom=1.com&gdpr=${GDPR}&gdpr_concent=${GDPR_CONSENT_883}\"></script>",
		},
		{
			name: "valid_WithKVAndGDPR",
			args: args{
				mod: &models.Targeting{
					PublisherID: "1",
					Domain:      "1.com",
					UnitSize:    "300x250",
					KV:          null.JSONFrom([]byte(`{"key_1":"value_1","key_2":"value_2"}`)),
				},
				addGDPR: true,
			},
			want: "<!-- HTML Tag for publisher='', domain='1.com', size='300x250', key_1='value_1', key_2='value_2', exported='" + now + "' -->\n" +
				"<script src=\"https://rt.marphezis.com/js?pid=1&size=300x250&dom=1.com&key_1=value_1&key_2=value_2&gdpr=${GDPR}&gdpr_concent=${GDPR_CONSENT_883}\"></script>",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := GetJSTagString(tt.args.mod, tt.args.addGDPR)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
