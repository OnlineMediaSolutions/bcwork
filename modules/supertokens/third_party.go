package supertokens

import (
	"errors"
	"log"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const (
	googleDomain = "google.com"
	omsDomain    = "onlinemediasolutions.com"
)

var (
	errNotAllowedDomain = errors.New("provided domain not allowed to sign in/up")
)

func isAllowedDomain(email string) bool {
	log.Print(email)
	emailArr := strings.Split(email, "@")
	// if emailArr[1] != omsDomain {
	return emailArr[1] != googleDomain
}

func getThirdPartySignInUpFunctionOverride() *tpmodels.OverrideStruct {
	return &tpmodels.OverrideStruct{
		Functions: func(originalImplementation tpmodels.RecipeInterface) tpmodels.RecipeInterface {
			originalThirdPartySignInUp := *originalImplementation.SignInUp

			(*originalImplementation.SignInUp) = func(thirdPartyID, thirdPartyUserID, email string, oAuthTokens map[string]interface{}, rawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider, tenantId string, userContext supertokens.UserContext) (tpmodels.SignInUpResponse, error) {
				if !isAllowedDomain(email) {
					return tpmodels.SignInUpResponse{}, errNotAllowedDomain
				}
				// We allow the sign in / up operation
				return originalThirdPartySignInUp(thirdPartyID, thirdPartyUserID, email, oAuthTokens, rawUserInfoFromProvider, tenantId, userContext)
			}

			return originalImplementation
		},

		APIs: func(originalImplementation tpmodels.APIInterface) tpmodels.APIInterface {
			originalSignInUpPOST := *originalImplementation.SignInUpPOST

			(*originalImplementation.SignInUpPOST) = func(provider *tpmodels.TypeProvider, input tpmodels.TypeSignInUpInput, tenantId string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.SignInUpPOSTResponse, error) {

				resp, err := originalSignInUpPOST(provider, input, tenantId, options, userContext)

				if errors.Is(err, errNotAllowedDomain) {
					// this error was thrown from our function override above.
					// so we send a useful message to the user
					return tpmodels.SignInUpPOSTResponse{
						GeneralError: &supertokens.GeneralErrorResponse{
							Message: "Sign ups are disabled. Please contact the admin.",
						},
					}, nil
				}

				return resp, err
			}

			return originalImplementation
		},
	}
}

/*
   We use different credentials for different platforms when required. For example the redirect URI for Github
   is different for Web and mobile. In such a case we can provide multiple providers with different client Ids.

   When the frontend makes a request and wants to use a specific clientId, it needs to send the clientId to use in the
   request. In the absence of a clientId in the request the SDK uses the default provider, indicated by `isDefault: true`.
   When adding multiple providers for the same type (Google, Github etc), make sure to set `isDefault: true`.
*/

func getThirdPartyProviderGoogle() tpmodels.ProviderInput {
	// TODO: replace credentials with config
	return tpmodels.ProviderInput{
		Config: tpmodels.ProviderConfig{
			ThirdPartyId: "google",
			Clients: []tpmodels.ProviderClientConfig{
				{
					ClientID:     "1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com",
					ClientSecret: "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",
				},
			},
		},
	}
}

func getThirdPartyProviderGithub() tpmodels.ProviderInput {
	// TODO: replace credentials with config
	return tpmodels.ProviderInput{
		Config: tpmodels.ProviderConfig{
			ThirdPartyId: "github",
			Clients: []tpmodels.ProviderClientConfig{
				{
					ClientID:     "467101b197249757c71f",
					ClientSecret: "e97051221f4b6426e8fe8d51486396703012f5bd",
				},
			},
		},
	}
}

func getThirdPartyProviderApple() tpmodels.ProviderInput {
	// TODO: replace credentials with config
	return tpmodels.ProviderInput{
		Config: tpmodels.ProviderConfig{
			ThirdPartyId: "apple",
			Clients: []tpmodels.ProviderClientConfig{
				{
					ClientID: "4398792-io.supertokens.example.service",
					AdditionalConfig: map[string]interface{}{
						"keyId":      "7M48Y4RYDL",
						"privateKey": "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgu8gXs+XYkqXD6Ala9Sf/iJXzhbwcoG5dMh1OonpdJUmgCgYIKoZIzj0DAQehRANCAASfrvlFbFCYqn3I2zeknYXLwtH30JuOKestDbSfZYxZNMqhF/OzdZFTV0zc5u5s3eN+oCWbnvl0hM+9IW0UlkdA\n-----END PRIVATE KEY-----",
						"teamId":     "YWQCXGJRJL",
					},
				},
			},
		},
	}
}
