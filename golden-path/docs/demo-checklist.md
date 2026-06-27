# 📋 Panduan & Checklist Demo Live (8 Menit)

Dokumen operasional ini disusun oleh Anggota 7 (Documentation & Demo Lead) untuk memastikan pelaksanaan demo live di hadapan dosen penguji berjalan lancar tanpa kendala teknis.

---

## 📅 H-1: Persiapan Teknis & Verifikasi

*   [ ] **Verifikasi Kode & Manifes**: Pastikan tidak ada konflik git pada branch `main`. Lakukan `git pull origin main` untuk sinkronisasi akhir.
*   [ ] **Uji Coba Mandiri (Dry-Run)**: Jalankan seluruh alur demo secara lokal minimal dua kali untuk memastikan tidak ada error yang tidak terduga.
*   [ ] **Cadangan Rekaman Video**: Rekam seluruh alur demo yang berhasil menggunakan perekam layar (durasi maks. 8 menit) dan simpan sebagai file `.mp4` lokal di laptop. Ini akan digunakan sebagai fallback jika klaster lokal mengalami masalah saat demo.
*   [ ] **Pembagian Tugas Presentasi**: Konfirmasikan kesiapan setiap anggota kelompok untuk mempresentasikan bagian masing-masing sesuai pembagian waktu.

---

## ⏰ H-30 Menit Sebelum Zoom

*   [ ] **Bebaskan Resource Laptop**: Tutup semua aplikasi berat yang tidak diperlukan (seperti browser dengan banyak tab, IntelliJ, Slack desktop, dll.) untuk menyisakan RAM minimal 8GB untuk Docker & Minikube.
*   [ ] **Nyalakan Minikube**:
    ```bash
    minikube start --driver=docker --memory=4096 --cpus=2
    ```
*   [ ] **Bersihkan Lingkungan Klaster**: Hapus semua resource sisa uji coba sebelumnya agar demo dimulai dari kondisi bersih:
    ```bash
    kubectl delete ns golden-path-dev golden-path-prod --ignore-not-found
    ```
*   [ ] **Persiapkan Tata Letak Terminal (Split Screen)**:
    *   **Terminal Kiri**: Jalankan perintah monitoring Pod secara real-time:
        ```bash
        kubectl get pods -n golden-path-dev -w
        ```
    *   **Terminal Kanan**: Kosongkan layar (`clear`) untuk mengeksekusi perintah Helm dan kubectl.
*   [ ] **Siapkan Browser**: Buka satu jendela browser bersih dan arahkan tab ke:
    *   Halaman repositori GitLab kelompok.
    *   Tab kosong untuk mengakses endpoint aplikasi nanti.
*   [ ] **Tunjuk Anggota Cadangan (Backup Presenter)**: Tunjuk satu anggota kelompok lain yang juga memiliki Docker/Minikube aktif untuk bersiap melakukan *screen share* jika presenter utama mengalami mati listrik atau putus koneksi internet.

---

## ⏱️ Alur Demo Live (Total Durasi: 8 Menit)

| Waktu | Durasi | Topik Demo | Presenter | Tindakan Teknis / Perintah |
| :---: | :---: | :--- | :---: | :--- |
| **00:00 - 01:00** | 1.0 Menit | Pembukaan & Analisis Gap | Anggota 1 | Menjelaskan latar belakang proyek, gap operasional sebelum adanya Golden Path, dan mengenalkan struktur folder `golden-path/`. |
| **01:00 - 02:30** | 1.5 Menit | Pipeline CI/CD & Keamanan | Anggota 3 | Memperlihatkan file `.gitlab-ci.yml` dan log sukses GitLab CI yang menjalankan security scans (Trivy, Gosec, Gitleaks) dan push image ke registry. |
| **02:30 - 04:30** | 2.0 Menit | Deploy dengan Helm Golden Path | Anggota 2 & 5 | 1. Tunjukkan folder Helm Chart.<br>2. Jalankan perintah di Terminal Kanan:<br>`helm upgrade --install golden-path implementation/helm/golden-path-chart/`<br>3. Tunjukkan di Terminal Kiri bahwa Pod langsung terbuat dan berstatus `Running`.<br>4. Jalankan `minikube service taskflow-api -n golden-path-dev` untuk membuka aplikasi di browser.<br>5. Akses endpoint `/health` untuk membuktikan aplikasi aktif. |
| **04:30 - 06:00** | 1.5 Menit | Simulasi Self-Healing & Isolasi Namespace | Anggota 4 | 1. Hapus satu pod secara paksa:<br>`kubectl delete pod <nama-pod> -n golden-path-dev`<br>2. Tunjukkan ke dosen bahwa Kubernetes otomatis membuat pod baru secara instan (*self-healing*).<br>3. Jelaskan pemisahan namespace dev dan prod. |
| **06:00 - 07:30** | 1.5 Menit | Demo Pengukuran & Evaluasi | Anggota 6 | 1. Jalankan skrip pengukuran durasi:<br>`./evaluation/measure.sh`<br>2. Tunjukkan hasil pengukuran waktu (sekitar 5.89 detik) dan bandingkan dengan baseline manual (15.1 detik) sesuai isi berkas `evaluation/analysis.md`. |
| **07:30 - 08:00** | 0.5 Menit | Penutup & Refleksi | Anggota 7 | Memberikan kesimpulan singkat mengenai kontribusi platform engineering terhadap produktivitas tim pengembang dan menutup presentasi. |

---

## 🚨 Rencana Cadangan (Fallback Procedures)

### Kasus A: Minikube lokal macet atau hang saat demo live.
*   **Penyebab**: Alokasi RAM laptop host terlalu tinggi atau terjadi kebocoran memori pada Docker Desktop.
*   **Langkah Solusi**:
    1.  Jangan mencoba restart Minikube secara langsung (memakan waktu > 2 menit).
    2.  Presenter utama segera menghentikan *screen share* dan meminta maaf secara profesional.
    3.  Putar **Video Cadangan (Dry-run)** yang telah direkam pada H-1 untuk menunjukkan alur kerja yang berhasil. Jelaskan langkah demi langkah melalui video tersebut.

### Kasus B: Koneksi Internet lambat sehingga pipeline GitLab menggantung (*stuck*).
*   **Langkah Solusi**:
    1.  Lewati pemutaran log GitLab CI secara langsung.
    2.  Beralih ke Terminal lokal dan lakukan deployment secara manual menggunakan file manifes lokal:
        ```bash
        kubectl apply -f implementation/kubernetes/deployment.yaml
        ```
    3.  Jelaskan bahwa otomatisasi lokal menggunakan manifes yang sama dengan yang dijalankan oleh pipeline GitLab.

### Kasus C: Port tabrakan (*Port Conflict*) saat menjalankan Service.
*   **Gejala**: Perintah `minikube service` memunculkan error port sudah digunakan oleh aplikasi lain di laptop.
*   **Langkah Solusi**:
    1.  Buka file `values.yaml` pada Helm chart secara cepat.
    2.  Ubah nilai `service.nodePort` dari `30080` menjadi `30081` (atau port kosong lainnya).
    3.  Lakukan upgrade release Helm secara instan:
        ```bash
        helm upgrade --install golden-path implementation/helm/golden-path-chart/
        ```
    4.  Jelaskan kepada dosen bahwa ini membuktikan kelenturan (*modularity*) dari sistem *Self-Service* Golden Path.
