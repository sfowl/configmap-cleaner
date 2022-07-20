package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1-configmap,mutating=true,failurePolicy=fail,groups="",resources=configmaps,verbs=create;update,versions=v1,name=vconfigmap.kb.io,admissionReviewVersions=v1,sideEffects=none

// configMapMutator mutates configMaps
type configMapMutator struct {
	Client  client.Client
	decoder *admission.Decoder
}

// configMapMutator strips private keys from ca-bundles in configmaps
func (v *configMapMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	configMap := &corev1.ConfigMap{}
	log := log.Log.WithName("configmap-cleaner")

	err := v.decoder.Decode(req, configMap)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	log.Info(fmt.Sprintf("configmap %s create/update request received", configMap.Name))
	if configMap.Name == "oauth-serving-cert" {
		if cert, ok := configMap.Data["ca-bundle.crt"]; ok {
			if strings.Contains(cert, "PRIVATE KEY") {
				log.Info(fmt.Sprintf("configmap %s ca-bundle contains private key, cleaning", configMap.Name))

				// shallow copy
				mutatedConfigMap := configMap
				mutatedConfigMap.Data["ca-bundle.crt"] = cleanPrivateKey(cert)
				marshaledConfigMap, err := json.Marshal(mutatedConfigMap)
				if err != nil {
					return admission.Errored(http.StatusInternalServerError, err)
				}

				return admission.PatchResponseFromRaw(req.Object.Raw, marshaledConfigMap)
			}
		}
	}

	return admission.Allowed("")
}

// configMapMutator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (v *configMapMutator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}

func cleanPrivateKey(certBundle string) string {
	var certsFiltered string
	for block, rest := pem.Decode([]byte(certBundle)); block != nil; block, rest = pem.Decode(rest) {
		_, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			// if not a valid cert (i.e. private key), skip
			continue
		}

		certsFiltered += "\n" + string(pem.EncodeToMemory(block))
	}

	return certsFiltered
}
