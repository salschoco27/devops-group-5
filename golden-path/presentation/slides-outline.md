---
marp: true
theme: gaia
_class: lead
paginate: true
backgroundColor: #1a1a2e
color: #e2e2e2
---

# 🚀 Self-Service Deployment Golden Path
### Implementasi Jalur Emas & Otomatisasi Kubernetes (Minikube) & Helm

**Kelompok 5 — Operasional Pengembang (DevOps) A**
*Chelsea V. H. · Salsabila R. · Farida Q. A. · Nayla R. Z. · Adlya I. A. · Aisyah R. · Nayyara A.*

---

## 🧭 Segmen 1: Latar Belakang & Analisis Gap
### (Durasi: 1 Menit | 1 Slide)

*   **Masalah Operasional Tradisional**:
    *   Pengembang (*developer*) terbebani oleh penulisan manifes Kubernetes (`YAML`) manual.
    *   *Configuration drift* antar lingkungan (`dev` vs `prod`).
    *   Waktu setup lokal (*onboarding*) memakan waktu berjam-jam hingga berhari-hari.
*   **Solusi Jalur Emas (*Golden Path*)**:
    *   Menyediakan template standar yang aman dan patuh aturan (*self-service*).
    *   Mengotomatiskan pipeline deployment tanpa mengharuskan developer menjadi pakar infra.

---

## 📚 Segmen 2: Sorotan Paper 1 (van de Kamp et al., 2024)
### (Durasi: 1.5 Menit | Slide 1 dari 2)

**Paving the Path Towards Platform Engineering Using a Comprehensive Reference Model**
*   **Konsep Platform Engineering Reference Model (PE-RM)**:
    *   Meningkatkan efisiensi organisasi melalui penyediaan portal mandiri (*self-service developer portal*).
    *   **Beban Kognitif (*Cognitive Load*)**: Fokus utama adalah mengurangi kebingungan pengembang dengan menyembunyikan kompleksitas infrastruktur.
    *   **Evolusi Bersama**: Pemeliharaan template *Golden Path* dilakukan secara kolaboratif bersama perwakilan developer (*developer guilds*).

---

## 📚 Segmen 2: Penerapan Paper 1 Pada Proyek Kita
### (Durasi: 1.5 Menit | Slide 2 dari 2)

*   **Implementasi Konsep PE-RM**:
    *   Penyediaan **Helm Chart** terstandar sebagai *Self-Service package*.
    *   Menyembunyikan kerumitan manifestasi Pod, Service, dan Namespace di balik berkas konfigurasi sederhana `values.yaml`.
*   **Hasil Organisasional**:
    *   Memotong birokrasi request manual ke tim operasional.
    *   Pengembang dapat memicu deployment instan melalui Git-flow.

---

## 📖 Segmen 3: Sorotan Paper 2 (Ghanbari et al., 2026)
### (Durasi: 1.5 Menit | Slide 1 dari 2)

**Using development environment as code for enhancing developer experience**
*   **Development Environment as Code (DEaC)**:
    *   Penggunaan kode deklaratif untuk mengotomatiskan setup perkakas mesin lokal pengembang.
    *   Stres pengembang sering dipicu oleh "pengetahuan tersembunyi" (*hidden knowledge*) dan instruksi Slack yang usang.
*   **Siklus ADR (Action Design Research)**:
    *   Membuktikan bahwa otomatisasi lingkungan kerja lokal meningkatkan kepuasan pengembang (*developer experience* - DX) secara signifikan.

---

## 📖 Segmen 3: Penerapan Paper 2 Pada Proyek Kita
### (Durasi: 1.5 Menit | Slide 2 dari 2)

*   **Penerapan 4 Design Principles (DP1 - DP4)**:
    *   *DP1 (Automated Setup)*: Satu perintah Helm/Docker menyingkirkan jam instalasi manual.
    *   *DP2 (Consistency)*: Penyelarasan environment lokal Minikube dengan replika cloud staging/prod.
    *   *DP3 (Modularity)*: Pengembang bebas menyesuaikan alokasi resource tanpa bentrok dengan OS host.
    *   *DP4 (Embedded Knowledge)*: Mengintegrasikan tabel troubleshooting langsung ke repositori kode.

---

## 🏗️ Segmen 4: Arsitektur & Alur Kerja Sistem
### (Durasi: 2 Menit | Slide 1 dari 2)

```text
Developer ──► Push (GitLab Repo) ──► GitLab CI/CD Pipeline
                                           │
       ┌───────────────────────────────────┼───────────────────────────────────┐
       ▼ (Stage: Validate)                 ▼ (Stage: Build & Push)             ▼ (Stage: Deploy)
   Helm Lint & Go Vet                  Multi-stage (scratch)             Helm Upgrade / Install
                                       GitLab Container Registry         Minikube Cluster
```
*   **Prinsip Keamanan Non-Root**:
    *   Kontainer berjalan dengan UID `10001` (`runAsNonRoot: true`).
    *   Mencegah eskalasi hak akses kontainer host.

---

## 🏗️ Segmen 4: Detail Implementasi Teknis
### (Durasi: 2 Menit | Slide 2 dari 2)

*   **Pilar Otomatisasi Kita**:
    *   **Helm Chart**: Orkestrasi Namespace, Deployment (dengan CPU/Memory limits), dan NodePort Service.
    *   **GitLab CI Engine**: Pipeline 5-tahap (`validate`, `build`, `push`, `deploy-dev`, `deploy-prod`).
    *   **Mekanisme Keamanan**: Terintegrasi Trivy SCA, Gosec SAST, dan Gitleaks secret scanning.
    *   **Skrip Pengukuran**: `measure.sh` berbasis CLI kubectl wait untuk menghitung latensi deployment secara presisi.

---

## 💻 Segmen 5: Live Demo (8 Menit)
### (Tanpa Slide - Beralih ke Terminal & Browser)

*   *Agenda Demo Terjadwal*:
    *   Eksekusi `helm upgrade --install` secara instan.
    *   Verifikasi status Pod dan Service di namespace `golden-path-dev`.
    *   Akses endpoint HTTP `/health` dan `/api/v1/stats`.
    *   Simulasi *Self-Healing* (menghapus Pod aktif dan melihat auto-recovery).
    *   Eksekusi skrip pengukur performa `measure.sh`.

---

## 📊 Segmen 6: Evaluasi Perbandingan Performa
### (Durasi: 3 Menit | Slide 1 dari 2)

*   **Hasil Pengukuran Empiris**:
    *   **Sebelum (Manual)**: **15.10 detik** (tiga perintah apply beruntun).
    *   **Sesudah (Helm Golden Path)**: **5.89 detik** (satu rilis kompilasi terpadu).
    *   **Efisiensi Waktu**: Peningkatan kecepatan sebesar **61%**.

*   **Keberhasilan Kualitatif**:
    *   Mengurangi beban kognitif pengembang (dari menulis puluhan baris manifes YAML menjadi cukup beberapa parameter di `values.yaml`).

---

## 🔍 Segmen 6: Ancaman Validitas & Hipotesis Alternatif
### (Durasi: 3 Menit | Slide 2 dari 2)

*   **Ancaman Validitas Pengukuran**:
    *   *Internal*: Ketergantungan pada cache lokal Docker (melewati latensi pull network).
    *   *Eksternal*: Pengujian pada klaster single-node Minikube, bukan klaster cloud AWS EKS multi-zona.
*   **Hipotesis Alternatif**:
    *   Jika kecepatan melambat di produksi, hambatan utama diprediksi berada pada *bandwidth* penarikan image (*Registry bottleneck*), bukan pada performa pemrosesan manifes Helm.

---

## 🏁 Segmen 7: Penutup & Pembelajaran (*Lessons Learned*)
### (Durasi: 1 Menit | 1 Slide)

*   **Kesimpulan Proyek**:
    *   Platform Engineering terbukti meningkatkan kebahagiaan pengembang (*developer experience*) dengan menyingkirkan aktivitas operasional berulang.
    *   Keberhasilan Golden Path bergantung pada kolaborasi sosiologis antara tim Platform dan tim Developer.
*   **Rencana Masa Depan (Future Work)**:
    *   Membangun Internal Developer Portal berbasis Spotify Backstage.
    *   Menerapkan *Policy-as-Code* (OPA/Kyverno) untuk validasi kepatuhan infrastruktur cloud.
