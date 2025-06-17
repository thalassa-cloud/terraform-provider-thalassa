package thalassa

import (
	"github.com/thalassa-cloud/client-go/dbaas/dbaasalphav1"
	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/kubernetes"
	"github.com/thalassa-cloud/client-go/me"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

type Client interface {
	IaaS() *iaas.Client
	Kubernetes() *kubernetes.Client
	Me() *me.Client
	DbaaSAlphaV1() *dbaasalphav1.Client

	// SetOrganisation sets the organisation for the client
	SetOrganisation(organisation string)
}

type thalassaCloudClient struct {
	client client.Client
}

// Option is a function that modifies the Client.
type Option func(*Client) error

// NewClient applies all options, configures authentication, and returns the client.
func NewClient(opts ...client.Option) (Client, error) {
	c, err := client.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &thalassaCloudClient{
		client: c,
	}, nil
}

func (c *thalassaCloudClient) SetOrganisation(organisation string) {
	c.client.SetOrganisation(organisation)
}

func (c *thalassaCloudClient) IaaS() *iaas.Client {
	iaasClient, err := iaas.New(c.client)
	if err != nil {
		panic(err)
	}
	return iaasClient
}

func (c *thalassaCloudClient) Kubernetes() *kubernetes.Client {
	kubernetesClient, err := kubernetes.New(c.client)
	if err != nil {
		panic(err)
	}
	return kubernetesClient
}

func (c *thalassaCloudClient) Me() *me.Client {
	meClient, err := me.New(c.client)
	if err != nil {
		panic(err)
	}
	return meClient
}

func (c *thalassaCloudClient) DbaaSAlphaV1() *dbaasalphav1.Client {
	dbaasClient, err := dbaasalphav1.New(c.client)
	if err != nil {
		panic(err)
	}
	return dbaasClient
}
