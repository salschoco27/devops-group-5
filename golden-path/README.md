# 🚀 Golden Path Self-Service Deployment Framework

Framework ini mengimplementasikan konsep **Golden Path** (Jalur Emas) dan **Self-Service Deployment** berbasis Kubernetes dan Helm. Melalui framework ini, tim pengembang (*development team*) dapat mendeploy aplikasi secara mandiri tanpa harus terbebani oleh kompleksitas pengelolaan manifes Kubernetes secara manual. Tata kelola, kepatuhan keamanan, dan limitasi resource telah diintegrasikan langsung ke dalam template standar oleh tim platform (*platform team*).

---

## 📋 Prerequisites

Sebelum memulai, pastikan perangkat lokal Anda telah terinstal perkakas berikut:

*   **Docker Desktop** (atau Docker Engine): Digunakan sebagai *container runtime* lokal dan *driver* untuk Minikube. [Unduh Docker](https://www.docker.com/products/docker-desktop/)
*   **Minikube v1.30+**: Klaster Kubernetes lokal satu node untuk simulasi deployment. [Panduan Minikube](https://minikube.sigs.k8s.io/docs/start/)
*   **kubectl v1.27+**: CLI Kubernetes untuk mengontrol klaster. [Panduan kubectl](https://kubernetes.io/docs/tasks/tools/)
*   **Helm v3+**: Package manager Kubernetes untuk pengelolaan aplikasi. [Unduh Helm](https://helm.sh/docs/intro/install/)
*   **Akun GitLab** (atau GitLab Self-Hosted): Untuk menjalankan otomatisasi pipeline CI/CD.

### Verifikasi Instalasi Lokal
Pastikan instalasi berhasil dengan menjalankan perintah berikut di terminal:
```bash
docker --version
minikube version
kubectl version --client
helm version
```

---

## ⚡ Quick Start — Setup dari Nol

Ikuti 5 langkah berikut untuk menjalankan dan mendeploy aplikasi menggunakan Golden Path dari awal:

### 1. Clone Repositori
Clone repositori proyek DevOps Anda dan arahkan terminal ke folder `golden-path`:
```bash
git clone https://github.com/salschoco27/devops-group-5.git
cd devops-group-5/golden-path
```

### 2. Jalankan Minikube Cluster
Mulai klaster Minikube lokal dengan alokasi resource minimum yang disarankan:
```bash
minikube delete   # Hapus klaster lama jika ada konflik state
minikube start --driver=docker --memory=4096 --cpus=2

# Verifikasi status klaster
kubectl get nodes
```
*Pastikan status node menunjukkan `Ready`.*

### 3. Buat Kubernetes Namespaces
Golden Path mengisolasi lingkungan menggunakan Namespace. Terapkan manifes namespace standar untuk memisahkan lingkungan `dev` dan `prod`:
```bash
kubectl apply -f implementation/kubernetes/namespace.yaml

# Verifikasi pembuatan namespace
kubectl get namespaces
```
*Pastikan namespace `golden-path-dev` dan `golden-path-prod` muncul dengan label `managed-by=golden-path`.*

### 4. Gunakan Golden Path Pertamamu (Local Apply)
Developer dapat melakukan deployment secara mandiri ke lingkungan lokal menggunakan template standardisasi:
```bash
# 1. Salin template dev ke file konfigurasi layanan Anda
cp implementation/golden-path-templates/deploy-dev.yaml my-service-deploy.yaml

# 2. Buka file 'my-service-deploy.yaml' dan lakukan penggantian variabel:
#    - Ganti {{ SERVICE_NAME }} dengan nama aplikasi Anda (misal: taskflow-app)
#    - Ganti {{ PORT }} dengan port aplikasi Anda (misal: 8080)
#    - Ganti {{ IMAGE_REPO }}/{{ SERVICE_NAME }}:{{ IMAGE_TAG }} dengan image uji coba 
#      (contoh: nginxinc/nginx-unprivileged:alpine)

# 3. Terapkan manifes ke klaster
kubectl apply -f my-service-deploy.yaml

# 4. Verifikasi status pod dan service Anda
kubectl get pods -n golden-path-dev
kubectl get svc -n golden-path-dev
```

### 5. (Opsional) Jalankan via GitLab CI/CD Pipeline
Untuk menjalankan otomatisasi penuh via GitLab:
1. Pastikan Anda telah membuat variabel lingkungan di GitLab repo Anda (**Settings > CI/CD > Variables**):
   *   `KUBECONFIG_BASE64`: File konfigurasi Kubeconfig Anda yang di-encode ke Base64 (`cat ~/.kube/config | base64`).
2. Dorong (*push*) perubahan kode Anda ke branch `develop` (untuk deploy ke dev) atau ke branch `main` (untuk manual trigger deploy ke prod).
3. Pipeline GitLab CI (`implementation/gitlab-ci/.gitlab-ci.yml`) akan secara otomatis mengeksekusi tahapan `validate` ➡️ `build` ➡️ `push` ➡️ `deploy-dev`.

---

## 🛠️ Panduan Integrasi Developer (Go Application)

Agar aplikasi Go Anda dapat berjalan secara optimal di dalam lingkungan klaster Kubernetes Golden Path, pastikan aplikasi Anda memenuhi standar integrasi berikut:

### 1. Port Binding
Aplikasi Go Anda harus dikonfigurasi untuk mendengarkan (*listen*) lalu lintas HTTP pada port dinamis yang diambil dari environment variable. Golden Path secara default menggunakan port **8080**.
```go
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
log.Printf("Starting server on port %s", port)
log.Fatal(http.ListenAndServe(":"+port, nil))
```

### 2. Health Checks (Wajib)
Kubernetes memantau kesehatan kontainer secara otomatis menggunakan Probes. Aplikasi Anda harus menyediakan endpoint HTTP berikut:
*   **Liveness Probe (`GET /health`)**: Mengembalikan status HTTP 200 OK jika aplikasi masih hidup. Jika gagal, Kubernetes akan merestart kontainer.
*   **Readiness Probe (`GET /ready` atau root `/`)**: Mengembalikan status HTTP 200 OK jika aplikasi telah siap melayani traffic (misal: koneksi database telah berhasil diinisialisasi).

### 3. Graceful Shutdown
Untuk mencegah kegagalan koneksi (*connection reset*) saat melakukan rolling update, aplikasi harus menangani sinyal **SIGTERM** untuk menutup koneksi HTTP yang aktif secara elegan sebelum proses mati:
```go
// Menangkap sinyal OS
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)
<-c

// Menghentikan server dengan timeout grace period 15 detik
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()
server.Shutdown(ctx)
```

---

## 📦 Panduan Helm (Golden Path Self-Service)

Helm mempermudah pengelolaan parameter konfigurasi tanpa perlu mengubah manifes dasar Kubernetes. Seluruh konfigurasi dikelola secara deklaratif melalui `values.yaml`.

Arahkan terminal ke direktori Helm Chart:
```bash
cd implementation/helm/golden-path-chart
```

### 1. Validasi & Linting Chart
Verifikasi sintaks dan struktur Helm chart untuk memastikan tidak ada kesalahan konfigurasi:
```bash
helm lint .
```

### 2. Render Uji Coba (Dry-Run)
Cetak manifes YAML hasil kompilasi Helm secara lokal tanpa menerapkannya ke klaster untuk melakukan inspeksi:
```bash
helm template golden-path-test .
```

### 3. Deploy dan Upgrade Aplikasi
Gunakan perintah tunggal ini untuk mendeploy aplikasi pertama kali atau melakukan pembaruan konfigurasi:
```bash
helm upgrade --install golden-path .
```

### 4. Konfigurasi Swalayan (Self-Service) oleh Developer
Jika developer ingin mengubah alokasi resource limit, jumlah replika, atau tipe service, developer cukup mengedit file `values.yaml`:
```yaml
# Contoh peningkatan replika untuk antisipasi traffic tinggi
replicaCount: 3

# Contoh penyesuaian memory limit
resources:
  limits:
    memory: "256Mi"
```
Setelah disimpan, cukup jalankan kembali `helm upgrade --install golden-path .` untuk menerapkan perubahan.

---

## ❓ Troubleshooting

Berikut adalah tabel panduan solusi cepat terhadap kendala yang sering ditemui selama setup dan eksekusi deployment:

| Error | Kemungkinan Penyebab | Solusi |
| :--- | :--- | :--- |
| **`ImagePullBackOff`** / **`ErrImagePull`** | 1. Docker image tidak ditemukan di registry.<br>2. Kredensial penarikan (*imagePullSecrets*) salah atau tidak ada. | 1. Periksa nama registry dan tag image pada template.<br>2. Ganti image sementara ke `nginxinc/nginx-unprivileged:alpine` untuk memverifikasi deployment lokal.<br>3. Pastikan token akses registry GitLab telah dikonfigurasi dengan benar pada klaster. |
| **`Pending pod`** | Resource (CPU/Memory) pada klaster Minikube tidak mencukupi untuk memenuhi alokasi request. | 1. Hentikan aplikasi lain yang berjalan di Docker Desktop Anda.<br>2. Berikan resource lebih besar pada Minikube: `minikube start --memory=4096 --cpus=2` (atau lebih tinggi).<br>3. Sesuaikan limit pada file `values.yaml` agar lebih rendah. |
| **`CrashLoopBackOff`** | 1. Aplikasi mengalami *panic* saat dijalankan.<br>2. Port binding salah (misal: listen di port 3000 tapi probe mengarah ke 8080). | 1. Periksa log kontainer Anda untuk mencari *stack trace* error: `kubectl logs <pod-name> -n golden-path-dev`.<br>2. Pastikan port pada kontainer sama dengan port liveness/readiness probe. |
| **`Dial tcp: connection refused`** | Sinyal jaringan terputus ke klaster Minikube atau status VM/Docker driver Minikube mati. | 1. Jalankan `minikube status` untuk memeriksa keadaan klaster.<br>2. Jika mati, jalankan `minikube start`.<br>3. Hubungkan kubectl kembali dengan menjalankan `kubectl cluster-info`. |
| **`Service not accessible via NodePort`** | NodePort di Kubernetes dibatasi secara default pada rentang `30000-32767`. Port di luar rentang ini akan ditolak. | 1. Pastikan nilai `nodePort` di dalam `service.yaml` atau `values.yaml` berada pada rentang tersebut (misal: `30080`).<br>2. Gunakan perintah `minikube service <service-name> -n <namespace>` untuk membuka jalur akses secara otomatis di browser lokal. |
| **`Failed to apply namespace (Forbidden)`** | Anda mencoba menerapkan konfigurasi tanpa hak akses administratif (RBAC) pada klaster. | Pastikan Anda terhubung ke klaster lokal Minikube dengan hak akses penuh, bukan ke klaster produksi perusahaan yang memiliki kebijakan keamanan ketat. Periksa konteks aktif: `kubectl config current-context`. |
