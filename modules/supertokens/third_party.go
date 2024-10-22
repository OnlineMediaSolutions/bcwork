package supertokens

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/spf13/viper"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	MaxDaysForTemporaryPassword = 30

	supertokensThirdPartyRootKeyConfig = "thirdParty"
	// google
	supertokensGoogleRootKeyConfig         = supertokensThirdPartyRootKeyConfig + "." + "google"
	supertokensGoogleClientIDKeyConfig     = supertokensGoogleRootKeyConfig + "." + "clientID"
	supertokensGoogleClientSecretKeyConfig = supertokensGoogleRootKeyConfig + "." + "clientSecret"
	// apple
	supertokensAppleRootKeyConfig       = supertokensThirdPartyRootKeyConfig + "." + "apple"
	supertokensAppleClientIDKeyConfig   = supertokensAppleRootKeyConfig + "." + "clientID"
	supertokensAppleKeyIDKeyConfig      = supertokensAppleRootKeyConfig + "." + "keyID"
	supertokensApplePrivateKeyKeyConfig = supertokensAppleRootKeyConfig + "." + "privateKey"
	supertokensAppleTeamIDKeyConfig     = supertokensAppleRootKeyConfig + "." + "teamID"
)

var (
	errNotAllowed                        = errors.New("provided email not allowed to sign in/up using third-party providers")
	errUserDisabled                      = errors.New("user disabled")
	errTemporaryPasswordNeedsToBeChanged = errors.New("temporary password needs to be changed")
)

func GetThirdPartyEmailPasswordFunctionsOverride() *tpepmodels.OverrideStruct {
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

func getThirdPartyProviderGoogle() tpmodels.ProviderInput {
	supertokensEnv := viper.GetString(supertokensRootKeyConfig + "." + supertokensEnvKeyConfig)

	return tpmodels.ProviderInput{
		Config: tpmodels.ProviderConfig{
			ThirdPartyId: "google",
			Clients: []tpmodels.ProviderClientConfig{
				{
					ClientID:     supertokensConfig(supertokensEnv, supertokensGoogleClientIDKeyConfig),
					ClientSecret: supertokensConfig(supertokensEnv, supertokensGoogleClientSecretKeyConfig),
				},
			},
		},
	}
}

func getThirdPartyProviderApple() tpmodels.ProviderInput {
	supertokensEnv := viper.GetString(supertokensRootKeyConfig + "." + supertokensEnvKeyConfig)

	return tpmodels.ProviderInput{
		Config: tpmodels.ProviderConfig{
			ThirdPartyId: "apple",
			Clients: []tpmodels.ProviderClientConfig{
				{
					ClientID: supertokensConfig(supertokensEnv, supertokensAppleClientIDKeyConfig),
					AdditionalConfig: map[string]interface{}{
						"keyId":      supertokensConfig(supertokensEnv, supertokensAppleKeyIDKeyConfig),
						"privateKey": supertokensConfig(supertokensEnv, supertokensApplePrivateKeyKeyConfig),
						"teamId":     supertokensConfig(supertokensEnv, supertokensAppleTeamIDKeyConfig),
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
