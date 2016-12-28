package instance

import (
	"encoding/json"
	"fmt"
	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
	"k8s.io/kubernetes/pkg/api"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/render"
)

const (
	CSRKey = "archon.kubeup.com/csr"
)

type CertificateControlInterface interface {
	GenerateCertificate(secret *api.Secret, instance *cluster.Instance) error
}

type CertificateControl struct {
	signer *local.Signer
}

func NewCertificateControl(caCertFile, caKeyFile string) (*CertificateControl, error) {
	policy := &config.Signing{
		Default: config.DefaultConfig(),
	}
	ca, err := local.NewSignerFromFile(caCertFile, caKeyFile, policy)
	if err != nil {
		return nil, err
	}
	cc := &CertificateControl{
		signer: ca,
	}
	return cc, nil
}

func (cc *CertificateControl) GenerateCertificate(secret *api.Secret, instance *cluster.Instance) error {
	csrTemplate := secret.Annotations[CSRKey]
	if len(csrTemplate) == 0 {
		return fmt.Errorf("No CSR template in secret annotations")
	}
	renderer, err := render.NewInstanceRenderer(instance)
	if err != nil {
		return fmt.Errorf("Failed to initialize renderer: %v", err)
	}

	csrString, err := renderer.Render("csr", csrTemplate)
	if err != nil {
		return fmt.Errorf("Failed to render csr template: %v", err)
	}

	csrReq := csr.New()
	err = json.Unmarshal([]byte(csrString), csrReq)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal csr: %v", err)
	}

	g := &csr.Generator{Validator: genkey.Validator}
	csrBytes, key, err := g.ProcessRequest(csrReq)
	if err != nil {
		return fmt.Errorf("Failed to process csr request: %v", err)
	}
	req := signer.SignRequest{Request: string(csrBytes)}
	certBytes, err := cc.signer.Sign(req)
	if err != nil {
		return fmt.Errorf("Failed to sign csr request: %v", err)
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	secret.Data["tls-key"] = key
	secret.Data["tls-cert"] = certBytes

	secret.Annotations["archon.kubeup.com/status"] = "Ready"

	return nil
}
