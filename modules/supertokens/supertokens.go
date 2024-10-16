package supertokens

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	httpclient "github.com/m6yf/bcwork/modules/http_client"
	"github.com/spf13/viper"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/supertokens/supertokens-golang/recipe/userroles"
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

	UserIDContextKey   = "user_id"
	TenantIDContextKey = "tenant_id"
	RoleContextKey     = "role"

	UserRoleName  = "user"
	AdminRoleName = "admin"
)

type SuperTokensClient struct {
	baseURL    string
	httpClient httpclient.Doer
}

func NewSuperTokensClient() (*SuperTokensClient, error) {
	baseURL := viper.GetString("supertokens.appInfo.apiDomain") + viper.GetString("supertokens.appInfo.apiBasePath")

	err := initSuperTokens()
	if err != nil {
		return nil, fmt.Errorf("failed to init supertokens: %w", err)
	}

	return &SuperTokensClient{
		baseURL:    baseURL,
		httpClient: httpclient.New(),
	}, nil
}

func (c *SuperTokensClient) CreateUser(ctx context.Context, email, password string) (string, error) {
	payload := fmt.Sprintf(`{"formFields": [{"id": "email","value": "%v"},{"id": "password","value": "%v"}]}`, email, password)
	url := fmt.Sprintf("%v/signup", c.baseURL)

	body, err := c.httpClient.Do(ctx, http.MethodPost, url, payload)
	if err != nil {
		return "", err
	}

	type createUserResponse struct {
		Status string `json:"status"`
		User   struct {
			ID         string `json:"id"`
			Email      string `json:"email"`
			TimeJoined int    `json:"timeJoined"`
		} `json:"user"`
	}

	var resp createUserResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", fmt.Errorf("can't unmarshal [%v] to createUserResponse: %w", string(body), err)
	}

	if resp.Status != "OK" {
		return "", fmt.Errorf("error creating user in supertoken: status [%v] not equal 'OK'", resp.Status)
	}

	return resp.User.ID, nil
}

func (c *SuperTokensClient) AddRoleToUser(userID, newRole string) error {
	_, err := userroles.AddRoleToUser(supertokens.DefaultTenantId, userID, newRole)
	if err != nil {
		return fmt.Errorf("error adding role to supertoken user: %w", err)
	}

	return nil
}

func (c *SuperTokensClient) UpdateUserRole(userID, newRole, previousRole string) error {
	_, err := userroles.AddRoleToUser(supertokens.DefaultTenantId, userID, newRole)
	if err != nil {
		return fmt.Errorf("error adding new role to supertoken user: %w", err)
	}

	_, err = userroles.RemoveUserRole(supertokens.DefaultTenantId, userID, previousRole)
	if err != nil {
		return fmt.Errorf("error removing previous role from supertoken user: %w", err)
	}

	return nil
}

func (c *SuperTokensClient) RevokeAllSessionsForUser(userID string) error {
	tenantID := supertokens.DefaultTenantId

	_, err := session.RevokeAllSessionsForUser(userID, &tenantID)
	if err != nil {
		return fmt.Errorf("error revoking all sessions for user: %w", err)
	}

	return nil
}

func (c *SuperTokensClient) UpdateUserMetadata(userID string, metadata map[string]interface{}) error {
	_, err := usermetadata.UpdateUserMetadata(userID, metadata)
	if err != nil {
		return fmt.Errorf("error revoking all sessions for user: %w", err)
	}

	return nil
}

func getUserRole(userID string) (string, error) {
	tenantID := supertokens.DefaultTenantId

	roles, err := userroles.GetRolesForUser(tenantID, userID, nil)
	if err != nil {
		return "", err
	}

	return roles.OK.Roles[0], nil
}
