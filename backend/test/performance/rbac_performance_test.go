package performance

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

func getPerfAdminToken(t *testing.T, baseURL string, c *http.Client) string {
	t.Helper()
	payload := map[string]string{"username": "admin", "password": "123456"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		t.Skipf("后端不可达，跳过 RBAC 性能测试: %v", err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	var v struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(b, &v)
	if v.Token == "" {
		t.Skipf("登录未拿到 token，跳过 RBAC 性能测试: %s", string(b))
	}
	return v.Token
}

func TestRBACAdminAccessLatency(t *testing.T) {
	baseURL := "http://localhost:8889"
	client := &http.Client{Timeout: 10 * time.Second}
	token := getPerfAdminToken(t, baseURL, client)

	iterations := 30
	start := time.Now()
	for i := 0; i < iterations; i++ {
		req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/admin/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("unexpected status: %d", resp.StatusCode)
		}
	}
	duration := time.Since(start)
	avg := duration / time.Duration(iterations)

	// 给宽松阈值，作为回归报警而非极限压测。
	if avg > 800*time.Millisecond {
		t.Fatalf("rbac admin access avg too slow: %v", avg)
	}
}

func TestRBACUnauthorizedResponse(t *testing.T) {
	baseURL := "http://localhost:8889"
	client := &http.Client{Timeout: 10 * time.Second}

	req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/admin/users", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Skipf("后端不可达，跳过 RBAC 性能测试: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusForbidden {
		t.Fatalf("unexpected status for unauthorized request: %d", resp.StatusCode)
	}
}
