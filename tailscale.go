package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TailConfig struct {
	AuthKey string
	Account string
}

func NewTailAPI(apiKey, orgID string) (*tailAPI, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	if orgID == "" {
		return nil, fmt.Errorf("Org ID is required")
	}

	return &tailAPI{
		config: TailConfig{
			AuthKey: apiKey,
			Account: orgID,
		},
	}, nil
}

type tailAPI struct {
	config TailConfig
}

type listDeviceResponse struct {
	Devices []Device `json:"devices"`
}

type setDeviceRoutesRequest struct {
	Routes []string `json:"routes"`
}

type setDeviceRoutesResponse struct {
	AdvertisedRoutes []string `json:"advertisedRoutes"`
	EnabledRoutes    []string `json:"enabledRoutes"`
}

type DeviceRoutes struct {
	AdvertisedRoutes []string
	EnabledRoutes    []string
}

func (t *tailAPI) Devices() ([]Device, error) {
	path := fmt.Sprintf(
		"tailnet/%s/devices?fields=all",
		t.config.Account,
	)
	res, err := t.tailAPIGetRequest(path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp listDeviceResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Devices, nil
}

func (t *tailAPI) SetDeviceRoutes(deviceID string, routes []string) (DeviceRoutes, error) {
	reqBody := setDeviceRoutesRequest{
		Routes: routes,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return DeviceRoutes{}, err
	}

	res, err := t.tailAPIPostRequest(
		fmt.Sprintf("device/%s/routes",
			deviceID,
		),
		bytes.NewReader(body),
	)
	if err != nil {
		return DeviceRoutes{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return DeviceRoutes{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var resp setDeviceRoutesResponse
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return DeviceRoutes{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return DeviceRoutes{
		AdvertisedRoutes: resp.AdvertisedRoutes,
		EnabledRoutes:    resp.EnabledRoutes,
	}, nil
}

func (t *tailAPI) tailAPIRequest(path string, method string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf(
		"https://api.tailscale.com/api/v2/%s",
		path,
	)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if method == "POST" {
		req.Header.Add("Content-Type", "application/json")
	}

	req.Header.Add("Authorization", "Bearer "+t.config.AuthKey)

	return http.DefaultClient.Do(req)
}

func (t *tailAPI) tailAPIGetRequest(path string) (*http.Response, error) {
	return t.tailAPIRequest(path, "GET", nil)
}

func (t *tailAPI) tailAPIPostRequest(path string, body io.Reader) (*http.Response, error) {
	return t.tailAPIRequest(path, "POST", body)
}
