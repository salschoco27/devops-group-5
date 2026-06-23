# Gap Analysis: Manual Deployment vs. Golden Path Self-Service

Dokumen ini menganalisis kesenjangan (*gap*) antara metode deployment aplikasi Kubernetes secara manual (tradisional) dengan pendekatan *Golden Path* yang kami rancang. Analisis ini disusun berdasarkan landasan ilmiah dari penelitian platform engineering terbaru.

## Gap yang Diidentifikasi
Dalam siklus hidup pengembangan perangkat lunak (SDLC), proses rilis dan deployment ke lingkungan Kubernetes sering kali menjadi hambatan utama bagi produktivitas pengembang. Berdasarkan penelitian oleh van de Kamp et al. (2024) mengenai *Platform Engineering Reference Model* (PE-RM) dan Ghanbari et al. (2026) mengenai *Development Environment as Code* (DEaC), serta refleksi atas pengalaman tim kami sendiri, kami mengidentifikasi tiga gap utama dalam deployment tradisional:
1. Tidak adanya standarisasi manifes Kubernetes yang berujung pada kerentanan keamanan dan pemborosan resource klaster.
2. Proses onboarding pengembang baru yang memakan waktu lama akibat informasi konfigurasi yang tersebar (*scattered information*).
3. Ketiadaan mekanisme *self-service* mandiri yang memaksa pengembang bergantung penuh pada tim operasional (Ops) untuk setiap perilisan aplikasi.

---

## Gap 1: Tidak Ada Standarisasi Deployment
van de Kamp et al. (2024) menemukan bahwa **kurangnya pemahaman bersama (*lack of a shared understanding*) dalam mendefinisikan strategi platform, serta ketiadaan batas tanggung jawab yang jelas antara platform team dan development team (seperti dibahas pada Section 3.2), memicu hambatan komunikasi yang secara langsung menurunkan produktivitas tim pengembang dan kualitas perangkat lunak (sebagaimana dialami Wehkamp pada studi kasus Section 4.1)**. 

Dalam konteks tim kami, gap ini muncul sebagai kekacauan tata kelola manifestasi YAML. Tanpa adanya *Golden Path*, setiap pengembang menulis manifes Kubernetes mereka sendiri secara *free-form* dari awal. Hal ini berakibat pada:
- **Ketidakpatuhan Keamanan:** Banyak pengembang yang lupa mengonfigurasi *security context* (seperti `runAsNonRoot: true` atau pembatasan user ID), sehingga kontainer berjalan sebagai root yang rentan terhadap serangan eksploitasi kontainer (*container escape*).
- **Miskonfigurasi Resource:** Tidak adanya batas kuota resource (`resources.limits` dan `resources.requests`) yang seragam membuat beberapa aplikasi memakan seluruh memori klaster (efek *noisy neighbor*), sementara aplikasi lainnya kekurangan daya komputasi.
- **Ketiadaan Probes:** Pengembang sering mengabaikan konfigurasi pemeriksaan kesehatan otomatis (`livenessProbe` dan `readinessProbe`), sehingga Kubernetes tidak dapat mendeteksi pod yang mengalami *crash* atau *hang*, merusak garansi zero-downtime saat melakukan rolling update.

---

## Gap 2: Onboarding Developer Baru Terlalu Lama
Ghanbari et al. (2026) menunjukkan bahwa setup environment manual menyebabkan **beban kerja yang berulang-ulang, memakan waktu lama (bisa memakan waktu hingga satu hari penuh untuk mengatasi berbagai error konfigurasi), serta menurunkan kepuasan kerja pengembang (Developer Experience/DX) akibat informasi konfigurasi yang tersebar di berbagai tempat (*scattered information*) dan pengetahuan yang tersembunyi (*hidden knowledge*) seperti dibahas di Section 4.1.2**.

Relevansinya dalam deployment workflow kami sangat nyata. Sebelum menerapkan *Golden Path*, pengembang baru harus melewati proses belajar yang curam (*steep learning curve*) untuk sekadar dapat mendeploy satu *microservice* sederhana ke Kubernetes. Mereka harus:
- Mencari tahu registry mana yang digunakan untuk menyimpan Docker image kelompok.
- Menanyakan kredensial akses klaster ke administrator.
- Memahami struktur namespace yang diperbolehkan.
- Meniru manifes YAML dari proyek lama yang belum tentu menggunakan konfigurasi terbaru yang aman.

Sebagai perbandingan, dalam eksperimen usabilitas van de Kamp et al. (2024) (Section 4.3), waktu onboarding yang semula diestimasikan rata-rata **7 jam** berhasil ditekan menjadi hitungan menit (**23,2 menit** tanpa tutorial dan **19 menit** dengan tutorial seperti ditunjukkan pada Fig. 6) ketika beralih menggunakan platform berbasis PE-RM. Lebih lanjut, pada evaluasi empiris di Ghanbari et al. (2026), otomatisasi menggunakan pendekatan DEaC terbukti mampu menyederhanakan konfigurasi lingkungan yang semula memakan waktu berhari-hari menjadi cukup dengan menjalankan **dua perintah Docker** atau setara dengan mengkloning repositori dalam waktu **15 menit** saja (Section 4.1.3 & Section 4.2.2).

---

## Gap 3: Tidak Ada Self-Service — Developer Harus Minta Bantuan Ops
Pada sistem tradisional yang tidak menerapkan prinsip platform engineering, pengembang tidak memiliki akses mandiri untuk merilis aplikasi mereka ke lingkungan pengujian (*development*) maupun produksi. Setiap kali ada perubahan kode atau kebutuhan rilis fitur baru, pengembang harus membuat tiket permintaan bantuan (*ticket creation*) atau meminta insinyur operasional (Ops) secara manual untuk:
1. Membuatkan namespace baru di klaster Kubernetes.
2. Melakukan deployment manifes YAML yang baru menggunakan kredensial admin milik tim Ops.
3. Mengonfigurasi DNS, routing, atau LoadBalancer eksternal.

Hal ini menciptakan hambatan komunikasi (*inter-team dependency bottleneck*) yang memperlambat frekuensi rilis aplikasi. Pengembang harus menunggu antrean pengerjaan dari tim Ops, sementara tim Ops sendiri terbebani oleh tugas-tugas administratif berulang yang seharusnya bisa diotomatisasi. Pendekatan *Golden Path* memutus rantai ketergantungan ini dengan menyediakan jalur mandiri (*self-service*) berbasis git-push, di mana pengembang bertindak sebagai pemicu deployment melalui GitLab CI tanpa memiliki akses administratif langsung ke dalam klaster Minikube.

---

## Gap yang TIDAK Diselesaikan Proyek Ini
Sebagai refleksi jujur atas keterbatasan waktu pengerjaan (skala implementasi 1 minggu) dan sumber daya komputasi lokal, proyek kami secara eksplisit membatasi cakupan dan tidak menyelesaikan gap berikut:
1. **Penyediaan Infrastruktur Fisik (IaC):** Proyek kami tidak mengotomatisasi penyediaan klaster Kubernetes di cloud menggunakan Terraform atau sejenisnya. Klaster Minikube harus disiapkan secara manual di komputer masing-masing anggota sebelum deployment berjalan.
2. **Platform Portal Grafis (GUI):** Kami tidak membangun *Internal Developer Portal* (IDP) lengkap dengan antarmuka grafis seperti Backstage (sebagaimana didefinisikan dalam *Technology Viewpoint* van de Kamp et al., 2024). Interaksi developer dengan platform masih bersifat *text-based* melalui file konfigurasi git repository.
3. **Penyelarasan Lingkungan Pengembang Lokal (DEaC):** Proyek ini fokus pada otomatisasi deployment di klaster Kubernetes bersama, bukan otomatisasi penyelarasan lingkungan kerja lokal pengembang (*local IDE/Docker setup*) secara menyeluruh seperti konsep DEaC yang diusulkan oleh Ghanbari et al. (2026).

---

## Kesimpulan: Mengapa Golden Path Adalah Solusi yang Tepat

Menerapkan *Golden Path* adalah langkah yang tepat karena solusi ini secara langsung menurunkan beban kognitif pengembang (*cognitive load reduction*) sekaligus menjaga integritas tata kelola sistem infrastruktur (*infrastructure governance*). 

Dengan menyediakan template siap pakai (`deploy-dev.yaml` dan `deploy-prod.yaml`) serta mesin otomatisasi pipeline (`.gitlab-ci.yml`), pengembang dapat berfokus sepenuhnya pada penulisan kode bisnis aplikasi Go mereka. Semua aspek kompleksitas Kubernetes—mulai dari isolasi namespace, konfigurasi non-root, *health check*, hingga pembatasan kuota resource—telah dikelola di balik layar secara transparan dan aman oleh tim platform. Hasilnya adalah siklus rilis yang lebih cepat, stabilitas sistem yang terjamin, dan pengalaman pengembang (*Developer Experience*) yang jauh lebih menyenangkan.