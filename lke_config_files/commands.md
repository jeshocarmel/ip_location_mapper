helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-release bitnami/redis --values values-production.yml