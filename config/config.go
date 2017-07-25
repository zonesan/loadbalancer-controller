/*
Copyright 2017 Caicloud authors. All rights reserved.

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

package config

import (
	"strings"

	"k8s.io/client-go/kubernetes"

	"github.com/caicloud/loadbalancer-controller/pkg/toleration"
	"github.com/caicloud/loadbalancer-controller/pkg/tprclient"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	defaultIpvsdrImage       = "cargo.caicloud.io/caicloud/loadbalancer-provider-ipvsdr:v0.1.0"
	defaultHTTPBackendImage  = "cargo.caicloud.io/caicloud/default-http-backend:v0.1.0"
	defaultNginxIngressImage = "cargo.caicloud.io/caicloud/nginx-ingress-controller:0.9.0-beta.10"
)

type additionalTolerations []string

func (a *additionalTolerations) Set(value string) error {
	values := strings.Split(value, ",")
	if len(values) == 0 {
		return nil
	}
	*a = append(*a, values...)
	// add additional keys
	toleration.AddAdditionalTolerationKeys(*a)
	return nil
}

func (a *additionalTolerations) String() string {
	return strings.Join(*a, ",")
}

// Configuration contains the global config of controller
type Configuration struct {
	Client                kubernetes.Interface
	TPRClient             tprclient.Interface
	AdditionalTolerations additionalTolerations
	Proxies               Proxies
	Providers             Providers
}

// Proxies contains all cli flags of proxies
type Proxies struct {
	DefaultHTTPBackend    string
	DefaultSSLCertificate string
	Nginx                 ProxyNginx
}

// ProxyNginx contains all cli flags of nginx proxy
type ProxyNginx struct {
	Image string
}

// Providers contains all cli flags of providers
type Providers struct {
	Ipvsdr ProviderIpvsdr
}

// ProviderIpvsdr contains all cli flags of ipvsdr providers
type ProviderIpvsdr struct {
	Image string
}

// AddFlags add flags to app
func (c *Configuration) AddFlags(app *cli.App) {

	flags := []cli.Flag{
		// other
		cli.GenericFlag{
			Name:   "additional-tolerations",
			Usage:  "A comma separated list of k8s `TolerationKeys`",
			EnvVar: "ADDITIONAL_TOLERATIONS",
			Value:  &c.AdditionalTolerations,
		},
		// proxies
		cli.StringFlag{
			Name:        "default-http-backend",
			Usage:       "Default http backend `Image` for ingress controller",
			EnvVar:      "DEFAULT_HTTP_BACKEND",
			Value:       defaultHTTPBackendImage,
			Destination: &c.Proxies.DefaultHTTPBackend,
		},
		cli.StringFlag{
			Name:        "default-ssl-certificate",
			Usage:       "Name of the secret that contains a SSL `certificate` to be used as default for a HTTPS catch-all server",
			EnvVar:      "DEFAULT_SSL_CERTIFICATE",
			Destination: &c.Proxies.DefaultSSLCertificate,
		},
		// nginx
		cli.StringFlag{
			Name:        "proxy-nginx",
			Usage:       "`Image` of nginx ingress controller",
			EnvVar:      "PROXY_NGINX",
			Value:       defaultNginxIngressImage,
			Destination: &c.Proxies.Nginx.Image,
		},
		// ipvsdr
		cli.StringFlag{
			Name:        "provider-ipvsdr",
			Usage:       "`Image` of ipvsdr provider",
			EnvVar:      "PROVIDER_IPVS_DR",
			Value:       defaultIpvsdrImage,
			Destination: &c.Providers.Ipvsdr.Image,
		},
	}
	app.Flags = append(app.Flags, flags...)
}
