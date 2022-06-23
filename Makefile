.PHONY: kube-release kube-remove
kube-release:
	docker build -t app .
	kubectl apply -f scripts/mongo-deployment.yml
	kubectl apply -f deployment.yml

kube-remove:
	kubectl delete -f scripts/mongo-deployment.yml
	kubectl delete -f deployment.yml

tests:
	go test ./...