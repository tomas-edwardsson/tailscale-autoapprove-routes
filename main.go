package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AuthKey    string `envconfig:"TAILSCALE_AUTHKEY" required:"true"`
	Account    string `envconfig:"TAILSCALE_ACCOUNT" required:"true"`
	RouterName string `envconfig:"ROUTER_NAME" required:"true"`
}

type Device struct {
	ID               string   `json:"id"`
	Hostname         string   `json:"hostname"`
	EnabledRoutes    []string `json:"enabledRoutes"`
	AdvertisedRoutes []string `json:"advertisedRoutes"`
}

func missingAFromListB(a, b []string) []string {
	var missing []string
	for _, route := range a {
		found := false
		for _, r := range b {
			if route == r {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, route)
		}
	}
	return missing
}

func main() {

	var config Config
	envconfig.MustProcess("", &config)

	tailAPI, err := NewTailAPI(config.AuthKey, config.Account)
	if err != nil {
		fmt.Println(err)
		return
	}

	devices, err := tailAPI.Devices()
	var router *Device
	for _, device := range devices {
		if device.Hostname == config.RouterName {
			router = &device
		}
	}

	advertisedRoutes := NewSet(router.AdvertisedRoutes...)
	enabledRoutes := NewSet(router.EnabledRoutes...)

	routesToAdd := advertisedRoutes.Difference(enabledRoutes)
	routesToDelete := enabledRoutes.Difference(advertisedRoutes)

	if len(routesToAdd) > 0 || len(routesToDelete) > 0 {
		fmt.Println("Updating Routes")
		if len(routesToAdd) > 0 {
			fmt.Println("  Adding:", routesToAdd)
		}
		if len(routesToDelete) > 0 {
			fmt.Println("  Deleting:", routesToDelete)
		}

		resp, err := tailAPI.SetDeviceRoutes(
			router.ID,
			advertisedRoutes.ToSlice(),
		)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("  Advertised Routes:", resp.AdvertisedRoutes)
		fmt.Println("  Enabled Routes:", resp.EnabledRoutes)
	} else {
		fmt.Println("Routes are in sync")

	}
}
