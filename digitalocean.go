// Package digitalocean is a DigitalOcean deploy driver for togo: provisions a
// Droplet (Docker via cloud-init) running the app image. Select with
// deploy.provider=digitalocean; needs DIGITALOCEAN_TOKEN.
package digitalocean

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/digitalocean/godo"
	"github.com/togo-framework/deploy"
	"github.com/togo-framework/togo"
)

func init() { deploy.RegisterDriver("digitalocean", New) }

func New(_ *togo.Kernel) (deploy.Deployer, error) {
	tok := os.Getenv("DIGITALOCEAN_TOKEN")
	if tok == "" {
		return nil, errors.New("deploy-digitalocean: DIGITALOCEAN_TOKEN not set")
	}
	return &driver{c: godo.NewFromToken(tok)}, nil
}

type driver struct{ c *godo.Client }

func cloudInit(image string) string {
	return fmt.Sprintf("#cloud-config\nruncmd:\n  - curl -fsSL https://get.docker.com | sh\n  - docker run -d --name app --restart always -p 80:8080 %s\n", image)
}

func (d *driver) Provision(ctx context.Context, spec deploy.Spec) (*deploy.Result, error) {
	region := spec.Region
	if region == "" {
		region = "fra1"
	}
	size := "s-1vcpu-1gb"
	if v, ok := spec.Options["size"].(string); ok && v != "" {
		size = v
	}
	dr, _, err := d.c.Droplets.Create(ctx, &godo.DropletCreateRequest{
		Name:     spec.App,
		Region:   region,
		Size:     size,
		Image:    godo.DropletCreateImage{Slug: "docker-20-04"},
		UserData: cloudInit(spec.Image),
	})
	if err != nil {
		return nil, fmt.Errorf("digitalocean provision: %w", err)
	}
	return &deploy.Result{Message: "droplet creating; app boots via cloud-init", Raw: map[string]any{"id": dr.ID}}, nil
}

func (d *driver) byName(ctx context.Context, name string) (*godo.Droplet, error) {
	ds, _, err := d.c.Droplets.List(ctx, &godo.ListOptions{PerPage: 200})
	if err != nil {
		return nil, err
	}
	for i := range ds {
		if ds[i].Name == name {
			return &ds[i], nil
		}
	}
	return nil, nil
}

func (d *driver) Deploy(ctx context.Context, spec deploy.Spec) (*deploy.Result, error) {
	dr, err := d.byName(ctx, spec.App)
	if err != nil {
		return nil, err
	}
	if dr == nil {
		return d.Provision(ctx, spec)
	}
	ip, _ := dr.PublicIPv4()
	return &deploy.Result{URL: "http://" + ip, Message: "droplet up; redeploy the container via CI/SSH", Raw: map[string]any{"id": dr.ID}}, nil
}

func (d *driver) Destroy(ctx context.Context, spec deploy.Spec) error {
	dr, err := d.byName(ctx, spec.App)
	if err != nil || dr == nil {
		return err
	}
	_, err = d.c.Droplets.Delete(ctx, dr.ID)
	return err
}

func (d *driver) Status(ctx context.Context, spec deploy.Spec) (*deploy.Status, error) {
	dr, err := d.byName(ctx, spec.App)
	if err != nil {
		return nil, err
	}
	if dr == nil {
		return &deploy.Status{Healthy: false, Detail: "no droplet"}, nil
	}
	ip, _ := dr.PublicIPv4()
	return &deploy.Status{Healthy: dr.Status == "active", Detail: dr.Status, Raw: map[string]any{"ip": ip}}, nil
}
