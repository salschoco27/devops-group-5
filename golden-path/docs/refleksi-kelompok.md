# 📝 Refleksi Kelompok — Implementasi Golden Path DevOps

Dokumen ini berisi refleksi kritis kelompok 5 mengenai proses riset, perancangan, dan implementasi *Golden Path* dan *Development Environment as Code* (DEaC) berdasarkan pemahaman terhadap literatur ilmiah acuan dan pengalaman praktis selama pengerjaan proyek.

---

### Pertanyaan 1: Hal apa yang paling mengejutkan atau tidak terduga bagi kelompok Anda dari paper yang telah dibaca? (Minimal 250 kata)

Hal yang paling mengejutkan bagi kelompok kami adalah kesimpulan bahwa beban kerja terbesar dan hambatan utama produktivitas dalam DevOps bukanlah masalah keterbatasan teknologi pipeline, melainkan masalah kognitif manusia (*human cognitive load*) dan sosiologi organisasi. 

Sebelum membaca paper van de Kamp et al. (2024) mengenai *Platform Engineering Reference Model* (PE-RM), kami berasumsi bahwa tim *Platform* memiliki wewenang diktatorial mutlak untuk memaksakan standarisasi perkakas kepada pengembang. Namun, kenyataannya kepemilikan dan evolusi *Golden Path* harus dibagi secara demokratis dengan komunitas pengembang melalui *Dev Guilds* atau *Chapters*. Jika tim platform bertindak otoriter tanpa mendengarkan umpan balik, pengembang akan cenderung menolak standarisasi tersebut dan membuat jalur pintas mereka sendiri (*shadow IT*), yang pada akhirnya merusak tata kelola keamanan organisasi.

Selain itu, temuan dari Ghanbari et al. (2026) mengenai dampak *setup* lingkungan lokal terhadap kesejahteraan psikologis pengembang (*Developer Experience* - DX) sangat membuka mata kami. Kami terkejut melihat data empiris yang menunjukkan seberapa besar kontribusi "pengetahuan tersembunyi" (*hidden knowledge*)—seperti instruksi tidak terdokumentasi yang hanya dibagikan secara informal di Slack—terhadap rasa frustrasi, kemarahan, dan stres pengembang. Paper tersebut membuktikan secara ilmiah bahwa troubleshooting manual untuk masalah sepele seperti versi library yang bentrok dapat menurunkan *self-efficacy* (rasa kompetensi diri) pengembang. 

Kami menyadari bahwa menyediakan otomatisasi bukan sekadar membuat pekerjaan menjadi lebih cepat beberapa detik, melainkan secara aktif melindungi kesehatan mental pengembang dan menjaga fokus kognitif mereka agar tetap berada dalam kondisi kerja optimal (*flow state*). Hubungan langsung antara otomatisasi lingkungan dengan kebahagiaan pengembang (*developer happiness*) adalah wawasan berharga yang belum pernah kami sadari sebelumnya.

---

### Pertanyaan 2: Apa perbedaan mendasar antara apa yang diusulkan di paper dengan apa yang kelompok Anda implementasikan? (Minimal 250 kata)

Perbedaan mendasar antara teori yang diusulkan dalam kedua paper dengan implementasi praktis kelompok kami terletak pada tingkat skala arsitektur (*architectural scale*), kompleksitas teknologi, dan fokus target otomatisasi.

Dalam paper van de Kamp et al. (2024), model referensi rekayasa platform mengasumsikan penggunaan portal pengembang internal (*Internal Developer Portal* - IDP) yang canggih dan tersentralisasi, seperti Spotify Backstage atau Otomi. Portal ini berfungsi sebagai antarmuka grafis di mana pengembang dapat mengklik tombol untuk membuat repositori baru, meluncurkan database, dan memantau status mikroservis secara visual. Sementara itu, implementasi kelompok kami menggunakan pendekatan *Git-based self-service* yang jauh lebih ringan. Kami tidak membangun portal web khusus, melainkan mengandalkan kustomisasi parameter deklaratif pada file `values.yaml` di Helm Chart dan otomatisasi pipeline berbasis git-flow pada GitLab CI/CD (`.gitlab-ci.yml`). Pengembang kami berinteraksi langsung dengan file konfigurasi kode, bukan dengan portal visual.

Di sisi lain, paper Ghanbari et al. (2026) berfokus pada konsep *Development Environment as Code* (DEaC) untuk menyelaraskan **mesin lokal pengembang** (mengintegrasikan IntelliJ IDE dengan kontainer Docker lokal untuk kompilasi kode). Sebaliknya, fokus utama proyek *Golden Path* kelompok kami bergeser ke arah **infrastruktur deployment runtime**. Kami mengotomatiskan siklus deployment aplikasi dari repositori Git langsung ke klaster Kubernetes (Minikube lokal) dengan pengawasan ketat terhadap kebijakan tata kelola sumber daya (*Resource Limits*) dan aturan keamanan (*Security Context* non-root). 

Meskipun prinsip dasar kami sama—yaitu menggunakan kode deklaratif untuk standardisasi—solusi Ghanbari berfokus pada kebebasan berkreasi di mesin lokal pengembang, sedangkan solusi kelompok kami berfokus pada kepatuhan keamanan dan pembatasan operasional saat aplikasi berjalan di lingkungan klaster.

---

### Pertanyaan 3: Jika Anda diberikan waktu tambahan 1 bulan dan akses ke cluster Kubernetes nyata di cloud, perbaikan atau pengembangan apa saja yang akan Anda lakukan? (Minimal 250 kata)

Jika kelompok kami diberikan waktu tambahan satu bulan dan akses ke klaster Kubernetes nyata di penyedia cloud publik (seperti AWS EKS atau Google GKE), kami akan memperluas proyek ini ke dalam tiga pilar pengembangan utama:

Pertama, kami akan mengimplementasikan **Internal Developer Portal (IDP) nyata menggunakan Spotify Backstage**. Kami akan mengintegrasikan Backstage dengan template GitLab kami sehingga pengembang dapat melakukan *onboarding* aplikasi baru melalui antarmuka web sekali klik. Di dalam portal ini, kami akan membangun dasbor pemantauan produktivitas pengembang secara otomatis, yang mengukur metrik DORA (*Deployment Frequency*, *Lead Time for Changes*, *Mean Time to Recovery*, dan *Change Failure Rate*) secara langsung dengan menarik data aktivitas API GitLab CI/CD dan Kubernetes.

Kedua, kami akan meningkatkan **Tata Kelola Keamanan (Governance) berbasis Policy-as-Code**. Di klaster nyata, kami akan memasang agen kebijakan seperti OPA (*Open Policy Agent*) atau Kyverno. Agen ini akan secara otomatis menolak (*block*) setiap deployment pengembang yang tidak mematuhi standar *Golden Path* kelompok kami—misalnya jika kontainer mencoba berjalan sebagai user root atau jika tidak mencantumkan batas alokasi memori di `values.yaml`. Kami juga akan mengintegrasikan *ExternalDNS* dan *cert-manager* agar setiap kali pengembang membuat Service baru, sistem otomatis mendaftarkan nama domain resmi dan menerbitkan sertifikat SSL/TLS HTTPS secara instan dari Let's Encrypt.

Ketiga, kami akan melakukan pengujian **Multi-Node Scheduling dan Latensi Jaringan Nyata**. Kami akan menguji bagaimana mekanisme *self-healing* Kubernetes berperilaku di klaster multi-zona saat terjadi kegagalan node fisik, serta mengukur durasi deployment tanpa *cached image* untuk memahami dampak latensi jaringan internet sesungguhnya saat menarik kontainer dari GitLab Registry ke mesin cloud. Ini akan memberikan data evaluasi performa yang jauh lebih valid secara eksternal.
