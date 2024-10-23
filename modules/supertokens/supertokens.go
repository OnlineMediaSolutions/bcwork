package supertokens

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	httpclient "github.com/m6yf/bcwork/modules/http_client"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const (
	SuperTokensIDKey         = "id"
	SuperTokensEmailKey      = "email"
	SuperTokensTimeJoinedKey = "timeJoined"

	SuperTokensMetaDataFirstNameKey        = "first_name"
	SuperTokensMetaDataLastNameKey         = "last_name"
	SuperTokensMetaDataOrganizationNameKey = "organization_name"
	SuperTokensMetaDataAddressKey          = "address"
	SuperTokensMetaDataPhoneKey            = "phone"
	SuperTokensMetaDataEnabledKey          = "enabled"
	SuperTokensMetaDataDisabledAtKey       = "disabled_at"

	UserEmailContextKey = "email"
	RoleContextKey      = "role"

	UserRoleName  = "user"
	AdminRoleName = "admin"

	CreateUserSupertokenPath     = "/signup"
	ChangePasswordSupertokenPath = "/forgot-password"
)

type TokenManagementSystem interface {
	CreateUser(ctx context.Context, email, password string) (string, error)
	RevokeAllSessionsForUser(userID string) error
	GetWebURL() string
	VerifySession(next http.Handler) http.Handler
	AdminRoleRequired(c *fiber.Ctx) error
}

type SuperTokensClient struct {
	apiURL                              string
	webURL                              string
	skipSessionVerificationForLocalHost bool // for local development and workers
	httpClient                          httpclient.Doer
}

var _ TokenManagementSystem = (*SuperTokensClient)(nil)

func NewSuperTokensClient(
	apiURL string,
	webURL string,
	initFunc func() error,
	skipSessionVerificationForLocalHost bool,
) (*SuperTokensClient, error) {
	err := initFunc()
	if err != nil {
		return nil, fmt.Errorf("failed to init supertokens: %w", err)
	}

	return &SuperTokensClient{
		apiURL:                              apiURL,
		webURL:                              webURL,
		skipSessionVerificationForLocalHost: skipSessionVerificationForLocalHost,
		httpClient:                          httpclient.New(),
	}, nil
}

func (c *SuperTokensClient) GetWebURL() string {
	return c.webURL
}

func (c *SuperTokensClient) CreateUser(ctx context.Context, email, password string) (string, error) {
	payload := fmt.Sprintf(`{"formFields": [{"id": "email","value": "%v"},{"id": "password","value": "%v"}]}`, email, password)
	url := c.apiURL + CreateUserSupertokenPath

	body, err := c.httpClient.Do(ctx, http.MethodPost, url, payload)
	if err != nil {
		return "", fmt.Errorf("can't do request to supertokens API: %w", err)
	}

	var resp CreateUserResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", fmt.Errorf("can't unmarshal [%v] to createUserResponse: %w", string(body), err)
	}

	if resp.Status != "OK" {
		return "", fmt.Errorf("error creating user in supertoken: status [%v] not equal 'OK'", resp.Status)
	}

	return resp.User.ID, nil
}

func (c *SuperTokensClient) RevokeAllSessionsForUser(email string) error {
	tenantID := supertokens.DefaultTenantId

	users, err := thirdpartyemailpassword.GetUsersByEmail(tenantID, email)
	if err != nil {
		return fmt.Errorf("error revoking all sessions for user: %w", err)
	}

	for _, user := range users {
		_, err = session.RevokeAllSessionsForUser(user.ID, &tenantID)
		if err != nil {
			return fmt.Errorf("error revoking all sessions for user: %w", err)
		}
	}

	return nil
}

type CreateUserResponse struct {
	Status string `json:"status"`
	User   struct {
		ID         string `json:"id"`
		Email      string `json:"email"`
		TimeJoined int    `json:"timeJoined"`
	} `json:"user"`
}
