package thalassa

import (
	"github.com/thalassa-cloud/client-go/audit"
	"github.com/thalassa-cloud/client-go/containerregistry"
	"github.com/thalassa-cloud/client-go/dbaas"
	"github.com/thalassa-cloud/client-go/dns"
	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/iam"
	"github.com/thalassa-cloud/client-go/kms"
	"github.com/thalassa-cloud/client-go/kubernetes"
	"github.com/thalassa-cloud/client-go/me"
	"github.com/thalassa-cloud/client-go/objectstorage"
	"github.com/thalassa-cloud/client-go/observability/prometheus"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/projects"
	"github.com/thalassa-cloud/client-go/quicklaunch"
	"github.com/thalassa-cloud/client-go/quotas"
	"github.com/thalassa-cloud/client-go/secrets"
	"github.com/thalassa-cloud/client-go/tfs"
)

type Client interface {
	Audit() *audit.Client
	DBaaS() *dbaas.Client
	IaaS() *iaas.Client
	IAM() *iam.Client
	Kubernetes() *kubernetes.Client
	Me() *me.Client
	ObjectStorage() *objectstorage.Client
	Quotas() *quotas.Client
	QuickLaunch() *quicklaunch.Client
	Tfs() *tfs.Client
	ObservabilityPrometheus() *prometheus.Client
	ContainerRegistry() *containerregistry.Client
	// KMS returns a client for the Key Management Service.
	KMS() *kms.Client
	// Secrets returns a client for the Secrets Manager.
	Secrets() *secrets.Client
	// DNS returns a client for DNS zones and records.
	DNS() *dns.Client
	// Projects returns a client for organisation-scoped project management.
	Projects() *projects.Client
	// SetOrganisation sets the organisation for the client
	SetOrganisation(organisation string)
	GetClient() client.Client
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

func (c *thalassaCloudClient) DBaaS() *dbaas.Client {
	dbaasClient, err := dbaas.New(c.client)
	if err != nil {
		panic(err)
	}
	return dbaasClient
}

func (c *thalassaCloudClient) IAM() *iam.Client {
	iamClient, err := iam.New(c.client)
	if err != nil {
		panic(err)
	}
	return iamClient
}

func (c *thalassaCloudClient) ObjectStorage() *objectstorage.Client {
	objectStorageClient, err := objectstorage.New(c.client)
	if err != nil {
		panic(err)
	}
	return objectStorageClient
}

func (c *thalassaCloudClient) Quotas() *quotas.Client {
	quotasClient, err := quotas.New(c.client)
	if err != nil {
		panic(err)
	}
	return quotasClient
}

func (c *thalassaCloudClient) QuickLaunch() *quicklaunch.Client {
	ql, err := quicklaunch.New(c.client)
	if err != nil {
		panic(err)
	}
	return ql
}

func (c *thalassaCloudClient) Audit() *audit.Client {
	auditClient, err := audit.New(c.client)
	if err != nil {
		panic(err)
	}
	return auditClient
}

func (c *thalassaCloudClient) Tfs() *tfs.Client {
	tfsClient, err := tfs.New(c.client)
	if err != nil {
		panic(err)
	}
	return tfsClient
}

func (c *thalassaCloudClient) ObservabilityPrometheus() *prometheus.Client {
	prometheusClient, err := prometheus.New(c.client)
	if err != nil {
		panic(err)
	}
	return prometheusClient
}

// container registry
func (c *thalassaCloudClient) ContainerRegistry() *containerregistry.Client {
	containerRegistryClient, err := containerregistry.New(c.client)
	if err != nil {
		panic(err)
	}
	return containerRegistryClient
}

func (c *thalassaCloudClient) KMS() *kms.Client {
	kmsClient, err := kms.New(c.client)
	if err != nil {
		panic(err)
	}
	return kmsClient
}

func (c *thalassaCloudClient) Secrets() *secrets.Client {
	secretsClient, err := secrets.New(c.client)
	if err != nil {
		panic(err)
	}
	return secretsClient
}

func (c *thalassaCloudClient) DNS() *dns.Client {
	dnsClient, err := dns.New(c.client)
	if err != nil {
		panic(err)
	}
	return dnsClient
}

func (c *thalassaCloudClient) Projects() *projects.Client {
	projectsClient, err := projects.New(c.client)
	if err != nil {
		panic(err)
	}
	return projectsClient
}

func (c *thalassaCloudClient) GetClient() client.Client {
	return c.client
}
