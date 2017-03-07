/*
Copyright 2016 The Archon Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package certificate

import (
	"encoding/json"
	"fmt"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/validation"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/render"
	"net"
)

const (
	ResourceStatusKey   = "archon.kubeup.com/status"
	ResourceTypeKey     = "archon.kubeup.com/type"
	ResourceInstanceKey = "archon.kubeup.com/instance"
	ResourceCAKey       = "archon.kubeup.com/ca"
	CSRKey              = "archon.kubeup.com/csr"
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

func NewCertificateControlFromSecret(secret *api.Secret) (*CertificateControl, error) {
	cacert := secret.Data["tls-cert"]
	cakey := secret.Data["tls-key"]
	policy := &config.Signing{
		Default: config.DefaultConfig(),
	}
	parsedCa, err := helpers.ParseCertificatePEM(cacert)
	if err != nil {
		return nil, err
	}
	priv, err := helpers.ParsePrivateKeyPEMWithPassword(cakey, nil)
	if err != nil {
		return nil, err
	}

	ca, err := local.NewSigner(priv, parsedCa, signer.DefaultSigAlgo(priv), policy)
	if err != nil {
		return nil, err
	}
	cc := &CertificateControl{
		signer: ca,
	}
	return cc, nil
}

func isValidHostname(name string) bool {
	if len(validation.NameIsDNSSubdomain(name, false)) == 0 {
		return true
	}
	if ip := net.ParseIP(name); ip != nil {
		return true
	}
	return false
}

func validator(req *csr.CertificateRequest) error {
	for _, host := range req.Hosts {
		if isValidHostname(host) == false {
			return fmt.Errorf("Invalid hostname for csr: %s", host)
		}
	}
	return nil
}

func (cc *CertificateControl) GenerateCertificate(secret *api.Secret, instance *cluster.Instance) error {
	csrTemplate := secret.Annotations[CSRKey]
	if len(csrTemplate) == 0 {
		return fmt.Errorf("No CSR template in secret annotations")
	}

	if instance != nil {
		renderer, err := render.NewInstanceRenderer(instance)
		if err != nil {
			return fmt.Errorf("Failed to initialize renderer: %v", err)
		}

		csrTemplate, err = renderer.Render("csr", csrTemplate)
		if err != nil {
			return fmt.Errorf("Failed to render csr template: %v", err)
		}
	}

	csrReq := csr.New()
	err := json.Unmarshal([]byte(csrTemplate), csrReq)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal csr: %v", err)
	}

	g := &csr.Generator{Validator: validator}
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

func GenerateCA(secret *api.Secret) error {
	csrTemplate := secret.Annotations[CSRKey]
	if len(csrTemplate) == 0 {
		return fmt.Errorf("No CSR template in secret annotations")
	}

	csrReq := csr.New()
	err := json.Unmarshal([]byte(csrTemplate), csrReq)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal csr: %v", err)
	}

	cert, _, key, err := initca.New(csrReq)
	if err != nil {
		return fmt.Errorf("Failed to initialize ca: %v", err)
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	secret.Data["tls-key"] = key
	secret.Data["tls-cert"] = cert

	secret.Annotations["archon.kubeup.com/status"] = "Ready"

	return nil
}
