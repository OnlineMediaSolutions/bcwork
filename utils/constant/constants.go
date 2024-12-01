package constant

const (
	MinFactorValue          = 0.01
	MaxFactorValue          = 10.0
	MaxCountryCodeLength    = 2
	MinDPOFactorValue       = 0
	MaxDPOFactorValue       = 100
	MinBidCachingValue      = 1
	MaxRefreshCacheValue    = 500
	MinRefreshCacheValue    = 1
	ProductionApiUrl        = "http://localhost:8000"
	DpoGetEndpoint          = "/dpo/get"
	DpoSetEndpoint          = "/bulk/dpo"
	ConfigEndpoint          = "/config/get"
	PostgresTimestampLayout = "2006-01-02 15:04:05"
	PostgresTimestamp       = "2006-01-02"

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

	// Workers
	SellersJsonWorkerCount  = 5
	AdsTxtNotVerifiedStatus = "not verified"
	AdsTxtNotIncludedStatus = "not included"
	AdsTxtIncludedStatus    = "included"
)
