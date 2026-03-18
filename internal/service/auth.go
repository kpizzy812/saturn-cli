package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/saturn-platform/saturn-cli/internal/models"
)

// AuthService handles CLI device auth flow (unauthenticated endpoints)
type AuthService struct {
	httpClient     *http.Client
	pollHTTPClient *http.Client
	baseURL        string
}

// NewAuthService creates a new auth service
func NewAuthService(baseURL string) *AuthService {
	return &AuthService{
		// Used for initiation requests — may take longer
		httpClient: &http.Client{Timeout: 30 * time.Second},
		// Used for polling — must complete well within the poll interval
		pollHTTPClient: &http.Client{Timeout: 8 * time.Second},
		baseURL:        baseURL,
	}
}

// InitDeviceAuth initiates a device authorization session
func (a *AuthService) InitDeviceAuth(ctx context.Context) (*models.DeviceAuthResponse, error) {
	endpoint := a.baseURL + "/api/v1/cli/auth/init"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader([]byte("{}")))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", a.baseURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
	}

	var result models.DeviceAuthResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// CheckAuthStatus checks the current status of the device auth session
func (a *AuthService) CheckAuthStatus(ctx context.Context, secret string) (*models.DeviceAuthStatus, error) {
	endpoint := fmt.Sprintf("%s/api/v1/cli/auth/check?secret=%s", a.baseURL, url.QueryEscape(secret))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := a.pollHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("session not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
	}

	var result models.DeviceAuthStatus
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// PollForToken polls the auth check endpoint until approved, denied, or timeout
func (a *AuthService) PollForToken(ctx context.Context, secret string, interval, timeout time.Duration) (*models.DeviceAuthStatus, error) {
	deadline := time.After(timeout)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-deadline:
			return &models.DeviceAuthStatus{Status: "expired"}, nil
		case <-ticker.C:
			status, err := a.CheckAuthStatus(ctx, secret)
			if err != nil {
				// Transient errors — keep polling
				continue
			}

			switch status.Status {
			case "pending":
				continue
			case "approved", "denied", "expired":
				return status, nil
			default:
				return status, nil
			}
		}
	}
}
