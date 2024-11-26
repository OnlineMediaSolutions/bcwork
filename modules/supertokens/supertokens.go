package supertokens

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/config"
	httpclient "github.com/m6yf/bcwork/modules/http_client"
	"github.com/spf13/viper"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const (
	DeveloperRoleName   = "Developer"
	AdminRoleName       = "Admin"
	SupermemberRoleName = "Supermember"
	MemberRoleName      = "Member"
	PublisherRoleName   = "Publisher"
	ConsultantRoleName  = "Consultant"

	WorkerUserID   = -1
	WorkerUserName = "Internal Worker"

	AutomationUserID   = -2
	AutomationUserName = "Automation"

	CreateUserSupertokenPath     = "/signup"
	ChangePasswordSupertokenPath = "/forgot-password"
)

type TokenManagementSystem interface {
	CreateUser(ctx context.Context, email, password string) (string, error)
	RevokeAllSessionsForUser(userID string) error
	GetWebURL() string
	VerifySession(next http.Handler) http.Handler
	AdminRoleRequired(c *fiber.Ctx) error
	GetEmailByUserID(userID string) (string, error)
}

type SuperTokensClient struct {
	apiURL     string
	webURL     string
	apiKeys    []string
	httpClient httpclient.Doer
}

var _ TokenManagementSystem = (*SuperTokensClient)(nil)

func NewSuperTokensClient(
	apiURL string,
	webURL string,
	initFunc func() error,
) (*SuperTokensClient, error) {
	err := initFunc()
	if err != nil {
		return nil, fmt.Errorf("failed to init supertokens: %w", err)
	}

	awsKey := viper.GetString(config.AWSWorkerAPIKeyKey)
	cronKey := viper.GetString(config.CronWorkerAPIKeyKey)

	return &SuperTokensClient{
		apiURL:     apiURL,
		webURL:     webURL,
		apiKeys:    []string{awsKey, cronKey},
		httpClient: httpclient.New(true),
	}, nil
}

func (c *SuperTokensClient) GetWebURL() string {
	return c.webURL
}

func (c *SuperTokensClient) CreateUser(ctx context.Context, email, password string) (string, error) {
	payload := fmt.Sprintf(`{"formFields": [{"id": "email","value": "%v"},{"id": "password","value": "%v"}]}`, email, password)
	url := c.apiURL + CreateUserSupertokenPath

	body, _, err := c.httpClient.Do(ctx, http.MethodPost, url, strings.NewReader(payload))
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

func (c *SuperTokensClient) GetEmailByUserID(userID string) (string, error) {
	user, err := thirdpartyemailpassword.GetUserById(userID)
	if err != nil {
		return "", fmt.Errorf("can't get user by id from supertokens: %w", err)
	}

	if user == nil {
		return "", errors.New("user not found in supertokens")
	}

	return user.Email, nil
}

type CreateUserResponse struct {
	Status string `json:"status"`
	User   struct {
		ID         string `json:"id"`
		Email      string `json:"email"`
		TimeJoined int    `json:"timeJoined"`
	} `json:"user"`
}
