package supertokens

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const MaxDaysForTemporaryPassword = 30

var (
	errNotAllowed                        = errors.New("provided email not allowed to sign in/up using third-party providers")
	errUserDisabled                      = errors.New("user disabled")
	errTemporaryPasswordNeedsToBeChanged = errors.New("temporary password needs to be changed")
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
				err := validateUserThirdParty(email)
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
				err := validateUserEmailPassword(email)
				if err != nil {
					return tpepmodels.SignInResponse{}, err
				}

				return originalEmailPasswordSignIn(email, password, tenantId, userContext)
			}

			// create password reset
			originalCreateResetPasswordToken := *originalImplementation.CreateResetPasswordToken

			(*originalImplementation.CreateResetPasswordToken) = func(
				userID string,
				tenantId string,
				userContext supertokens.UserContext,
			) (epmodels.CreateResetPasswordTokenResponse, error) {
				resp, err := originalCreateResetPasswordToken(userID, tenantId, userContext)
				if err != nil {
					return epmodels.CreateResetPasswordTokenResponse{}, err
				}

				if resp.OK != nil {
					err := updateResetToken(userID, resp.OK.Token)
					if err != nil {
						return epmodels.CreateResetPasswordTokenResponse{}, err
					}
				}

				return resp, err
			}

			// password reset
			originalResetPasswordUsingToken := *originalImplementation.ResetPasswordUsingToken

			(*originalImplementation.ResetPasswordUsingToken) = func(
				token string,
				newPassword string,
				tenantId string,
				userContext supertokens.UserContext,
			) (epmodels.ResetPasswordUsingTokenResponse, error) {
				err := updatePasswordChanging(token)
				if err != nil {
					return epmodels.ResetPasswordUsingTokenResponse{}, err
				}
				return originalResetPasswordUsingToken(token, newPassword, tenantId, userContext)
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

func validateUserThirdParty(email string) error {
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

func validateUserEmailPassword(email string) error {
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

	if isPasswordNeedsToBeChanged(user.PasswordChanged, user.CreatedAt) {
		return errTemporaryPasswordNeedsToBeChanged
	}

	return nil
}

func isPasswordNeedsToBeChanged(passwordChanged bool, createdAt time.Time) bool {

	return !passwordChanged && createdAt.AddDate(0, 0, MaxDaysForTemporaryPassword).Before(time.Now())
}

func updateResetToken(userID, resetToken string) error {
	ctx := context.Background()

	mod, err := models.Users(
		models.UserWhere.UserID.EQ(userID),
		models.UserWhere.PasswordChanged.EQ(false),
	).One(ctx, bcdb.DB())
	if err != nil {
		// if password was already changed, skip that user
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	mod.ResetToken = null.StringFrom(resetToken)

	_, err = mod.Update(ctx, bcdb.DB(), boil.Whitelist(models.UserColumns.ResetToken))
	if err != nil {
		return err
	}

	return nil
}

func updatePasswordChanging(resetToken string) error {
	ctx := context.Background()
	modResetToken := null.StringFrom(resetToken)

	mod, err := models.Users(models.UserWhere.ResetToken.EQ(modResetToken)).One(ctx, bcdb.DB())
	if err != nil {
		return err
	}

	mod.ResetToken = null.NewString("", false)
	mod.PasswordChanged = true

	_, err = mod.Update(ctx, bcdb.DB(), boil.Whitelist(models.UserColumns.ResetToken, models.UserColumns.PasswordChanged))
	if err != nil {
		return err
	}

	return nil
}
