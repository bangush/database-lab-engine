/*
2020 © Postgres.ai
*/

// Package platform provides the Platform API client.
package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"gitlab.com/postgres-ai/database-lab/v3/pkg/log"
)

const (
	// PersonalType defines personal type of the Platform token.
	PersonalType = "personal"

	// AdminType defines administrative type of the Platform token.
	AdminType = "admin"

	// HashType defines hash type of the Platform token.
	HashType = "hash"
)

// TokenCheckRequest represents a token checking request.
type TokenCheckRequest struct {
	Token string `json:"token"`
}

// TokenCheckResponse represents response of a token checking request.
type TokenCheckResponse struct {
	APIResponse
	OrganizationID uint       `json:"org_id"`
	Personal       bool       `json:"is_personal"`
	TokenType      string     `json:"token_type"`
	ValidUntil     *time.Time `json:"valid_until"`
}

// CheckPlatformToken makes an HTTP request to check the Platform Access Token.
func (p *Client) CheckPlatformToken(ctx context.Context, request TokenCheckRequest) (TokenCheckResponse, error) {
	respData := TokenCheckResponse{}

	if err := p.doPost(ctx, "/rpc/dblab_token_check", request, &respData); err != nil {
		return respData, errors.Wrap(err, "failed to post request")
	}

	if respData.Code != "" || respData.Message != "" {
		log.Dbg(fmt.Sprintf("Unsuccessful response given. Request: %v", request))

		return respData, errors.Errorf("error: %v", respData)
	}

	return respData, nil
}
