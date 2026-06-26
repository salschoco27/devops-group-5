## Perbandingan Before vs After
| Metode | Deployment Time |
|---------|----------------:|
| Kubernetes Manifest | **15,1 detik** |
| Helm Golden Path | **5,89 detik** |

### Apakah Golden Path Berhasil?
Ya. Berdasarkan hasil pengujian, penggunaan **Helm Golden Path** berhasil mengurangi waktu *deployment* dari **15,1 detik** menjadi **5,89 detik**. Selain lebih cepat, proses *deployment* juga menjadi lebih sederhana karena cukup menggunakan satu perintah `helm install`.

### Hubungan dengan Paper
Hasil ini sejalan dengan penelitian **van de Kamp et al. (2024)** yang menyatakan bahwa *Golden Path* dapat menyederhanakan proses pengembangan dan mengurangi waktu yang dibutuhkan untuk melakukan *deployment*. Pada implementasi proyek ini, peningkatan terlihat dari berkurangnya langkah manual dan waktu *deployment*.

### Ancaman terhadap Validitas Pengukuran
Pengujian hanya dilakukan pada satu lingkungan (*Minikube*) dan satu aplikasi contoh, sehingga hasilnya dapat berbeda jika diterapkan pada lingkungan produksi atau aplikasi dengan ukuran yang lebih besar.

### Hipotesis
Apabila peningkatan performa tidak terlalu signifikan, kemungkinan penyebabnya adalah spesifikasi komputer, proses inisialisasi Kubernetes, atau *image* yang sudah tersimpan (*cached*), sehingga waktu *deployment* menjadi lebih cepat dibanding kondisi sebenarnya.