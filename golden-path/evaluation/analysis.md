# Analisis Performa dan Evaluasi Implementasi Golden Path

Evaluasi ini menyajikan analisis komparatif mendalam mengenai efisiensi waktu deployment dan beban kerja operasional sebelum dan sesudah penerapan *Self-Service Golden Path*. Evaluasi ini juga menghubungkan hasil empiris proyek dengan temuan literatur ilmiah dari paper van de Kamp et al. (2024) dan Ghanbari et al. (2026), serta mendiskusikan ancaman terhadap validitas pengukuran.

---

## 1. Analisis Komparatif Performa Sebelum & Sesudah Golden Path

Pengukuran durasi deployment dilakukan menggunakan skrip otomatis `measure.sh` yang menghitung selisih waktu (dalam detik) sejak perintah deployment dieksekusi hingga seluruh Pod aplikasi berada dalam status `Ready` dan lulus uji *health check*.

| Parameter Pengukuran | Sebelum (*Manual Kubernetes Manifest*) | Sesudah (*Helm Golden Path*) | Selisih / Peningkatan Efisiensi |
| :--- | :---: | :---: | :---: |
| **Durasi Deployment** | **15,1 detik** | **5,89 detik** | **⚡ 9,21 detik (Efisiensi 61%)** |
| **Jumlah Perintah Manual** | 3 Perintah (`kubectl apply` beruntun) | 1 Perintah (`helm upgrade --install`) | Mengurangi langkah operasional sebesar 66% |
| **Potensi *Human Error*** | Tinggi (Kesalahan penulisan namespace/port) | Sangat Rendah (Konfigurasi terpusat di `values.yaml`) | Tereliminasi lewat validasi skema |
| **Isolasi Lingkungan** | Manual (Rentan tercampur) | Otomatis (Terisolasi per namespace) | Terjamin oleh template standar |

### Analisis Hasil Durasi:
Dari data di atas, terlihat adanya penurunan durasi deployment yang signifikan dari **15,10 detik menjadi 5,89 detik** (peningkatan kecepatan sekitar **61%**). 

Pada metode sebelum (*manual*), pengembang harus menerapkan tiga manifes YAML secara manual (`namespace.yaml`, `deployment.yaml`, dan `service.yaml`) secara berurutan. Setiap perintah memerlukan waktu eksekusi CLI dan pemrosesan API server Kubernetes secara terpisah. 

Setelah menggunakan *Helm Golden Path*, seluruh manifest tersebut dikompilasi menjadi satu kesatuan (*single release unit*). Helm mengelola dependensi antar-resource secara internal dan mengirimkannya ke API server Kubernetes dalam satu payload terpadu, mengurangi overhead komunikasi jaringan dan pemrosesan transaksi API server.

---

## 2. Asesmen Kejujuran atas Keberhasilan Proyek (*Honest Success Assessment*)

### Apa yang Berhasil Sesuai Harapan:
1.  **Otomatisasi Penuh Sekali Klik**: Pengembang dapat mendeploy seluruh stack aplikasi (Namespace, Deployment dengan *Resource Limits* dan *Security Context* non-root, serta Service) hanya dengan satu perintah `helm upgrade --install` atau melalui pemicuan otomatis di pipeline GitLab CI.
2.  **Standardisasi Tata Kelola (*Governance*)**: Tim platform berhasil mengunci parameter keamanan kritis (seperti pencegahan privilege escalation dan pembatasan alokasi CPU/Memory minimum) di dalam template Helm Chart. Pengembang tidak perlu lagi mempelajari sintaks keamanan YAML Kubernetes yang rumit; mereka cukup menyesuaikan nilai-nilai parameter di `values.yaml` secara mandiri.
3.  **Konsistensi Deployment**: Uji coba berulang membuktikan bahwa struktur deployment di lingkungan `golden-path-dev` dan `golden-path-prod` konsisten secara arsitektur, mengurangi bug yang disebabkan oleh perbedaan konfigurasi antar-lingkungan (*config drift*).

### Keterbatasan dan Kendala yang Ditemui:
1.  **Ketergantungan pada Cache Lokal**: Pengukuran durasi **5,89 detik** sangat dipengaruhi oleh penggunaan kontainer yang sudah tersimpan di *local cache* Docker desktop pada mesin pengujian. Jika dijalankan pada klaster baru tanpa *cached image*, durasi sesungguhnya akan didominasi oleh waktu penarikan image (*image pull latency*) dari registry eksternal.
2.  **Kompleksitas Awal Helm**: Bagi pengembang yang belum akrab dengan pola templating Go di Helm, membaca file di folder `templates/` bisa meningkatkan beban kognitif awal. Standardisasi ini memindahkan beban kerja konfigurasi dari pengembang ke tim pemelihara platform (*platform team*).

---

## 3. Hubungan Mendalam dengan Temuan Ilmiah (*Literature Correlation*)

Hasil pengujian empiris kelompok kami menunjukkan korelasi yang kuat dengan teori-teori peningkatan produktivitas yang diajukan oleh kedua paper acuan:

### A. Hubungan dengan van de Kamp, Bakker, & Zhao (2024)
Paper pertama mengenai *Platform Engineering Reference Model* (PE-RM) menekankan bahwa penyediaan *Golden Path* bertujuan untuk mengurangi beban kognitif (*cognitive load*) pengembang dengan menyediakan portal mandiri (*self-service*). Dalam studi mereka, waktu orientasi (*onboarding time*) anggota baru dipotong dari **7 jam menjadi hanya hitungan menit** melalui otomatisasi template. 

Pada proyek kami, pengembang tidak perlu lagi menulis manifest Kubernetes dari nol atau mencari dokumentasi port mana yang harus dibuka. Keberhasilan menurunkan durasi deployment ke **5,89 detik** memvalidasi argumen van de Kamp bahwa meminimalkan hambatan operasional (*operational friction*) secara langsung memberikan ruang bagi pengembang untuk fokus pada nilai bisnis aplikasi, bukan pada konfigurasi infrastruktur.

### B. Hubungan dengan Ghanbari, Terimaa, & Koskinen (2026)
Paper kedua mengenai *Development Environment as Code* (DEaC) meneliti bagaimana otomatisasi lokal mengurangi stres dan frustrasi pengembang. Dalam penelitian ADR mereka di Finlandia, setup lingkungan lokal dipangkas dari **beberapa hari menjadi kurang dari 15 menit** melalui otomatisasi kontainer (2 perintah Docker). 

Implementasi *Golden Path* kami menerapkan prinsip-prinsip desain DEaC secara penuh (khususnya **DP1: Automated Setup** dan **DP2: Consistency**). Dengan membungkus port binding, health checks, dan isolasi namespace ke dalam kode deklaratif Helm, kami meminimalkan kesalahan manusia. Efisiensi yang kami peroleh (selisih 9,21 detik secara eksekusi) mencerminkan kepuasan pengembang (*developer satisfaction*) yang diidentifikasi oleh Ghanbari et al., di mana proses yang bebas hambatan meningkatkan moral kerja pengembang karena mereka tidak perlu melakukan troubleshooting manual pada konfigurasi jaringan yang rusak.

---

## 4. Ancaman terhadap Validitas Pengukuran (*Threats to Validity*)

Untuk menjaga objektivitas ilmiah, kami mengidentifikasi beberapa faktor yang dapat mengancam validitas hasil pengukuran kami:

### A. Validitas Internal (*Internal Validity*)
*   **Efek Image Caching**: Pengujian dilakukan pada mesin lokal di mana *image* aplikasi Go (`nginx` / `sample-app`) telah di-pull sebelumnya. Durasi deployment sesungguhnya di klaster produksi tanpa cache akan lebih lama karena bergantung pada *bandwidth* jaringan internet untuk menarik image dari GitLab Container Registry.
*   **Variasi Kinerja Host**: Skrip `measure.sh` mengukur durasi berdasarkan resource CPU dan memori laptop host penguji. Performa CPU throttling atau konsumsi memori background dari aplikasi lain dapat membiaskan hasil pengukuran waktu respon API Kubernetes.

### B. Validitas Eksternal (*External Validity*)
*   **Skala Aplikasi Tunggal**: Pengujian kami hanya melibatkan satu mikroservis sederhana dengan satu Pod. Pada lingkungan industri nyata yang memiliki arsitektur *microservices* kompleks (puluhan Pod, Service Mesh, Ingress Controller, PV/PVC, dan database eksternal), durasi kompilasi Helm dan penjadwalan pod (*scheduling overhead*) akan jauh lebih tinggi dan mungkin tidak linier dengan hasil pengujian lokal kami.
*   **Lingkungan Minikube vs Cloud Nyata**: Klaster lokal satu-node Minikube tidak merepresentasikan latensi jaringan multi-zone, otorisasi IAM cloud, atau overhead *cold-start* load balancer yang ada pada penyedia cloud publik seperti AWS EKS atau Google GKE.

---

## 5. Hipotesis Alternatif (*Alternative Hypotheses*)

Jika di masa mendatang tim pengembang tidak merasakan peningkatan kecepatan deployment sebesar 61% atau justru mengalami kelambatan, kami merumuskan hipotesis alternatif sebagai berikut:

*   **Hipotesis 1 (*Registry Bottleneck*)**: Latensi jaringan eksternal dan kecepatan autentikasi ke GitLab Container Registry memiliki dampak yang jauh lebih besar terhadap total durasi deployment dibandingkan efisiensi kompilasi manifes Helm.
*   **Hipotesis 2 (*Cluster Resource Constraint*)**: Ketika klaster Kubernetes mengalami kelebihan beban (*resource pressure*), durasi penjadwalan Pod (*pod scheduling state*) akan tetap lama meskipun file manifest dikirim dalam satu payload Helm, karena Kubernetes harus menunggu pod lain di-evict.
*   **Hipotesis 3 (*Stateful App Limitation*)**: Untuk aplikasi yang bersifat *stateful* (membutuhkan inisialisasi volume/database), Helm upgrade akan mengalami hambatan durasi pada proses *mounting* volume yang tidak dipengaruhi oleh otomatisasi jalur emas (*Golden Path*).
