# Reading Notes: Paving the Path Towards Platform Engineering Using a Comprehensive Reference Model

## Identitas Paper
- **Judul**: Paving the Path Towards Platform Engineering Using a Comprehensive Reference Model
- **Penulis**: Ruben van de Kamp, Kees Bakker, dan Zhiming Zhao
- **Tahun**: 2024
- **Venue**: EDOC 2023 Workshops, Lecture Notes in Business Information Processing (LNBIP), volume 498
- **DOI**: [10.1007/978-3-031-54712-6_11](https://doi.org/10.1007/978-3-031-54712-6_11)

## Klaim Utama & Cara Membuktikannya
Paper ini mengklaim bahwa kurangnya pemahaman bersama (*shared understanding*) adalah hambatan utama dalam adopsi *Platform Engineering*. Penulis mengajukan **Platform Engineering Reference Model (PE-RM)** sebagai kerangka kerja struktural untuk menciptakan strategi platform yang disesuaikan. 

Cara membuktikannya:
1. **Wawancara Pakar**: Melibatkan delapan pakar rekayasa perangkat lunak dari berbagai latar belakang.
2. **Studi Kasus**: Diimplementasikan pada Wehkamp, organisasi ritel online besar di Belanda, untuk memvalidasi model dalam konteks organisasi nyata.
3. **Eksperimen Usabilitas**: Melakukan pengujian produktivitas terhadap 10 partisipan (dari mahasiswa hingga pengembang senior) untuk mengukur efektivitas *Internal Developer Platform* (IDP) yang dipandu oleh PE-RM.

## Definisi Golden Path dalam Paper Ini
Penulis mendefinisikan *Golden Path* (Jalur Emas) sebagai sekumpulan praktik terbaik, alat, pustaka, dan pilihan arsitektur yang didukung oleh organisasi dan dipelihara oleh *development guilds*. Inti dari konsep ini adalah mitigasi **beban kognitif (cognitive load)** pengembang. Dengan mengikuti *Golden Path*, pengembang dapat membangun, menguji, dan menyebarkan aplikasi dengan konsistensi, keandalan, dan kecepatan yang lebih tinggi tanpa harus memahami kompleksitas infrastruktur di bawahnya secara mendalam.

## 5 Viewpoints PE-RM
1. **Enterprise Viewpoint**: Fokus pada konteks organisasi, peran (*Platform Team* vs *Development Team*), dan siklus hidup proses. Kami mengadopsi ini dengan memisahkan tanggung jawab pengelolaan templat (saya sebagai Anggota 2) dari implementasi fitur aplikasi Go.
2. **Information Viewpoint**: Mengelola objek informasi seperti metadata *Service Catalog* dan dokumentasi *Golden Path*. Proyek kami mengimplementasikan ini melalui standarisasi label metadata pada manifest YAML.
3. **Computational Viewpoint**: Menjelaskan fungsionalitas komponen platform dan interaksi antar antarmuka. Kami mengadopsinya dengan merancang templat YAML yang mengabstraksi fungsionalitas "Deploy" menjadi parameter sederhana bagi developer.
4. **Engineering Viewpoint**: Menggambarkan distribusi komponen ke dalam berbagai "plane" (misal: *Developer Control Plane*). Implementasi kami menggunakan GitLab CI sebagai *Integration & Delivery Plane* utama.
5. **Technology Viewpoint**: Menentukan standar teknologi konkret (Kubernetes, Docker, Helm). Kelompok kami secara eksplisit memilih stack Go, GitLab CI, dan Kubernetes sebagai wujud dari viewpoint ini.

## Hasil Eksperimen Paper
Eksperimen menunjukkan hasil empiris yang signifikan:
- **Onboarding Time**: Sebelum adanya IDP/Golden Path, partisipan memperkirakan butuh waktu rata-rata **7 jam** untuk melakukan *onboarding* aplikasi baru.
- **Efisiensi**: Setelah menggunakan platform berbasis PE-RM, terjadi penurunan dramatis waktu *onboarding* menjadi hitungan menit (sebagaimana divisualisasikan pada Fig. 6 dalam paper).
- **Kepuasan**: Platform mendapatkan rating kemudahan penggunaan **4 dari 5**, dengan partisipan merasa bahwa **90% pekerjaan** infrastruktur telah dikelola secara otomatis oleh platform.

## Komponen Kubernetes dalam Paper
Paper ini memandang Kubernetes sebagai bagian dari **Resource Plane**. Penulis menekankan bahwa platform rekayasa harus memisahkan fitur untuk "membangun aplikasi" dan "menjalankan aplikasi". Integrasi dilakukan melalui pipeline berbasis Git (*Git-driven*) yang menghubungkan *Developer Portal* (seperti Backstage) dengan klaster Kubernetes. Dibandingkan dengan implementasi kita, paper ini menyarankan penggunaan *Provisioning Service* yang lebih kompleks, sementara kita menyederhanakannya menggunakan mekanisme variabel substitusi pada `kubectl` di dalam GitLab CI.

## Asumsi dan Keterbatasan Paper
Berdasarkan Section 5 (Discussion), penulis mengakui:
1. **Konteks Tunggal**: Temuan didasarkan pada satu studi kasus di organisasi retail, sehingga generalisasi pada tipe organisasi lain perlu divalidasi lebih lanjut.
2. **Lingkungan Cloud-Centric**: Penelitian sangat terfokus pada lingkungan *Cloud Native* dan mungkin tidak selaras dengan infrastruktur *on-premises* tradisional.
3. **Trade-off Biaya**: Penelitian belum mengeksplorasi secara mendalam analisis biaya investasi awal pembangunan platform dibandingkan dengan keuntungan jangka panjang yang diperoleh.

## Satu Hal yang dipertanyakan
Secara kritis, kami mempertanyakan fleksibilitas PE-RM terhadap **"Shadow IT"**. Meskipun *Golden Path* dirancang untuk mempermudah, paper ini kurang mengeksplorasi mekanisme penanganan pengembang "pembangkang" yang ingin menggunakan teknologi di luar standar (di luar Jalur Emas) tanpa merusak integritas keamanan platform. Metodologi PE-RM tampak sangat kaku dalam hal tata kelola organisasi, yang mungkin sulit diterapkan pada startup dengan budaya kebebasan teknis yang ekstrem.

## Implikasi untuk Implementasi Kami
1. **Adopsi Template**: Kita mengadopsi penuh konsep *Golden Path* dengan menyediakan template `deploy-dev.yaml` dan `deploy-prod.yaml` untuk mengurangi *cognitive load* tim.
2. **Abstraksi Viewpoint**: Kita menggunakan *Computational Viewpoint* untuk menyembunyikan detail *resource limits* dan *probes* dari pengembang aplikasi Go kita.
3. **Pembatasan Scope**: Mengingat durasi proyek hanya 1 minggu, kita tidak membangun *Developer Portal* penuh (seperti Backstage yang dimention di paper), melainkan menggunakan fitur *GitLab CI Templates* sebagai pengganti portal swalayan tersebut.