package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saturn-platform/saturn-cli/internal/api"
	"github.com/saturn-platform/saturn-cli/internal/models"
)

func TestGitHubAppService_List(t *testing.T) {
	org := "my-org"
	apps := []models.GitHubApp{
		{
			ID:           1,
			UUID:         "github-app-uuid-1",
			Name:         "My GitHub App 1",
			Organization: &org,
			APIURL:       "https://api.github.com",
			HTMLURL:      "https://github.com",
			CustomUser:   "git",
			CustomPort:   22,
		},
		{
			ID:      2,
			UUID:    "github-app-uuid-2",
			Name:    "My GitHub App 2",
			APIURL:  "https://api.github.com",
			HTMLURL: "https://github.com",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/github-apps", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(apps)
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-token")
	svc := NewGitHubAppService(client)

	result, err := svc.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "github-app-uuid-1", result[0].UUID)
	assert.Equal(t, "My GitHub App 1", result[0].Name)
	assert.Equal(t, "my-org", *result[0].Organization)
	assert.Equal(t, "github-app-uuid-2", result[1].UUID)
	assert.Equal(t, "My GitHub App 2", result[1].Name)
}

func TestGitHubAppService_List_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"internal server error"}`))
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-token")
	svc := NewGitHubAppService(client)

	result, err := svc.List(context.Background())
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to list GitHub Apps")
}

func TestGitHubAppService_Get(t *testing.T) {
	app := models.GitHubApp{
		ID:      1,
		UUID:    "github-app-uuid-1",
		Name:    "My GitHub App",
		APIURL:  "https://api.github.com",
		HTMLURL: "https://github.com",
		AppID:   12345,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/github-apps/github-app-uuid-1", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(app)
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-token")
	svc := NewGitHubAppService(client)

	result, err := svc.Get(context.Background(), "github-app-uuid-1")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "github-app-uuid-1", result.UUID)
	assert.Equal(t, "My GitHub App", result.Name)
	assert.Equal(t, 12345, result.AppID)
}

func TestGitHubAppService_Get_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"GitHub App not found"}`))
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-token")
	svc := NewGitHubAppService(client)

	result, err := svc.Get(context.Background(), "nonexistent-uuid")
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get GitHub App nonexistent-uuid")
}

func TestGitHubAppService_Create(t *testing.T) {
	org := "my-org"
	isSystemWide := false
	req := &models.GitHubAppCreateRequest{
		Name:           "New GitHub App",
		Organization:   &org,
		APIURL:         "https://api.github.com",
		HTMLURL:        "https://github.com",
		AppID:          99999,
		InstallationID: 11111,
		ClientID:       "client-id-123",
		ClientSecret:   "client-secret-abc",
		PrivateKeyUUID: "private-key-uuid-1",
		IsSystemWide:   &isSystemWide,
	}

	createdApp := models.GitHubApp{
		ID:           5,
		UUID:         "new-github-app-uuid",
		Name:         "New GitHub App",
		Organization: &org,
		APIURL:       "https://api.github.com",
		HTMLURL:      "https://github.com",
		AppID:        99999,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/github-apps", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		var body models.GitHubAppCreateRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "New GitHub App", body.Name)
		assert.Equal(t, 99999, body.AppID)
		assert.Equal(t, "client-id-123", body.ClientID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(createdApp)
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-token")
	svc := NewGitHubAppService(client)

	result, err := svc.Create(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-github-app-uuid", result.UUID)
	assert.Equal(t, "New GitHub App", result.Name)
	assert.Equal(t, 99999, result.AppID)
}

func TestGitHubAppService_Update(t *testing.T) {
	newName := "Updated GitHub App"
	req := &models.GitHubAppUpdateRequest{
		Name: &newName,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/github-apps/github-app-uuid-1", r.URL.Path)
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		var body models.GitHubAppUpdateRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "Updated GitHub App", *body.Name)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "GitHub App updated successfully"})
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-token")
	svc := NewGitHubAppService(client)

	err := svc.Update(context.Background(), "github-app-uuid-1", req)
	require.NoError(t, err)
}

func TestGitHubAppService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/github-apps/github-app-uuid-1", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-token")
	svc := NewGitHubAppService(client)

	err := svc.Delete(context.Background(), "github-app-uuid-1")
	require.NoError(t, err)
}

func TestGitHubAppService_ListRepositories(t *testing.T) {
	repos := []models.GitHubRepository{
		{
			ID:       101,
			Name:     "my-repo",
			FullName: "my-org/my-repo",
			Private:  false,
			HTMLURL:  "https://github.com/my-org/my-repo",
			CloneURL: "https://github.com/my-org/my-repo.git",
		},
		{
			ID:       102,
			Name:     "private-repo",
			FullName: "my-org/private-repo",
			Private:  true,
			HTMLURL:  "https://github.com/my-org/private-repo",
			CloneURL: "https://github.com/my-org/private-repo.git",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/github-apps/github-app-uuid-1/repositories", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"repositories": repos,
		})
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-token")
	svc := NewGitHubAppService(client)

	result, err := svc.ListRepositories(context.Background(), "github-app-uuid-1")
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "my-org/my-repo", result[0].FullName)
	assert.Equal(t, false, result[0].Private)
	assert.Equal(t, "my-org/private-repo", result[1].FullName)
	assert.Equal(t, true, result[1].Private)
}

func TestGitHubAppService_ListBranches(t *testing.T) {
	branches := []models.GitHubBranch{
		{Name: "main", Protected: true},
		{Name: "develop", Protected: false},
		{Name: "feature/my-feature", Protected: false},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/github-apps/github-app-uuid-1/repositories/my-org/my-repo/branches", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"branches": branches,
		})
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-token")
	svc := NewGitHubAppService(client)

	result, err := svc.ListBranches(context.Background(), "github-app-uuid-1", "my-org", "my-repo")
	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "main", result[0].Name)
	assert.Equal(t, true, result[0].Protected)
	assert.Equal(t, "develop", result[1].Name)
	assert.Equal(t, false, result[1].Protected)
	assert.Equal(t, "feature/my-feature", result[2].Name)
}
