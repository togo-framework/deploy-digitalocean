# deploy-digitalocean — docs

**DigitalOcean deploy.** Provision a Droplet (cloud-init Docker) and run the app image.

## Install

```bash
togo install togo-framework/deploy-digitalocean
```

Registers on the [`deploy`](https://github.com/togo-framework/deploy) base; select it with **deploy.provider in togo.yaml (or DEPLOY_PROVIDER)**, then use **`togo deploy`**.

## Interface

`Deployer` — `Provision`/`Deploy`/`Destroy`/`Status` over a `Spec{App,Dir,BuildCmd,Host,User,Image,Region,Domain}` built from your `togo.yaml`.

## Configuration

| Env var | Description |
|---|---|
| `DIGITALOCEAN_TOKEN` | DigitalOcean API token (required). |

## Usage & notes

Uses the godo SDK to create a Droplet whose cloud-init runs `spec.Image` on :8080→:80. `Destroy` deletes it.

## Example

```bash
togo deploy --provider digitalocean --dry-run   # preview the plan
togo deploy --provider digitalocean
```

## Links

- [godo](https://github.com/digitalocean/godo)
- [Marketplace](https://to-go.dev/marketplace)
- [Source](https://github.com/togo-framework/deploy-digitalocean)
