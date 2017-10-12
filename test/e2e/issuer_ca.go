/*
Copyright 2017 Jetstack Ltd.
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

package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/jetstack-experimental/cert-manager/pkg/apis/certmanager/v1alpha1"
	"github.com/jetstack-experimental/cert-manager/test/e2e/framework"
	"github.com/jetstack-experimental/cert-manager/test/util"
)

var _ = framework.CertManagerDescribe("CA Issuer", func() {
	f := framework.NewDefaultFramework("create-ca-issuer")

	podName := "test-cert-manager"
	issuerName := "test-ca-issuer"
	secretName := "ca-issuer-signing-keypair"

	BeforeEach(func() {
		By("Creating a cert-manager pod")
		pod, err := f.KubeClientSet.CoreV1().Pods(f.Namespace.Name).Create(NewCertManagerControllerPod(podName))
		Expect(err).NotTo(HaveOccurred())
		err = framework.WaitForPodRunningInNamespace(f.KubeClientSet, pod)
		Expect(err).NotTo(HaveOccurred())
		By("Creating a signing keypair fixture")
		_, err = f.KubeClientSet.CoreV1().Secrets(f.Namespace.Name).Create(newSigningKeypairSecret(secretName))
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		By("Deleting the cert-manager pod")
		err := f.KubeClientSet.CoreV1().Pods(f.Namespace.Name).Delete(podName, nil)
		Expect(err).NotTo(HaveOccurred())
		err = f.KubeClientSet.CoreV1().Secrets(f.Namespace.Name).Delete(secretName, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should generate a signing keypair", func() {
		By("Creating an Issuer")
		_, err := f.CertManagerClientSet.CertmanagerV1alpha1().Issuers(f.Namespace.Name).Create(newCertManagerCAIssuer(issuerName, secretName))
		Expect(err).NotTo(HaveOccurred())
		By("Waiting for Issuer to become Ready")
		err = util.WaitForIssuerCondition(f.CertManagerClientSet.CertmanagerV1alpha1().Issuers(f.Namespace.Name),
			issuerName,
			v1alpha1.IssuerCondition{
				Type:   v1alpha1.IssuerConditionReady,
				Status: v1alpha1.ConditionTrue,
			})
		Expect(err).NotTo(HaveOccurred())
	})
})

func newCertManagerCAIssuer(name, secretName string) *v1alpha1.Issuer {
	return &v1alpha1.Issuer{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1alpha1.IssuerSpec{
			IssuerConfig: v1alpha1.IssuerConfig{
				CA: &v1alpha1.CAIssuer{
					SecretRef: v1alpha1.LocalObjectReference{
						Name: secretName,
					},
				},
			},
		},
	}
}

func newSigningKeypairSecret(name string) *apiv1.Secret {
	return &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		StringData: map[string]string{
			apiv1.TLSCertKey: `-----BEGIN CERTIFICATE-----
MIID4DCCAsigAwIBAgIJAJzTROInmDkQMA0GCSqGSIb3DQEBCwUAMFMxCzAJBgNV
BAYTAlVLMQswCQYDVQQIEwJOQTEVMBMGA1UEChMMY2VydC1tYW5hZ2VyMSAwHgYD
VQQDExdjZXJ0LW1hbmFnZXIgdGVzdGluZyBDQTAeFw0xNzA5MTAxODMzNDNaFw0y
NzA5MDgxODMzNDNaMFMxCzAJBgNVBAYTAlVLMQswCQYDVQQIEwJOQTEVMBMGA1UE
ChMMY2VydC1tYW5hZ2VyMSAwHgYDVQQDExdjZXJ0LW1hbmFnZXIgdGVzdGluZyBD
QTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAM+Q2AO4hARav0qwjk7I
4mEh5R201HS8s7HpaLOXBNvvh7qJ9yJz6jLqYg6EvP0K/bK56Cp2oe2igd7GOxpV
3YPOc3CG0CCqHMprEcvxj2xBKX00Rtcn4oVLhDPhAb0BV/R7NFLeWxzh+ggvPI1X
m1qLaWYqYZEJ5bBsYXD3tPdS4GGINRz8Zvih46f0Z2wVkCGoTpsbX8HO74sa2Day
UjzAsWGlO5bZGiMSHjDEnf9yek2TcjEyVoohoOLaQg/ng21T5RWzeZKTl1cznwuG
Vr9tZfHFqxQ5qeaId+1ICtxNvkEjbTnZl6Wy9Cthn0dxwOeS5TqMJ7SFNXy1gp4j
f/MCAwEAAaOBtjCBszAdBgNVHQ4EFgQUBtrjvWfbkLA0iX6sKVRhKUo864kwgYMG
A1UdIwR8MHqAFAba471n25CwNIl+rClUYSlKPOuJoVekVTBTMQswCQYDVQQGEwJV
SzELMAkGA1UECBMCTkExFTATBgNVBAoTDGNlcnQtbWFuYWdlcjEgMB4GA1UEAxMX
Y2VydC1tYW5hZ2VyIHRlc3RpbmcgQ0GCCQCc00TiJ5g5EDAMBgNVHRMEBTADAQH/
MA0GCSqGSIb3DQEBCwUAA4IBAQCR+jXhup5tCKwhAf8xgvp589BczQOjmotuZGEL
Dcint2y263ChEdsoLhyJfvFCAZfTSm+UT95Hl+ZKVuoVEcAS7udaFUFpC/gIYVOi
H4/uvJps4SpVCB7+T/orcTjZ2ewT23mQAQg+B+iwX9VCof+fadkYOg1XD9/eaj6E
9McXID3iuCXg02RmEOwVMrTggHPwHrOGAilSaZc58cJZHmMYlT5rGrJcWS/AyXnH
VOodKC004yjh7w9aSbCCbAL0tDEnhm4Jrb8cxt7pDWbdEVUeuk9LZRQtluYBnmJU
kQ7ALfUfUh/RUpCV4uI6sEI3NDX2YqQbOtsBD/hNaL1F85FA
-----END CERTIFICATE-----`,
			apiv1.TLSPrivateKeyKey: `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAz5DYA7iEBFq/SrCOTsjiYSHlHbTUdLyzselos5cE2++Huon3
InPqMupiDoS8/Qr9srnoKnah7aKB3sY7GlXdg85zcIbQIKocymsRy/GPbEEpfTRG
1yfihUuEM+EBvQFX9Hs0Ut5bHOH6CC88jVebWotpZiphkQnlsGxhcPe091LgYYg1
HPxm+KHjp/RnbBWQIahOmxtfwc7vixrYNrJSPMCxYaU7ltkaIxIeMMSd/3J6TZNy
MTJWiiGg4tpCD+eDbVPlFbN5kpOXVzOfC4ZWv21l8cWrFDmp5oh37UgK3E2+QSNt
OdmXpbL0K2GfR3HA55LlOowntIU1fLWCniN/8wIDAQABAoIBAQCYvGvIKSG0FpbG
vi6pmLbEZO20s1jW4fiUxT2PUWR49sR4pocdahB/EOvA5TowNcNDnftSK+Ox+q/4
HwRkt6R+Fg/qULmcH7F53dnFqeYw8a42/J3YOvg7v7rzdfISg4eWVobFJ+wBz+Nt
3FyBYWLm+MlBLZSH5rGG5em59/zJNHWIhH+oQPfCxAkYEvd8tXOTUzjhqvEfjaJy
FZghnT9xto4MwDdNCPbtzdNjTMhiv0AHkcZGGtRJfkehXX2qhXOQ2UzzO9XrMZnv
5KgYf+bXKJsyS3SPl6TTl7vg2gKBciRvsdFhMy5I5GyIADrEDJnNNmXQRtiaFLfd
k/aqfPT5AoGBAPquMouZUbVS/Qh+qbls7G4zAuznfCiqdctcKmUGPRP4sTTjWdUp
fjI+UTt1e8hncmr4RY7Oa9kUV/kDwzS5spUZZ+u0PczS3XKxOwNOleoH00dfc9vt
cxctHdPdDTndRi8Z4k3m931jIX7jB/Pyx8qeNYB3pj0k3ThktwMbAVLnAoGBANP4
beI5zpbvtAdExJcuxx2mRDGF0lIdKC0bvQaeqM3Lwqnmc0Fz1dbP7KXDa+SdJWPd
res+NHPZoEPeEJuDTSngXOLNECZe4Ja9frn1TeY858vMJBwIkyc8zu+sgXxjQUM+
TWUlTUhtXyybkRnxAEny4OT2TTgmXITJaKOmV1UVAoGAHaXSlo4YitB42rNYUXTf
dZ0U4H30Qj7+1YFeBjq5qI4GL1IgQsS4hyq1osmfTTFm593bJCunt7HfQbU/NhIs
W9P4ZXkYwgvCYxkw+JAnzNkGFO/mHQG1Ve1hFLiVIt3XuiRejoYdiTfbM02YmDKD
jKQvgbUk9SBSBaRrvLNJ8csCgYAYnrZEnGo+ZcEHRxl+ZdSCwRkSl3SCTRiphJtD
9ZGttYj6quWgKJAhzyyxZC1X9FivbMQSmrsE6bYPq+9J4MpJnuGrBh5mFocHeyMI
/lD5+QEDTsay6twMpqdydxrjE7Q01zuuD9MWIn33dGo6FR/vduJgNatqZipA0hPx
ThS+sQKBgQDh0+cVo1mfYiCkp3IQPB8QYiJ/g2/UBk6pH8ZZDZ+A5td6NveiWO1y
wTEUWkX2qyz9SLxWDGOhdKqxNrLCUSYSOV/5/JQEtBm6K50ArFtrY40JP/T/5KvM
tSK2ayFX1wQ3PuEmewAogy/20tWo80cr556AXA62Utl2PzLK30Db8w==
-----END RSA PRIVATE KEY-----`,
		},
	}
}