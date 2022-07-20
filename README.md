# ConfigMap Cleaner Webhook

A Kubernetes mutating webhook that cleans ConfigMaps by removing private keys. This is a *demonstration* of a mitigation option for [CVE-2022-2403](https://access.redhat.com/security/cve/CVE-2022-2403).

**THIS IS ONLY DEMO CODE, NOT OFFICIALLY SUPPORTED SOFTWARE. USE AT YOUR OWN RISK.**


## Deploy

Log in to OpenShift cluster, then:

```bash
$ make deploy IMG=quay.io/sfowler/configmap-cleaner:v1
```

OR 

```bash
$ oc apply -k ./config/default
```

Afterwards all ConfigMap creation/update requests named `ouath-serving-cert` will have their `ca-bundle.crt` keys checked for private keys, and if found, those keys will be removed.

Pre-existing ConfigMap objects that contain private keys will need to be removed manually.


## Monitor

After deployment, monitor the webhook status with:

```bash
$ oc logs -n configmap-cleaner -f deploy/configmap-cleaner-manager -c manager
```

Succesful cleanings will look like:

```bash
{"level":"info","ts":1658276836.782465,"logger":"configmap-cleaner","msg":"configmap oauth-serving-cert create/update request received"}
{"level":"info","ts":1658276836.7825415,"logger":"configmap-cleaner","msg":"configmap oauth-serving-cert ca-bundle contains private key, cleaning"}
```

Respective configmap updates by the openshift authentication-operator can be seen with:

```bash
$ oc logs -f --tail=20 -n openshift-authentication-operator deploy/authentication-operator | grep -v throttling
I0720 00:23:44.218155       1 event.go:285] Event(v1.ObjectReference{Kind:"Deployment", Namespace:"openshift-authentication-operator", Name:"authentication-operator", UID:"3752d2ea-1c9c-40be-a61f-7e7f47588428", APIVersion:"apps/v1", ResourceVersion:"", FieldPath:""}): type: 'Normal' reason: 'ConfigMapUpdated' Updated ConfigMap/oauth-serving-cert -n openshift-config-managed:
cause by changes in data.ca-bundle.crt
```

These updates attempt to restore the private key to the configmap, but the webhook successfully filters them out.

Finally, check the oauth-serving-cert configmap to ensure that no private key is present:

```bash
$ oc get cm -n openshift-config-managed oauth-serving-cert -o yaml | grep -i private
$
```

## Preview

To preview the deployment, use `--dry-run` to show the objects that will be created.

```bash
$ oc apply -k config/default/ --dry-run=client
namespace/configmap-cleaner configured (dry run)
serviceaccount/configmap-cleaner-manager configured (dry run)
role.rbac.authorization.k8s.io/configmap-cleaner-manager configured (dry run)
clusterrole.rbac.authorization.k8s.io/configmap-cleaner-metrics-reader configured (dry run)
clusterrole.rbac.authorization.k8s.io/configmap-cleaner-proxy-role configured (dry run)
clusterrolebinding.rbac.authorization.k8s.io/configmap-cleaner-manager-rolebinding configured (dry run)
clusterrolebinding.rbac.authorization.k8s.io/configmap-cleaner-proxy-rolebinding configured (dry run)
configmap/configmap-cleaner-manager-config configured (dry run)
service/configmap-cleaner-manager-metrics-service configured (dry run)
service/configmap-cleaner-webhook-service configured (dry run)
deployment.apps/configmap-cleaner-manager configured (dry run)
mutatingwebhookconfiguration.admissionregistration.k8s.io/configmap-cleaner-webhook configured (dry run)
```

## Build From Source

```bash
$ make container-build IMG=example/image/name
$ make container-push IMG=example/image/name
$ make deploy IMG=example/image/name
```
