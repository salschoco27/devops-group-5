1. Buka **Docker Desktop** (pastikan menyala).
2. Jalankan command ini berurutan di terminal:

```bash
minikube start --memory=3000 --cpus=2
kubectl apply -f golden-path/implementation/kubernetes/namespace.yaml
```

# Panduan Golden Path Deployment (Go Developer)
Folder ini berisi template manifest Kubernetes standar (Golden Path) yang telah diatur untuk keamanan, stabilitas, dan skalabilitas. Developer hanya perlu fokus pada kode, sementara infrastruktur ditangani oleh templat ini via CI/CD.

## Persyaratan Aplikasi Go
Agar aplikasi dapat berjalan sempurna di lingkungan Kubernetes, pastikan kode Go memenuhi kriteria berikut:

**1. Port Binding**: Aplikasi harus mendengarkan *(listen)* pada port yang didefinisikan dalam variabel `{{ PORT }}` (Default: 8080).

**2. Health Checks (WAJIB)**: Aplikasi wajib mengekspos endpoint HTTP untuk pemeriksaan kesehatan:
- `GET /health` : Digunakan sebagai **Liveness Probe** (Pengecekan apakah container hidup).
- `GET /ready` atau root port: Digunakan sebagai **Readiness Probe** (Pengecekan apakah aplikasi siap menerima traffic).
**3. Graceful Shutdown**: Disarankan aplikasi menangani **SIGTERM** agar Rolling Update berjalan tanpa error *connection reset*.

## Cara Kerja Otomatisasi (CI/CD)

Tidak perlu mengubah file YAML secara manual. Pipeline GitLab CI akan melakukan substitusi otomatis:
- `{{ SERVICE_NAME }}`: Akan otomatis terisi nama repositori.
- `{{ IMAGE_TAG }}`: Akan otomatis menggunakan 8 karakter pertama **Commit SHA** (contoh: `a1b2c3d4`).
- `{{ IMAGE_REPO }}`: Diambil dari registry internal kelompok kita.

## Aturan Tata Kelola (Governance)
* **Metadata**: Jangan menghapus label `managed-by: golden-path`. Ini digunakan oleh tim Platform untuk audit resource klaster.
* **Resources**: Templat **Production** menggunakan *limit* resource yang ketat. Jika aplikasi Anda membutuhkan memori lebih dari *512Mi*, segera hubungi (Anggota 2🐈‍⬛).
* **Deployment**:
  - Push ke branch `develop` → Auto-deploy ke namespace `golden-path-dev`.
  - Push ke branch `main` → Manual trigger (klik 'Play' di GitLab) ke `golden-path-prod`.