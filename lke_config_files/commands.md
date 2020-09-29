helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-release bitnami/redis --values values-production.yml

helm repo add stable https://kubernetes-charts.storage.googleapis.com/
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install nginx-ingress ingress-nginx/ingress-nginx
