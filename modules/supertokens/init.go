package supertokens

import (
	"github.com/m6yf/bcwork/utils/pointer"
	"github.com/spf13/viper"
	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/supertokens/supertokens-golang/recipe/userroles"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func initSuperTokens() error {
	antiCsrf := "NONE" // Should be one of "NONE" or "VIA_CUSTOM_HEADER" or "VIA_TOKEN"

	// supertokensPrefix := viper.GetString("supertokens.env") // TODO: prefix depend on enviroment

	return supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: viper.GetString("supertokens.uri"),
			APIKey:        viper.GetString("supertokens.key"),
		},
		AppInfo: supertokens.AppInfo{
			AppName:         viper.GetString("supertokens.appInfo.appName"),
			APIDomain:       viper.GetString("supertokens.appInfo.apiDomain"),
			APIBasePath:     pointer.String(viper.GetString("supertokens.appInfo.apiBasePath")),
			WebsiteDomain:   viper.GetString("supertokens.appInfo.websiteDomain"),
			WebsiteBasePath: pointer.String(viper.GetString("supertokens.appInfo.websiteBasePath")),
		},
		RecipeList: []supertokens.Recipe{
			thirdparty.Init(&tpmodels.TypeInput{
				Override: getThirdPartySignInUpFunctionOverride(), // TODO: override didn't work
				// TODO: override if user not enabled
			}),
			thirdpartyemailpassword.Init(&tpepmodels.TypeInput{
				Providers: []tpmodels.ProviderInput{
					getThirdPartyProviderGoogle(),
					getThirdPartyProviderGithub(),
					getThirdPartyProviderApple(),
				},
			}),
			dashboard.Init(&dashboardmodels.TypeInput{
				ApiKey: viper.GetString("supertokens.dashboardApiKey"),
			}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &antiCsrf,
			}),
			userroles.Init(nil),
			usermetadata.Init(nil),
		},
	})
}
