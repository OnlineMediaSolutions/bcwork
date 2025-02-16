package constant

type ContextKey string

const (
	MinFactorValue          = 0.01
	MaxFactorValue          = 10.0
	MaxCountryCodeLength    = 2
	MinDPOFactorValue       = 0
	MaxDPOFactorValue       = 100
	MinBidCachingValue      = 1
	MaxRefreshCacheValue    = 500
	RefreshCacheDeleteValue = 500
	MinRefreshCacheValue    = 0
	MinThreshold            = 0.000
	MaxThreshold            = 0.010
	PostgresTimestampLayout = "2006-01-02 15:04:05"
	PostgresTimestamp       = "2006-01-02"

	CurrentTime = "NOW()"

	// Context
	UserIDContextKey      ContextKey = "user_id"
	UserEmailContextKey   ContextKey = "email"
	RoleContextKey        ContextKey = "role"
	RequestIDContextKey   ContextKey = "request_id"
	LoggerContextKey      ContextKey = "logger"
	RequestPathContextKey ContextKey = "request_path"

	// Global Factor Fee Type
	GlobalFactorConsultantFeeType = "consultant_fee"
	GlobalFactorTechFeeType       = "tech_fee"
	GlobalFactorTAMFeeType        = "tam_fee"

	// Workers
	SellersJsonWorkerCount       = 5
	AdsTxtNotVerifiedStatus      = "not verified"
	AdsTxtNotIncludedStatus      = "not included"
	AdsTxtIncludedStatus         = "included"
	AdsTxtRequestTimeout         = 60
	NewBidderAutomationThreshold = 5000
	ConversionToMillion          = 1000000

	// Request Types
	RequestTypeJS = "js"
)
