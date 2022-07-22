.PHONY: kubernetes
kubernetes:
	k3d cluster create terraform-operator --servers-memory 2024Mib --agents-memory 4048Mib --agents 1

.PHONY: monitoring
monitoring:
	helm upgrade -i prometheus-stack prometheus-community/kube-prometheus-stack -f kube-prometheus-stack.yaml

	sh create-svc.sh