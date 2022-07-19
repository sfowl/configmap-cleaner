# ConfigMap Cleaner Webhook

A Kubernetes mutating webhook that cleans ConfigMaps by removing private keys. This is a *demonstration* of a mitigation option for [CVE-2022-2403](https://access.redhat.com/security/cve/CVE-2022-2403).

**USE AT YOUR OWN RISK**


## Deploy

Log in to OpenShift cluster, then:

```bash
$ make deploy IMG=quay.io/sfowler/configmap-cleaner
```

OR 

```bash
$ oc apply -k ./config/default
```

Afterwards all ConfigMap creation/update requests named `ouath-serving-cert` will have their `ca-bundle.crt` keys checked for private keys, and if found those keys will be removed.

Pre-existing ConfigMap objects that contain private keys will need to be removed manually.


## Monitor

After deployment, monitor the webhook status with:

```bash
$ oc logs -n configmap-cleaner -f deploy/configmap-cleaner-manager
```

Succesfull cleanings will look like:

```bash
1.65821214371673e+09	INFO	entrypoint	config map create/update request received with private key, cleaning
1.6582121437176423e+09	DEBUG	controller-runtime.webhook.webhooks	wrote response	{"webhook": "/mutate-v1-configmap", "code": 200, "reason": "", "UID": "34ad694b-2424-470b-b252-1b239d168ad4", "allowed": true}
```
