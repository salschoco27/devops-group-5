# Design Decisions: Golden Path Implementation

## Keputusan 1: Mengapa golden path (bukan free-form deployment)
**Yang dipilih**: Mengadopsi pendekatan *Golden Path* berupa template yang sudah terstandarisasi untuk proses deployment.
**Alternatif yang dipertimbangkan**: *Free-form deployment* di mana setiap developer menulis konfigurasi Kubernetes (`Deployment`, `Service`, dll) dan pipeline CI/CD mereka sendiri dari awal.
**Justifikasi dari paper**:
  van de Kamp et al. (2024) menunjukkan bahwa *golden paths* memberikan "opinionated and supported approaches" yang memandu tim development dalam mengonfigurasi domain aplikasi tanpa harus memiliki pengetahuan mendalam tentang platform dasar (dibahas dalam *Computational Viewpoint*). Selain itu, Ghanbari et al. (2026) menekankan bahwa *Automated setup* (DP1) secara signifikan mengurangi porsi error manusia dan waktu onboarding. Dalam konteks pipeline GitLab CI + Minikube kami, ini berarti developer hanya perlu memikirkan kode mereka, sementara platform mengurus sisanya.
**Trade-off yang diterima**: Mengurangi fleksibilitas developer untuk melakukan kustomisasi infrastruktur yang ekstrem karena mereka dibatasi oleh standar template yang disediakan oleh *platform engineer*.

## Keputusan 2: Mengapa Minikube sebagai cluster target
**Yang dipilih**: Menggunakan Minikube sebagai *cluster* Kubernetes lokal untuk target eksekusi deployment.
**Alternatif yang dipertimbangkan**: Menggunakan *managed Kubernetes service* berbayar di cloud (seperti EKS AWS, GKE Google, atau AKS Azure).
**Justifikasi dari paper**:
  Mengacu pada Ghanbari et al. (2026) tentang prinsip *Consistent environment* (DP2) dan *Development Environment as Code* (DEaC), penggunaan *tool* lokal yang dapat direproduksi (seperti Minikube) sangat relevan untuk mensimulasikan arsitektur produksi di tahap development tanpa biaya cloud.
**Trade-off yang diterima**: Minikube tidak dapat merepresentasikan skenario beban kerja skala besar (*High Availability* atau *multi-node architecture*) karena sumber daya hanya terbatas pada laptop host lokal.

## Keputusan 3: Scope golden path: deployment saja, bukan full IDP
**Yang dipilih**: Membatasi cakupan *Golden Path* hanya pada fase otomatisasi *deployment* ke Kubernetes.
**Alternatif yang dipertimbangkan**: Membangun *Internal Developer Portal* (IDP) secara penuh lengkap dengan antarmuka grafis (GUI) seperti Backstage.
**Justifikasi dari paper**:
  Berdasarkan van de Kamp et al. (2024) pada *Technology Viewpoint*, implementasi IDP utuh membutuhkan investasi yang besar dan waktu yang panjang. Untuk tahap awal proyek ini, memprioritaskan fungsi *Integration & Delivery Plane* (melalui CI/CD) adalah langkah paling pragmatis untuk segera memberikan nilai tambah *self-service* tanpa beban operasional dari perancangan web portal.
**Trade-off yang diterima**: Pengalaman *self-service* dari developer masih berbasis interaksi dengan git repository (mengisi file YAML), bukan pengalaman *click-and-deploy* yang lebih intuitif via portal web.

## Keputusan 4: Mengapa template-based approach (YAML/Helm)
**Yang dipilih**: Menggunakan template YAML/Helm untuk menyembunyikan kompleksitas konfigurasi Kubernetes dari developer.
**Alternatif yang dipertimbangkan**: Mewajibkan developer untuk menulis manifest Kubernetes (*vanilla*) secara utuh untuk setiap layanannya.
**Justifikasi dari paper**:
  Ghanbari et al. (2026) melalui prinsip *Modular environment* (DP3) dan *Embedded knowledge* (DP4) menyarankan penggunaan *infrastructure-as-code* (IaC) yang mengkapsulasi kerumitan. Dalam kerangka van de Kamp et al. (2024) (*Information Viewpoint*), template bertindak sebagai objek informasi standar. Dengan ini, *best practice* keamanan (seperti batasan resource dan non-root user) sudah *embedded* dalam template, dan developer hanya mengubah variabel sederhana seperti nama aplikasi.
**Trade-off yang diterima**: Beban perwatan kini jatuh kepada tim platform (Anggota 2 dan Anggota 5) yang harus memastikan *template-template* ini terus relevan jika ada pembaruan versi API Kubernetes.

## Keputusan 5: Struktur namespace: satu namespace per environment
**Yang dipilih**: Memisahkan beban kerja lingkungan (dev dan prod) menggunakan isolasi secara logis dengan *Namespace* Kubernetes (misal: `golden-path-dev`).
**Alternatif yang dipertimbangkan**: Menyediakan *cluster* Minikube yang benar-benar terpisah untuk dev dan prod secara fisik, atau mencampur semuanya dalam `default` namespace.
**Justifikasi dari paper**:
  Mengacu pada *Engineering Viewpoint* dari van de Kamp et al. (2024) tentang komponen *Resource Plane*, pemisahan *workload* dapat dilakukan melalui pengaturan komputasi (Compute). *Namespace* memberikan metode isolasi yang efisien tanpa memerlukan *overhead* menjalankan dua virtual machine Minikube sekaligus.
**Trade-off yang diterima**: Kegagalan yang bersifat klaster (seperti *node/VM crash* pada Minikube) akan merusak seluruh lingkungan secara bersamaan, karena mereka sebenarnya hanya terpisah secara virtual di dalam satu mesin fisik yang sama.

## Keputusan 6: Mengapa GitLab CI sebagai trigger golden path
**Yang dipilih**: Menjadikan proses *push* kode ke GitLab sebagai pemicu (trigger) eksekusi deployment melalui GitLab CI.
**Alternatif yang dipertimbangkan**: Developer yang telah selesai mengisi template akan mengeksekusi `kubectl apply` secara manual di terminal mereka ke arah klaster.
**Justifikasi dari paper**:
  Ghanbari et al. (2026) sangat menekankan *Automated setup* (DP1). Di samping itu, berdasarkan *Enterprise Viewpoint* (van de Kamp et al., 2024), ada garis pemisah jelas antara peran developer (membuat fitur) dan platform (mengoperasikan sistem). Memanfaatkan sistem eksternal seperti GitLab CI memastikan bahwa setiap rilis didokumentasikan, dijalankan secara standar/terisolasi, dan sesuai prinsip keamanan, tanpa membiarkan developer memiliki akses administrasi langsung ke Kubernetes.
**Trade-off yang diterima**: Waktu tunggu (*feedback loop*) menjadi sedikit lebih lama karena developer harus menanti antrean *job* *runner* GitLab berjalan alih-alih melihat efek langsung dalam satu detik seperti *command* di terminal lokal.
