package token

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"ecs_exporter/config"

	"github.com/sirupsen/logrus"
)

type TokenManager struct {
	cfg         *config.APIConfig
	token       string
	tokenMux    sync.RWMutex
	refreshTick *time.Ticker
	stopCh      chan struct{}
}

func NewTokenManager(cfg *config.APIConfig) *TokenManager {
	return &TokenManager{
		cfg:    cfg,
		stopCh: make(chan struct{}),
	}
}

func (tm *TokenManager) Start() {
	// 初次刷新
	if err := tm.refreshToken(); err != nil {
		logrus.Fatalf("initial token fetch failed: %v", err)
	}

	// 定时刷新
	tm.refreshTick = time.NewTicker(time.Duration(tm.cfg.RefreshTokenHours) * time.Hour)

	go func() {
		for {
			select {
			case <-tm.refreshTick.C:
				if err := tm.refreshToken(); err != nil {
					logrus.Errorf("token refresh failed: %v", err)
				}
			case <-tm.stopCh:
				return
			}
		}
	}()
}

func (tm *TokenManager) Stop() {
	close(tm.stopCh)
	tm.refreshTick.Stop()
}

func (tm *TokenManager) GetToken() string {
	tm.tokenMux.RLock()
	defer tm.tokenMux.RUnlock()
	return tm.token
}

func (tm *TokenManager) refreshToken() error {
	body := map[string]interface{}{
		"auth": map[string]interface{}{
			"identity": map[string]interface{}{
				"methods": []string{"password"},
				"password": map[string]interface{}{
					"user": map[string]interface{}{
						"domain": map[string]string{
							"name": tm.cfg.Domain,
						},
						"name":     tm.cfg.Username,
						"password": tm.cfg.Password,
					},
				},
			},
			"scope": map[string]interface{}{
				"domain": map[string]string{
					"name": tm.cfg.Domain,
				},
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", strings.TrimRight(tm.cfg.IAMEndpoint, "/")+"/v3/auth/tokens", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bs, _ := io.ReadAll(resp.Body)
		return errors.New("token request failed, status: " + resp.Status + ", body: " + string(bs))
	}

	token := resp.Header.Get("X-Subject-Token")
	if token == "" {
		return errors.New("token not found in response header")
	}

	tm.tokenMux.Lock()
	tm.token = token
	tm.tokenMux.Unlock()

	logrus.Infof("token refreshed successfully")
	return nil
}
