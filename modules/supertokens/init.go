package supertokens

import (
	"fmt"

	"github.com/m6yf/bcwork/utils/pointer"
	"github.com/spf13/viper"
	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const (
	supertokensRootKeyConfig = "supertokens"
	supertokensEnvKeyConfig  = "env"
	// connection info
	supertokensURIKeyConfig    = "uri"
	supertokensAPIKeyKeyConfig = "key"
	// application info
	supertokensAppInfoRootKeyConfig     = "appInfo"
	supertokensAppNameKeyConfig         = supertokensAppInfoRootKeyConfig + "." + "appName"
	supertokensAPIDomainKeyConfig       = supertokensAppInfoRootKeyConfig + "." + "apiDomain"
	supertokensAPIBasePathKeyConfig     = supertokensAppInfoRootKeyConfig + "." + "apiBasePath"
	supertokensWebsiteDomainKeyConfig   = supertokensAppInfoRootKeyConfig + "." + "websiteDomain"
	supertokensWebsiteBasePathKeyConfig = supertokensAppInfoRootKeyConfig + "." + "websiteBasePath"
	// dashboard
	supertokensDashboardApiKeyKeyConfig = "dashboardApiKey"
)

func GetSuperTokensConfig() (string, string, func() error) {
	supertokensEnv := viper.GetString(supertokensRootKeyConfig + "." + supertokensEnvKeyConfig)

	apiURL := "http://localhost:" + viper.GetString("ports.http") + supertokensConfig(supertokensEnv, supertokensAPIBasePathKeyConfig)
	webURL := supertokensConfig(supertokensEnv, supertokensWebsiteDomainKeyConfig) + supertokensConfig(supertokensEnv, supertokensWebsiteBasePathKeyConfig)

	return apiURL, webURL, initSuperTokens
}

func initSuperTokens() error {
	antiCsrf := "NONE" // Should be one of "NONE" or "VIA_CUSTOM_HEADER" or "VIA_TOKEN"
	supertokensEnv := viper.GetString(supertokensRootKeyConfig + "." + supertokensEnvKeyConfig)

	return supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: supertokensConfig(supertokensEnv, supertokensURIKeyConfig),
			APIKey:        supertokensConfig(supertokensEnv, supertokensAPIKeyKeyConfig),
		},
		AppInfo: supertokens.AppInfo{
			AppName:         supertokensConfig(supertokensEnv, supertokensAppNameKeyConfig),
			APIDomain:       supertokensConfig(supertokensEnv, supertokensAPIDomainKeyConfig),
			APIBasePath:     pointer.String(supertokensConfig(supertokensEnv, supertokensAPIBasePathKeyConfig)),
			WebsiteDomain:   supertokensConfig(supertokensEnv, supertokensWebsiteDomainKeyConfig),
			WebsiteBasePath: pointer.String(supertokensConfig(supertokensEnv, supertokensWebsiteBasePathKeyConfig)),
		},
		RecipeList: []supertokens.Recipe{
			thirdpartyemailpassword.Init(&tpepmodels.TypeInput{
				Override: GetThirdPartyEmailPasswordFunctionsOverride(),
				Providers: []tpmodels.ProviderInput{
					getThirdPartyProviderGoogle(),
					getThirdPartyProviderApple(),
				},
			}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &antiCsrf,
			}),
			dashboard.Init(&dashboardmodels.TypeInput{
				ApiKey: supertokensConfig(supertokensEnv, supertokensDashboardApiKeyKeyConfig),
			}),
		},
	})
}

func supertokensConfig(env string, key string) string {
	return viper.GetString(fmt.Sprintf("%s.%s.%s", supertokensRootKeyConfig, env, key))
}
