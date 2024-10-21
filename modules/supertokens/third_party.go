package supertokens

import (
	"context"
	"database/sql"
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var (
	errNotAllowed   = errors.New("provided email not allowed to sign in/up using third-party providers")
	errUserDisabled = errors.New("user disabled")
)

func getThirdPartyEmailPasswordFunctionsOverride() *tpepmodels.OverrideStruct {
	return &tpepmodels.OverrideStruct{
		Functions: func(originalImplementation tpepmodels.RecipeInterface) tpepmodels.RecipeInterface {
			// sign in/up for third party providers
			originalThirdPartySignInUp := *originalImplementation.ThirdPartySignInUp

			(*originalImplementation.ThirdPartySignInUp) = func(
				thirdPartyID string,
				thirdPartyUserID string,
				email string,
				oAuthTokens tpmodels.TypeOAuthTokens,
				rawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider,
				tenantId string,
				userContext supertokens.UserContext,
			) (tpepmodels.SignInUpResponse, error) {
				err := validateUser(email)
				if err != nil {
					return tpepmodels.SignInUpResponse{}, err
				}

				return originalThirdPartySignInUp(thirdPartyID, thirdPartyUserID, email, oAuthTokens, rawUserInfoFromProvider, tenantId, userContext)
			}

			// sign in for email/password
			originalEmailPasswordSignIn := *originalImplementation.EmailPasswordSignIn

			(*originalImplementation.EmailPasswordSignIn) = func(
				email string,
				password string,
				tenantId string,
				userContext supertokens.UserContext,
			) (tpepmodels.SignInResponse, error) {
				err := validateUser(email)
				if err != nil {
					return tpepmodels.SignInResponse{}, err
				}

				return originalEmailPasswordSignIn(email, password, tenantId, userContext)
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

func validateUser(email string) error {
	user, err := getUserByEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errNotAllowed
		}
		return err
	}

	if !user.Enabled {
		return errUserDisabled
	}

	return nil
}
