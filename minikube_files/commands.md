helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-release bitnami/redis --values minikube_config_files/values-minikube.yml