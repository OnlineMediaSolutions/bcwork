package validations

import (
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/m6yf/bcwork/dto"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
	"github.com/m6yf/bcwork/utils/constant"
)

const (
	// Keys
	targetingPriceModelValidationKey = "targetingPriceModel"
	targetingStatusValidationKey     = "targetingStatus"
	countriesValidationKey           = "countries"
	devicesValidationKey             = "devices"
	emailValidationKey               = "email"
	phoneValidationKey               = "phone"
	userRoleValidationKey            = "userRole"
	userTypesValidationKey           = "userTypes"
	approvalProcessKey               = "approvalProcess"
	dpBlocksKey                      = "dpBlocks"
	dpThresholdKey                   = "dpThreshold"
	bidCachingControlPercentageKey   = "bccp"
	mediaTypeValidationKey           = "mediaType"
	intergrationTypeValidationKey    = "integrationType"
	adsTxtDomainStatusValidationKey  = "adsTxtDomainStatus"
	adsTxtDemandStatusValidationKey  = "adsTxtDemandStatus"
	// Error messages
	countryValidationErrorMessage            = "country code must be 2 characters long and should be in the allowed list"
	deviceValidationErrorMessage             = "device should be in the allowed list"
	targetingCostModelValidationErrorMessage = "targeting price model should be 'CPM' or 'Rev Share'"
	targetingStatusValidationErrorMessage    = "targeting status should be 'Active', 'Paused' or 'Archived'"
	emailValidationErrorMessage              = "email not valid"
	phoneValidationErrorMessage              = "phone not valid"
	userRoleValidationErrorMessage           = "user role must be in allowed list"
	userTypesValidationErrorMessage          = "user types must be in allowed list"
	approvalProcessErrorMessage              = "approval process must be in allowed list"
	dpBlocksErrorMessage                     = "dp blocks must be in allowed list"
	bidCachingControlPercentageErrorMessage  = "bid caching control percentage must be from 0 to 1"
	mediaTypeErrorMessage                    = "media type must be in allowed list"
	intergrationTypeErrorMessage             = "integration type must be in allowed list"
	adsTxtDomainStatusErrorMessage           = "ads.txt domain status must be in allowed list"
	adsTxtDemandStatusErrorMessage           = "ads.txt demand status must be in allowed list"
)

var (
	Validator = validator.New()

	globalFactorKeyTypes = []string{"tech_fee", "consultant_fee", "tam_fee"}
	targetingCostModels  = []string{dto.TargetingPriceModelCPM, dto.TargetingPriceModelRevShare}
	targetingStatuses    = []string{dto.TargetingStatusActive, dto.TargetingStatusPaused, dto.TargetingStatusArchived}
	userRoles            = []string{
		supertokens_module.DeveloperRoleName, supertokens_module.AdminRoleName,
		supertokens_module.SupermemberRoleName, supertokens_module.MemberRoleName,
		supertokens_module.PublisherRoleName, supertokens_module.ConsultantRoleName,
	}
	userTypes = []string{
		dto.UserTypeAccountManager, dto.UserTypeCampaignManager, dto.UserTypeMediaBuyer,
	}
	approvalProcesses = []string{
		dto.EmailApprovalProcess, dto.DemandPartnerPlatformApprovalProcess,
		dto.GDocApprovalProcess, dto.OtherApprovalProcess,
	}
	dpBlocks = []string{
		dto.EmailApprovalProcess, dto.DemandPartnerPlatformApprovalProcess,
		dto.GDocApprovalProcess, dto.OtherApprovalProcess,
	}

	phoneClearRegExp = regexp.MustCompile(`[ ()-]*`)
	phoneFindRegExp  = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

	integrationTypes = []string{
		dto.ORTBIntergrationType, dto.PrebidServerIntergrationType, dto.AmazonAPSIntergrationType,
	}
	mediaTypes = []string{
		dto.WebBannersMediaType, dto.VideoMediaType, dto.InAppMediaType,
	}

	adsTxtDomainStatuses = []string{
		dto.DomainStatusActive, dto.DomainStatusNew, dto.DomainStatusPaused,
	}
	adsTxtDemandStatuses = []string{
		dto.DPStatusPending, dto.DPStatusApproved, dto.DPStatusApprovedPaused,
		dto.DPStatusRejected, dto.DPStatusRejectedTQ, dto.DPStatusDisabledSPO,
		dto.DPStatusDisabledNoImps, dto.DPStatusHighDiscrepancy, dto.DPStatusNotSent,
		dto.DPStatusNoForm, dto.DPStatusWillNotBeSent,
	}
)

func init() {
	err := Validator.RegisterValidation("country", countryValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("device", deviceValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("placement_type", placementTypeValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("os", osValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("browser", browserValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("floor", floorValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("factor", factorValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("factorDpo", factorDpoValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("rate", rateValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("all", notAllowedTheWordAllValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("url", validateURL)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("globalFactorKey", globalFactorKeyValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(targetingPriceModelValidationKey, targetingCostModelValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(targetingStatusValidationKey, targetingStatusValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(countriesValidationKey, countriesValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(devicesValidationKey, devicesValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(emailValidationKey, emailValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(phoneValidationKey, phoneValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(userRoleValidationKey, userRoleValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(userTypesValidationKey, userTypesValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("bid_caching", bidCachingValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("refresh_cache", refreshCacheValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(approvalProcessKey, approvalProcessValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(dpBlocksKey, dpBlocksValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(dpThresholdKey, dpThresholdValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(bidCachingControlPercentageKey, bidCachingControlPercentageValidation, true)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(intergrationTypeValidationKey, integrationTypeValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(mediaTypeValidationKey, mediaTypeValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(adsTxtDemandStatusValidationKey, adsTxtDemandStatusValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(adsTxtDomainStatusValidationKey, adsTxtDomainStatusValidation)
	if err != nil {
		return
	}
}

func floorValidation(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	return val >= 0
}

func factorValidation(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	return val >= constant.MinFactorValue && val <= constant.MaxFactorValue
}

func factorDpoValidation(fieldLevel validator.FieldLevel) bool {
	val := fieldLevel.Field().Float()
	return val >= constant.MinDPOFactorValue && val <= constant.MaxDPOFactorValue
}

func rateValidation(fieldLevel validator.FieldLevel) bool {
	val := fieldLevel.Field().Float()
	return val >= constant.MinDPOFactorValue && val <= constant.MaxDPOFactorValue
}

func countryValidation(fl validator.FieldLevel) bool {
	country := fl.Field().String()
	return validateCountry(country)
}

func countriesValidation(fl validator.FieldLevel) bool {
	countries, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}

	for _, country := range countries {
		isValid := validateCountry(country)
		if !isValid {
			return false
		}
	}

	return true
}

func validateCountry(country string) bool {
	if country == "" || len(country) == 0 {
		return true
	}
	if len(country) != constant.MaxCountryCodeLength {
		return false
	}
	if _, ok := constant.AllowedCountries[country]; !ok {
		return false
	}
	return true
}

func placementTypeValidation(fl validator.FieldLevel) bool {
	placementType := fl.Field().String()
	if len(placementType) == 0 || placementType == "" {
		return true
	}

	if _, ok := constant.AllowedPlacemenTypes[placementType]; !ok {
		return false
	}
	return true
}

func osValidation(fl validator.FieldLevel) bool {
	os := fl.Field().String()
	if len(os) == 0 || os == "" {
		return true
	}

	if _, ok := constant.AllowedOses[os]; !ok {
		return false
	}
	return true
}

func browserValidation(fl validator.FieldLevel) bool {
	browser := fl.Field().String()
	if len(browser) == 0 || browser == "" {
		return true
	}

	if _, ok := constant.AllowedBrowsers[browser]; !ok {
		return false
	}
	return true
}

func deviceValidation(fl validator.FieldLevel) bool {
	device := fl.Field().String()
	return validateDevice(device)
}

func devicesValidation(fl validator.FieldLevel) bool {
	devices, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}

	for _, device := range devices {
		isValid := validateDevice(device)
		if !isValid {
			return false
		}
	}

	return true
}

func validateDevice(device string) bool {
	if device == "" {
		return true
	}
	if _, ok := constant.AllowedDevices[device]; !ok {
		return false
	}
	return true
}

func notAllowedTheWordAllValidation(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if field == "all" {
		return false
	}
	return true
}

func validateURL(fl validator.FieldLevel) bool {
	u, err := url.ParseRequestURI(fl.Field().String())
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

func globalFactorKeyValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	return slices.Contains(globalFactorKeyTypes, field.String())
}

func targetingCostModelValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	return slices.Contains(targetingCostModels, field.String())
}

func targetingStatusValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	return slices.Contains(targetingStatuses, field.String())
}

func emailValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	_, err := mail.ParseAddress(field.String())
	return err == nil
}

func phoneValidation(fl validator.FieldLevel) bool {
	field := phoneClearRegExp.ReplaceAllString(fl.Field().String(), "")
	return phoneFindRegExp.FindString(field) != "" || fl.Field().String() == ""
}

func userRoleValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	return slices.Contains(userRoles, field.String())
}

func userTypesValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.IsNil() {
		return true
	}

	providedUserTypes := field.Interface().([]string)
	for _, userType := range providedUserTypes {
		if !slices.Contains(userTypes, userType) {
			return false
		}
	}

	return true
}

func bidCachingValidation(fl validator.FieldLevel) bool {
	val := fl.Field().Int()
	return val >= constant.MinBidCachingValue
}

func refreshCacheValidation(fl validator.FieldLevel) bool {
	val := fl.Field().Int()
	return val <= constant.MaxRefreshCacheValue && val >= constant.MinRefreshCacheValue
}

func approvalProcessValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	return slices.Contains(approvalProcesses, field.String())
}

func dpBlocksValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	return slices.Contains(dpBlocks, field.String())
}

func dpThresholdValidation(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	return val >= constant.MinThreshold && val <= constant.MaxThreshold
}

func bidCachingControlPercentageValidation(fl validator.FieldLevel) bool {
	var val float64
	if fl.Field().Kind() == reflect.Ptr {
		if fl.Field().IsNil() {
			return true
		} else {
			val = fl.Field().Elem().Float()
		}
	} else {
		val = fl.Field().Float()
	}

	return val >= dto.BidCachingControlPercentageMin && val <= dto.BidCachingControlPercentageMax
}

func integrationTypeValidation(fl validator.FieldLevel) bool {
	field, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}

	if len(field) == 0 {
		switch getStructName(fl) {
		case dto.UpdatePublisherValuesStructName:
			return true
		default:
			return false
		}
	}

	for _, v := range field {
		if !slices.Contains(integrationTypes, v) {
			return false
		}
	}

	return true
}

func mediaTypeValidation(fl validator.FieldLevel) bool {
	field, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}

	if len(field) == 0 {
		switch getStructName(fl) {
		case dto.UpdatePublisherValuesStructName:
			return true
		default:
			return false
		}
	}

	for _, v := range field {
		if !slices.Contains(mediaTypes, v) {
			return false
		}
	}

	return true
}

func adsTxtDemandStatusValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	return slices.Contains(adsTxtDemandStatuses, field.String())
}

func adsTxtDomainStatusValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	return slices.Contains(adsTxtDomainStatuses, field.String())
}

func getStructName(fl validator.FieldLevel) string {
	return fl.Parent().Type().Name()
}
