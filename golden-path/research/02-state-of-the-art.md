# State of the Art: Platform Engineering dan Golden Path-Based Self-Service Deployment

Platform Engineering merupakan pendekatan modern dalam pengembangan perangkat lunak yang berfokus pada penyediaan platform internal yang memungkinkan developer melakukan deployment dan operasional aplikasi secara mandiri tanpa harus memahami kompleksitas infrastruktur yang mendasarinya. Van de Kamp et al. (2024) memperkenalkan **Platform Engineering Reference Model (PE-RM)** yang menjelaskan platform engineering melalui lima viewpoint utama, yaitu Enterprise, Information, Computational, Engineering, dan Technology Viewpoint.

### Enterprise Viewpoint
Enterprise Viewpoint berfokus pada tujuan organisasi, stakeholder, serta pembagian tanggung jawab antar tim. Dalam konteks proyek ini, terdapat pemisahan tanggung jawab antara platform team dan development team. Platform team menyediakan Golden Path berupa template deployment yang telah memenuhi standar organisasi, sedangkan development team cukup menggunakan template tersebut untuk melakukan deployment aplikasi. Pendekatan ini mengurangi kompleksitas operasional sekaligus meningkatkan konsistensi deployment antar aplikasi.

### Information Viewpoint
Information Viewpoint mendefinisikan artefak dan informasi yang digunakan dalam platform. Pada implementasi proyek ini, Golden Path direpresentasikan sebagai kumpulan artefak deployment yang terdiri dari Helm Chart, values.yaml, deployment template, service template, dan namespace configuration. Artefak tersebut menjadi sumber kebenaran (*source of truth*) yang digunakan developer saat melakukan deployment aplikasi ke Kubernetes.

### Computational Viewpoint
Computational Viewpoint menjelaskan layanan dan operasi yang disediakan platform kepada pengguna. Pada implementasi ini, layanan utama yang diberikan adalah **self-service deployment** menggunakan Helm. Developer cukup mengubah parameter konfigurasi pada `values.yaml`, seperti jumlah replika, versi image, atau resource limits tanpa perlu mengubah manifest Kubernetes secara langsung. Proses deployment kemudian dilakukan menggunakan perintah Helm yang telah distandarisasi.

### Engineering Viewpoint
Engineering Viewpoint menjelaskan bagaimana komponen platform diorganisasikan dan berinteraksi. Pada implementasi ini terdapat beberapa komponen utama yaitu Kubernetes Cluster (Minikube), Namespace Golden Path, Helm Chart, Deployment Template, dan Service Template. Namespace digunakan untuk mengisolasi lingkungan deployment, sedangkan Helm berperan sebagai abstraction layer yang menyederhanakan interaksi developer dengan Kubernetes.

### Technology Viewpoint
Technology Viewpoint menjelaskan teknologi konkret yang digunakan untuk membangun platform. Dalam proyek ini digunakan Kubernetes sebagai container orchestration platform, Helm sebagai package manager Kubernetes, GitLab CI/CD sebagai automation pipeline, serta Nginx container sebagai aplikasi contoh. Kombinasi teknologi tersebut membentuk fondasi implementasi Golden Path yang sederhana namun tetap merepresentasikan konsep platform engineering yang dijelaskan dalam PE-RM.

---

## Development Environment as Code (Ghanbari et al., 2026)

Ghanbari et al. (2026) memperkenalkan konsep **Development Environment as Code (DEaC)** sebagai evolusi dari Infrastructure as Code (IaC). Jika IaC berfokus pada otomatisasi penyediaan infrastruktur, maka DEaC memperluas konsep tersebut dengan mengotomatisasi konfigurasi lingkungan pengembangan secara menyeluruh.

Dalam DEaC, lingkungan pengembangan diperlakukan sebagai artefak yang dapat didefinisikan, disimpan, dibagikan, dan direproduksi menggunakan kode. Pendekatan ini bertujuan mengurangi konfigurasi manual, meningkatkan konsistensi lingkungan kerja, serta mempercepat proses onboarding developer.

Konsep DEaC sangat relevan dengan implementasi proyek ini karena Helm Chart yang dibangun berfungsi sebagai representasi deployment environment dalam bentuk kode. Semua konfigurasi deployment seperti namespace, resource allocation, replica count, service exposure, dan container image didefinisikan secara deklaratif di dalam Helm Chart dan values.yaml.

Dengan pendekatan ini, deployment dapat direproduksi secara konsisten pada berbagai lingkungan tanpa memerlukan konfigurasi manual yang berulang. Selain itu, perubahan konfigurasi dapat dikelola menggunakan version control sehingga memudahkan audit dan kolaborasi antar anggota tim.

Implementasi ini menunjukkan bagaimana prinsip DEaC dapat diterapkan dalam konteks Kubernetes melalui penggunaan Helm sebagai mekanisme konfigurasi dan deployment yang terstandarisasi.

---

## Tools dan Pendekatan yang Digunakan di Industri

### Backstage

Backstage merupakan Internal Developer Platform (IDP) open-source yang dikembangkan oleh Spotify. Backstage menyediakan portal terpusat bagi developer untuk mengakses dokumentasi, template aplikasi, service catalog, serta workflow deployment. Dalam konteks platform engineering, Backstage sering digunakan sebagai antarmuka utama yang mengimplementasikan konsep Golden Path.

### ArgoCD

ArgoCD merupakan platform GitOps yang melakukan sinkronisasi otomatis antara konfigurasi yang tersimpan di Git dan kondisi aktual cluster Kubernetes. ArgoCD banyak digunakan untuk mengimplementasikan deployment deklaratif dan continuous delivery dalam skala besar.

### Helm

Helm merupakan package manager untuk Kubernetes yang memungkinkan manifest Kubernetes dikemas menjadi template yang dapat digunakan ulang. Helm menyediakan parameterisasi melalui values.yaml sehingga deployment dapat disesuaikan tanpa mengubah manifest utama.

Dalam proyek ini, Helm dipilih karena memiliki tingkat kompleksitas yang lebih rendah dibandingkan Backstage maupun ArgoCD serta lebih realistis untuk diimplementasikan dalam ruang lingkup tugas perkuliahan. Helm juga secara langsung mendukung prinsip Golden Path melalui penyediaan template deployment yang terstandarisasi.

| Tool                | Fungsi Utama               | Tingkat Kompleksitas |
| ------------------- | -------------------------- | -------------------- |
| Backstage           | Internal Developer Portal  | Tinggi               |
| ArgoCD              | GitOps Deployment          | Menengah-Tinggi      |
| Helm                | Kubernetes Package Manager | Menengah             |
| Implementasi Proyek | Helm-Based Golden Path     | Menengah             |

---

## Posisi Implementasi Kita

Implementasi yang dikembangkan dalam proyek ini berada pada level awal Platform Engineering dengan fokus pada **Golden Path-Based Self-Service Deployment**.

Konsep utama yang diadopsi dari literatur adalah:

1. Penyediaan jalur deployment yang terstandarisasi (*Golden Path*).
2. Penggunaan Infrastructure as Code dan Development Environment as Code.
3. Self-service deployment melalui abstraksi konfigurasi menggunakan Helm.
4. Isolasi environment menggunakan namespace Kubernetes.
5. Otomatisasi deployment yang dapat diintegrasikan dengan pipeline CI/CD.

Developer cukup mengubah parameter pada `values.yaml` untuk melakukan perubahan deployment tanpa harus memahami detail manifest Kubernetes. Pendekatan ini meningkatkan *developer experience* sekaligus mengurangi risiko kesalahan konfigurasi.

Sebagai contoh, perubahan jumlah replika tidak lagi dilakukan dengan mengubah file `deployment.yaml`, tetapi cukup mengubah nilai berikut:

```yaml
replicaCount: 3
```

Kemudian deployment diperbarui menggunakan:

```bash
helm upgrade --install golden-path .
```

Pendekatan ini merupakan bentuk implementasi self-service deployment yang menjadi inti dari konsep Golden Path.

Meskipun demikian, implementasi ini masih memiliki beberapa keterbatasan dibandingkan platform engineering modern yang digunakan di industri. Implementasi ini belum mencakup Internal Developer Portal seperti Backstage, belum menerapkan GitOps menggunakan ArgoCD, belum menyediakan observability terintegrasi, dan belum memiliki mekanisme policy enforcement maupun governance otomatis.

Dengan demikian, implementasi yang dikembangkan dapat diposisikan sebagai fondasi platform engineering yang mengimplementasikan prinsip Golden Path dan self-service deployment menggunakan Kubernetes dan Helm. Pendekatan ini berhasil menunjukkan bagaimana konsep yang dijelaskan dalam Platform Engineering Reference Model (van de Kamp et al., 2024) dan Development Environment as Code (Ghanbari et al., 2026) dapat diterapkan dalam skala proyek yang realistis dan sesuai dengan ruang lingkup pembelajaran DevSecOps.
