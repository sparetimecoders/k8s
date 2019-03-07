package external_dns

/**
external_dns:create() {
  printf "Creating ${BLUE}external-dns${NC}\n"

  local cluster_name=$1
  local domain=$2
  cat <<EOF | ${KUBECTL_CMD} apply -f - &>/dev/null
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: external-dns
  namespace: kube-system
spec:
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: external-dns
    spec:
      containers:
      - name: external-dns
        image: registry.opensource.zalan.do/teapot/external-dns:latest
        args:
        - --source=ingress
        - --domain-filter=$domain
        - --provider=aws
        - --policy=sync
        - --aws-zone-type=public
        - --registry=txt
        - --txt-owner-id=$cluster_name
        - --txt-prefix=ext-dns
EOF

*/
