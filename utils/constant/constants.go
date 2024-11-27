package constant

const (
	MinFactorValue          = 0.01
	MaxFactorValue          = 10.0
	MaxCountryCodeLength    = 2
	MinDPOFactorValue       = 0
	MaxDPOFactorValue       = 100
	MinBidCachingValue      = 1
	MaxRefreshCacheValue    = 500
	ProductionApiUrl        = "http://localhost:8000"
	DpoGetEndpoint          = "/dpo/get"
	DpoSetEndpoint          = "/bulk/dpo"
	PostgresTimestampLayout = "2006-01-02 15:04:05"

	// Context
	UserIDContextKey      = "user_id"
	UserEmailContextKey   = "email"
	RoleContextKey        = "role"
	RequestIDContextKey   = "request_id"
	LoggerContextKey      = "logger"
	RequestPathContextKey = "request_path"

	// Global Factor Fee Type
	GlobalFactorConsultantFeeType = "consultant_fee"
	GlobalFactorTechFeeType       = "tech_fee"
	GlobalFactorTAMFeeType        = "tam_fee"
)
