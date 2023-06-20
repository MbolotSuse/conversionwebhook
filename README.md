## Conversion Webhook

An example conversion webhook which uses/is adapted from several sources, including:
- [An example conversion webhook](https://github.com/kubernetes/kubernetes/tree/v1.25.3/test/images/agnhost/crd-conversion-webhook)
- [Controller-Rutime Sig](https://github.com/kubernetes-sigs/controller-runtime/)
- [Kubebuilder](https://book.kubebuilder.io/multiversion-tutorial/tutorial.html)

## Usage

- Build/Push with `docker build . -t $REPO:$TAG Dockerfile && docker push $REPO:$TAG`.
- Update `manifest.yaml` with your image (replace `$REPO` with the repo and `$TAG` with the tag).
- Generate an ssl cert and key with: `openssl req -x509 -newkey rsa:2048 -keyout service.key -out service.crt -sha256 -nodes -days 365 -subj '/CN=conversion-webhook-svc.default.svc' --addext "subjectAltName = DNS:conversion-webhook-svc.default.svc"`.
- Create a secret for this key with `kubectl create secret tls conversion-webhook-secret --cert service.crt --key service.key`
- Deploy with `kubectl create -f manifest.yaml`.
- Use `kubectl get secret conversion-webhook-secret -o yaml`. Retrieve the `tls.crt` field.
- Update `caBundle: ""` in the `crdExamples/crd_webhook.yaml` to use the value retrieved in the previous step.
- Create the CRDs with `kubectl create -f crdExamples/crd_webhook.yaml`
- Create the CRs with `kubectl create -f crdExamples/new.yaml && kubectl create -f crdExamples/old.yaml`

## Purpose

A selective implementation of a CRD conversion webhook, done to explore how many of the components fit together and how they might be used in a wider setting.
