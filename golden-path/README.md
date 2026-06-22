1. Buka **Docker Desktop** (pastikan menyala).
2. Jalankan command ini berurutan di terminal:

```bash
minikube start --memory=3000 --cpus=2
kubectl apply -f golden-path/implementation/kubernetes/namespace.yaml
```
