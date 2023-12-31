/*
Copyright 2018 The Kubernetes Authors.

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

package main

import (
	"flag"
	"fmt"
	"net/http"

	v1 "github.com/Mbolotsuse/conversionwebhook/api/v1"
	v2 "github.com/Mbolotsuse/conversionwebhook/api/v2"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"
)

var (
	certFile string
	keyFile  string
	port     int
	scheme   = runtime.NewScheme()
)

// CmdCrdConversionWebhook is used by agnhost Cobra.
var CmdCrdConversionWebhook = &cobra.Command{
	Use:   "crd-conversion-webhook",
	Short: "Starts HTTP server on port 443 for testing CustomResourceConversionWebhook",
	Long: `The subcommand tests "CustomResourceConversionWebhook".

After deploying it to Kubernetes cluster, the administrator needs to create a "CustomResourceConversion.Webhook" in Kubernetes cluster to use remote webhook for conversions.

The subcommand starts a HTTP server, listening on port 443, and creating the "/crdconvert" endpoint.`,
	Args: cobra.MaximumNArgs(0),
	Run:  run,
}

// Config contains the server (the webhook) cert and key.
type Config struct {
	CertFile string
	KeyFile  string
}

func main() {
	CmdCrdConversionWebhook.Flags().StringVar(&certFile, "tls-cert-file", "",
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
			"after server cert.")
	CmdCrdConversionWebhook.Flags().StringVar(&keyFile, "tls-private-key-file", "",
		"File containing the default x509 private key matching --tls-cert-file.")
	CmdCrdConversionWebhook.Flags().IntVar(&port, "port", 443,
		"Secure port that the webhook listens on")
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(v2.AddToScheme(scheme))
	CmdCrdConversionWebhook.Execute()
}

func run(cmd *cobra.Command, args []string) {
	config := Config{CertFile: certFile, KeyFile: keyFile}

	http.HandleFunc("/crdconvert", conversion.NewWebhookHandler(scheme).ServeHTTP)
	http.HandleFunc("/readyz", func(w http.ResponseWriter, req *http.Request) { w.Write([]byte("ok")) })
	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klogFlags.Set("v", "6")
	klog.InitFlags(klogFlags)
	clientset := getClient()
	err := runControllers()
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		TLSConfig: configTLS(config, clientset),
	}
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		panic(err)
	}
}
