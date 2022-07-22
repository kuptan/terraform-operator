.PHONY: kubernetes
kubernetes:
	k3d cluster create terraform-operator --servers-memory 2024Mib --agents-memory 4048Mib --agents 1

.PHONY: monitoring
monitoring:

	sh create-svc.sh