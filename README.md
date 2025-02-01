Tailscale Autoapprove Routes
----------------------------

This is a simple script that will auto approve all routes for a specific
machine in a tailscale account.

## Usage

All arguments are passed via environment variables

| Variable                    | Description                                   |
|-----------------------------|-----------------------------------------------|
| TAILSCALE_AUTHKEY           | Your tailscale authkey                        |
| TAILSCALE_ACCOUNT           | Your tailscale account                        |
| ROUTER_NAME                 | The name of the machine to approve routes for |

## Example

```bash
export TAILSCALE_AUTHKEY="tskey-1234567890"
export TAILSCALE_ACCOUNT="example.com"
export ROUTER_NAME="my-router"

go run .
```

## Install

```bash
go mod tidy
go build
sudo mv tailscale-autoapprove-routes /usr/local/bin
```
