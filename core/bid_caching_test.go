package core

import (
	"encoding/json"
	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"testing"
)

func Test_BC_ToModel(t *testing.T) {
	t.Parallel()

	type args struct {
		bidCaching *BidCaching
	}

	tests := []struct {
		name     string
		args     args
		expected *models.BidCaching
	}{
		{
			name: "All fields populated",
			args: args{
				bidCaching: &BidCaching{
					RuleId:        "50afedac-d41a-53b0-a922-2c64c6e80623",
					Publisher:     "Publisher1",
					Domain:        "example.com",
					BidCaching:    1,
					OS:            "Windows",
					Country:       "US",
					Device:        "Desktop",
					PlacementType: "Banner",
					Browser:       "Chrome",
				},
			},
			expected: &models.BidCaching{
				RuleID:        "50afedac-d41a-53b0-a922-2c64c6e80623",
				Publisher:     "Publisher1",
				Domain:        "example.com",
				BidCaching:    1,
				Country:       null.StringFrom("US"),
				Os:            null.StringFrom("Windows"),
				Device:        null.StringFrom("Desktop"),
				PlacementType: null.StringFrom("Banner"),
				Browser:       null.StringFrom("Chrome"),
			},
		},
		{
			name: "Some fields empty",
			args: args{
				bidCaching: &BidCaching{
					RuleId:        "d823a92a-83e5-5c2b-a067-b982d6cdfaf8",
					Publisher:     "Publisher2",
					Domain:        "example.org",
					BidCaching:    1,
					OS:            "",
					Country:       "CA",
					Device:        "",
					PlacementType: "Sidebar",
					Browser:       "",
				},
			},
			expected: &models.BidCaching{
				RuleID:        "d823a92a-83e5-5c2b-a067-b982d6cdfaf8",
				Publisher:     "Publisher2",
				Domain:        "example.org",
				BidCaching:    1,
				Country:       null.StringFrom("CA"),
				Os:            null.String{},
				Device:        null.String{},
				PlacementType: null.StringFrom("Sidebar"),
				Browser:       null.String{},
			},
		},
		{
			name: "All fields empty",
			args: args{
				bidCaching: &BidCaching{
					RuleId:        "966affd7-d087-57a2-baff-55b926f4c32d",
					Publisher:     "",
					Domain:        "",
					BidCaching:    1,
					OS:            "",
					Country:       "",
					Device:        "",
					PlacementType: "",
					Browser:       "",
				},
			},
			expected: &models.BidCaching{
				RuleID:        "966affd7-d087-57a2-baff-55b926f4c32d",
				Publisher:     "",
				Domain:        "",
				BidCaching:    1,
				Country:       null.String{},
				Os:            null.String{},
				Device:        null.String{},
				PlacementType: null.String{},
				Browser:       null.String{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mod := tt.args.bidCaching.ToModel()
			assert.Equal(t, tt.expected, mod)
		})
	}
}

func TestCreateBidCachingMetadataGeneration(t *testing.T) {
	tests := []struct {
		name         string
		modBC        models.BidCachingSlice
		finalRules   []BidCachingRealtimeRecord
		expectedJSON string
	}{
		{
			name: "Sort By Correct Order",
			modBC: models.BidCachingSlice{
				{
					RuleID:     "",
					Publisher:  "20814",
					Domain:     "stream-together.org",
					Device:     null.StringFrom("mobile"),
					BidCaching: 12,
				},
				{
					RuleID:     "",
					Publisher:  "20814",
					Domain:     "stream-together.org",
					Device:     null.StringFrom("mobile"),
					Country:    null.StringFrom("il"),
					BidCaching: 11,
				},
				{
					RuleID:     "",
					Publisher:  "20814",
					Domain:     "stream-together.org",
					Device:     null.StringFrom("mobile"),
					Country:    null.StringFrom("us"),
					BidCaching: 14,
				},
			},
			finalRules:   []BidCachingRealtimeRecord{},
			expectedJSON: `{"rules":[{"rule":"(p=20814__d=stream-together.org__c=il__os=.*__dt=mobile__pt=.*__b=.*)","bid_caching":11,"rule_id":"cc11f229-1d4a-5bd2-a6d0-5fae8c7a9bf4"},{"rule":"(p=20814__d=stream-together.org__c=us__os=.*__dt=mobile__pt=.*__b=.*)","bid_caching":14,"rule_id":"a0d406cd-bf98-50ab-9ff2-1b314b27da65"},{"rule":"(p=20814__d=stream-together.org__c=.*__os=.*__dt=mobile__pt=.*__b=.*)","bid_caching":12,"rule_id":"cb45cb97-5ca2-503d-9008-317dbbe26d10"}]}`,
		},
		{
			name: "Device with null value",
			modBC: models.BidCachingSlice{
				{
					RuleID:     "",
					Publisher:  "20814",
					Domain:     "stream-together.org",
					Country:    null.StringFrom("us"),
					BidCaching: 11,
				},
			},
			finalRules:   []BidCachingRealtimeRecord{},
			expectedJSON: `{"rules": [{"rule": "(p=20814__d=stream-together.org__c=us__os=.*__dt=.*__pt=.*__b=.*)", "bid_caching": 11, "rule_id": "ad18394a-ee20-58c2-bb9b-dd459550a9f7"}]}`,
		},
		{
			name: "Same ruleId different input bid_caching",
			modBC: models.BidCachingSlice{
				{
					RuleID:     "",
					Publisher:  "20814",
					Domain:     "stream-together.org",
					Country:    null.StringFrom("us"),
					Device:     null.StringFrom("mobile"),
					BidCaching: 14,
				},
			},
			finalRules:   []BidCachingRealtimeRecord{},
			expectedJSON: `{"rules": [{"rule": "(p=20814__d=stream-together.org__c=us__os=.*__dt=mobile__pt=.*__b=.*)", "bid_caching": 14, "rule_id": "a0d406cd-bf98-50ab-9ff2-1b314b27da65"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateBidCachingMetadata(tt.modBC, tt.finalRules)

			resultJSON, err := json.Marshal(map[string]interface{}{"rules": result})
			if err != nil {
				t.Fatalf("Failed to marshal result to JSON: %v", err)
			}

			assert.JSONEq(t, tt.expectedJSON, string(resultJSON))
		})
	}
}