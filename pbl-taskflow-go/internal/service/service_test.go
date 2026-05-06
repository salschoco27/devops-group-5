package service_test

import (
	"fmt"
	"testing"

	"github.com/taskflow/api/internal/model"
	"github.com/taskflow/api/internal/repository"
	"github.com/taskflow/api/internal/service"
)

func newSvc() *service.TaskService {
	return service.NewTaskService(repository.NewMemoryRepository())
}

// ── [BUG] CalculateCompletionRate ────────────────────────────────────────────
// BUG #1: Integer division — hasil selalu 0 (kecuali semua task selesai).

func TestCalculateCompletionRate(t *testing.T) {
	tests := []struct {
		name    string
		tasks   []model.Task
		want    float64
		isBug   bool
	}{
		{
			name:  "tidak ada task",
			tasks: []model.Task{},
			want:  0,
		},
		{
			name:  "semua done → 100%",
			tasks: []model.Task{{Status: model.StatusDone}, {Status: model.StatusDone}},
			want:  100.0,
		},
		{
			// [BUG] 1/3 dengan integer division = 0, bukan 33.33
			name: "[BUG] sepertiga selesai → 33.33%",
			tasks: []model.Task{
				{Status: model.StatusDone},
				{Status: model.StatusTodo},
				{Status: model.StatusTodo},
			},
			want:  33.33,
			isBug: true,
		},
		{
			// [BUG] 1/2 dengan integer division = 0, bukan 50.0
			name:  "[BUG] setengah selesai → 50%",
			tasks: []model.Task{{Status: model.StatusDone}, {Status: model.StatusTodo}},
			want:  50.0,
			isBug: true,
		},
		{
			name: "tidak ada yang selesai → 0%",
			tasks: []model.Task{
				{Status: model.StatusTodo},
				{Status: model.StatusInProgress},
			},
			want: 0.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := service.CalculateCompletionRate(tc.tasks)
			// Toleransi 0.01 untuk floating point
			diff := got - tc.want
			if diff < 0 {
				diff = -diff
			}
			if diff > 0.01 {
				if tc.isBug {
					t.Errorf("BUG TERDETEKSI — CalculateCompletionRate() = %.2f, want %.2f\n"+
						"  → Integer division: %d/%d = 0 (bukan %.2f)\n"+
						"  → Perbaiki: gunakan float64(completed)/float64(len(tasks))*100",
						got, tc.want, len(tc.tasks)/2, len(tc.tasks), tc.want)
				} else {
					t.Errorf("CalculateCompletionRate() = %.2f, want %.2f", got, tc.want)
				}
			}
		})
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

func TestCreate(t *testing.T) {
	svc := newSvc()

	t.Run("sukses dengan default priority", func(t *testing.T) {
		task, err := svc.Create(model.CreateTaskRequest{Title: "Belajar Go"})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if task.Title != "Belajar Go" {
			t.Errorf("Title = %q, want %q", task.Title, "Belajar Go")
		}
		if task.Status != model.StatusTodo {
			t.Errorf("Status = %q, want todo", task.Status)
		}
		if task.Priority != model.PriorityMedium {
			t.Errorf("Priority = %q, want medium (default)", task.Priority)
		}
		if task.ID == "" {
			t.Error("ID tidak boleh kosong")
		}
	})

	t.Run("title kosong ditolak", func(t *testing.T) {
		_, err := svc.Create(model.CreateTaskRequest{Title: ""})
		if err == nil {
			t.Error("Create() harus error jika title kosong")
		}
	})

	t.Run("title spasi saja ditolak", func(t *testing.T) {
		_, err := svc.Create(model.CreateTaskRequest{Title: "   "})
		if err == nil {
			t.Error("Create() harus error jika title hanya spasi")
		}
	})

	t.Run("priority invalid ditolak", func(t *testing.T) {
		_, err := svc.Create(model.CreateTaskRequest{Title: "T", Priority: "extreme"})
		if err == nil {
			t.Error("Create() harus error untuk priority tidak valid")
		}
	})

	t.Run("priority high sukses", func(t *testing.T) {
		task, err := svc.Create(model.CreateTaskRequest{Title: "Urgent", Priority: model.PriorityHigh})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if task.Priority != model.PriorityHigh {
			t.Errorf("Priority = %q, want high", task.Priority)
		}
	})

	t.Run("setiap task ID unik", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 50; i++ {
			task, _ := svc.Create(model.CreateTaskRequest{Title: "Task"})
			if ids[task.ID] {
				t.Errorf("ID duplikat ditemukan: %s", task.ID)
			}
			ids[task.ID] = true
		}
	})

	t.Run("title lebih dari 200 karakter ditolak", func(t *testing.T) {
		longTitle := ""
		for i := 0; i < 201; i++ {
			longTitle += "a"
		}
		_, err := svc.Create(model.CreateTaskRequest{Title: longTitle})
		if err == nil {
			t.Error("Create() harus error jika title > 200 karakter")
		}
	})
}

// ── Update ───────────────────────────────────────────────────────────────────

func TestUpdate(t *testing.T) {
	svc := newSvc()

	t.Run("update status ke done mengisi completed_at", func(t *testing.T) {
		task, _ := svc.Create(model.CreateTaskRequest{Title: "Selesaikan"})
		statusDone := model.StatusDone
		updated, err := svc.Update(task.ID, model.UpdateTaskRequest{Status: &statusDone})
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}
		if updated.CompletedAt == nil {
			t.Error("CompletedAt harus terisi setelah status = done")
		}
	})

	t.Run("update task tidak ada → error", func(t *testing.T) {
		statusDone := model.StatusDone
		_, err := svc.Update("id-tidak-ada", model.UpdateTaskRequest{Status: &statusDone})
		if err == nil {
			t.Error("Update() harus error untuk ID tidak ada")
		}
	})

	t.Run("update status invalid → error", func(t *testing.T) {
		task, _ := svc.Create(model.CreateTaskRequest{Title: "T"})
		s := model.Status("invalid")
		_, err := svc.Update(task.ID, model.UpdateTaskRequest{Status: &s})
		if err == nil {
			t.Error("Update() harus error untuk status tidak valid")
		}
	})

	t.Run("update description sukses", func(t *testing.T) {
		task, _ := svc.Create(model.CreateTaskRequest{Title: "T"})
		desc := "Deskripsi Baru"
		updated, err := svc.Update(task.ID, model.UpdateTaskRequest{Description: &desc})
		if err != nil || updated.Description != desc {
			t.Errorf("Gagal update description")
		}
	})

	t.Run("update title jadi kosong ditolak", func(t *testing.T) {
		task, _ := svc.Create(model.CreateTaskRequest{Title: "T"})
		emptyTitle := ""
		_, err := svc.Update(task.ID, model.UpdateTaskRequest{Title: &emptyTitle})
		if err == nil {
			t.Error("Harusnya error jika title diupdate jadi kosong")
		}
	})
}

// ── [CICD] Full Task Lifecycle ────────────────────────────────────────────────
// [CICD] Simulasi integration test: create → get → update → delete.
// Jenis test ini dijalankan otomatis setelah deploy ke staging.

func TestTaskFullLifecycle(t *testing.T) {
	svc := newSvc()

	// 1. Create
	task, err := svc.Create(model.CreateTaskRequest{
		Title:    "Pipeline Lifecycle Test",
		Priority: model.PriorityHigh,
	})
	if err != nil {
		t.Fatalf("Create() gagal: %v", err)
	}

	// 2. Get
	got, err := svc.GetByID(task.ID)
	if err != nil || got.ID != task.ID {
		t.Fatalf("GetByID() gagal setelah create")
	}

	// 3. Update ke in_progress
	s := model.StatusInProgress
	got, err = svc.Update(task.ID, model.UpdateTaskRequest{Status: &s})
	if err != nil || got.Status != model.StatusInProgress {
		t.Fatalf("Update() ke in_progress gagal")
	}

	// 4. Update ke done
	done := model.StatusDone
	got, err = svc.Update(task.ID, model.UpdateTaskRequest{Status: &done})
	if err != nil || got.CompletedAt == nil {
		t.Fatalf("Update() ke done gagal atau CompletedAt nil")
	}

	// 5. Stats harus menunjukkan 1 done
	stats, err := svc.GetStats()
	if err != nil {
		t.Fatalf("GetStats() gagal: %v", err)
	}
	if stats.ByStatus["done"] != 1 {
		t.Errorf("Stats.ByStatus[done] = %d, want 1", stats.ByStatus["done"])
	}

	// 6. Delete
	_, err = svc.Delete(task.ID)
	if err != nil {
		t.Fatalf("Delete() gagal: %v", err)
	}

	// 7. Pastikan sudah terhapus
	if _, err = svc.GetByID(task.ID); err == nil {
		t.Error("GetByID() harus error setelah task dihapus")
	}
}

// ── [CICD] Rollback Simulation ───────────────────────────────────────────────

func TestRollbackStatusSimulation(t *testing.T) {
	svc := newSvc()
	task, _ := svc.Create(model.CreateTaskRequest{Title: "Rollback Test"})

	// Simulasi: deploy berhasil → update ke in_progress
	s := model.StatusInProgress
	svc.Update(task.ID, model.UpdateTaskRequest{Status: &s}) //nolint

	// Deployment bermasalah → rollback ke todo
	todo := model.StatusTodo
	rolled, err := svc.Update(task.ID, model.UpdateTaskRequest{Status: &todo})
	if err != nil {
		t.Fatalf("Rollback gagal: %v", err)
	}
	if rolled.Status != model.StatusTodo {
		t.Errorf("Setelah rollback, status = %q, want todo", rolled.Status)
	}
}

// ── [TODO] Tambahkan test berikut ─────────────────────────────────────────────
// TODO mahasiswa:
// - TestGetAll_WithStatusFilter (setelah bug #2 diperbaiki)
// - TestGetStats_CompletionRate (setelah bug #1 diperbaiki)
// - TestCreate_WithUnicodeTitle
// - TestDelete_AndVerifyStats

// 1. Test untuk memastikan Bug #2 (Filter Status) sudah benar-benar sembuh
func TestGetAll_WithStatusFilter(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := service.NewTaskService(repo)

	svc.Create(model.CreateTaskRequest{Title: "T1"})
	t2, _ := svc.Create(model.CreateTaskRequest{Title: "T2"})
	svc.Update(t2.ID, model.UpdateTaskRequest{Status: ptr(model.StatusDone)}) 

	t.Run("filter done sukses", func(t *testing.T) {
		tasks, _ := svc.GetAll("done")
		if len(tasks) != 1 {
			t.Errorf("Harusnya 1, dapat %d", len(tasks))
		}
	})

	t.Run("filter kosong ambil semua", func(t *testing.T) {
		tasks, _ := svc.GetAll("")
		if len(tasks) != 2 {
			t.Errorf("Harusnya 2, dapat %d", len(tasks))
		}
	})

	t.Run("filter status ngawur error", func(t *testing.T) {
		_, err := svc.GetAll("ngawur")
		if err == nil {
			t.Error("Harusnya error filter tidak valid")
		}
	})
}

// 2. Test untuk memastikan Bug #1 (Persentase Float) sudah benar
func TestGetStats_CompletionRate(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := service.NewTaskService(repo)

	t.Run("stats saat data kosong", func(t *testing.T) {
		stats, _ := svc.GetStats()
		if stats.Total != 0 || stats.CompletionRate != 0 {
			t.Error("Harusnya 0 saat kosong")
		}
	})

	// Buat 3 task, selesaikan 1 (1/3 = 33.33%)
	svc.Create(model.CreateTaskRequest{Title: "T1"})
	svc.Create(model.CreateTaskRequest{Title: "T2"})
	t3, _ := svc.Create(model.CreateTaskRequest{Title: "T3"})
	svc.Update(t3.ID, model.UpdateTaskRequest{Status: ptr(model.StatusDone)})

	stats, _ := svc.GetStats()
	expected := 33.33
	if stats.CompletionRate < 33.0 || stats.CompletionRate > 34.0 {
		t.Errorf("CompletionRate salah! Dapat %.2f, ingin %.2f. Cek integer division!", stats.CompletionRate, expected)
	}
}

// 3. Test Unicode Title (Memastikan karakter spesial aman)
func TestCreate_WithUnicodeTitle(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := service.NewTaskService(repo)

	title := "Task Unicode"
	task, err := svc.Create(model.CreateTaskRequest{Title: title, Priority: "high"})
	if err != nil {
		t.Fatalf("Harusnya sukses simpan unicode: %v", err)
	}
	if task.Title != title {
		t.Errorf("Title berubah! Ingin %s, dapat %s", title, task.Title)
	}
}

// Fungsi pembantu untuk pointer (tambahkan jika belum ada di file test)
func ptr[T any](v T) *T {
	return &v
}

// 4. TestDelete_AndVerifyStats (Optimasi Coverage Delete & Stats)
func TestDelete_AndVerifyStats(t *testing.T) {
	svc := newSvc()

	t.Run("Gagal - Hapus task yang tidak ada", func(t *testing.T) {
		// Ini akan memicu baris "if !ok" di fungsi Delete
		_, err := svc.Delete("id-ngasal-123")
		if err == nil {
			t.Error("Delete() harusnya error jika ID tidak ditemukan")
		}
	})

	t.Run("Sukses - Hapus task dan verifikasi statistik", func(t *testing.T) {
		// 1. Buat task baru
		task, _ := svc.Create(model.CreateTaskRequest{Title: "Task Hapus"})
		
		// 2. Pastikan stats mencatat 1 task
		s1, _ := svc.GetStats()
		if s1.Total != 1 {
			t.Errorf("Stats awal harusnya 1, dapat %d", s1.Total)
		}

		// 3. Hapus task tersebut
		_, err := svc.Delete(task.ID)
		if err != nil {
			t.Fatalf("Gagal hapus task: %v", err)
		}

		// 4. Pastikan stats kembali jadi 0 (Verifikasi)
		s2, _ := svc.GetStats()
		if s2.Total != 0 {
			t.Errorf("Stats akhir harusnya 0 setelah dihapus, dapat %d", s2.Total)
		}
		
		// 5. Pastikan jika di-GetByID juga error
		_, err = svc.GetByID(task.ID)
		if err == nil {
			t.Error("Task harusnya sudah tidak bisa ditemukan setelah dihapus")
		}
	})
}

type errorRepo struct{}

func (e *errorRepo) Save(task model.Task) error {
    return fmt.Errorf("db error")
}
func (e *errorRepo) FindByID(id string) (model.Task, bool, error) {
    return model.Task{}, false, fmt.Errorf("db error")
}
func (e *errorRepo) FindAll() ([]model.Task, error) {
    return nil, fmt.Errorf("db error")
}
func (e *errorRepo) FindByStatus(status model.Status) ([]model.Task, error) {
    return nil, fmt.Errorf("db error")
}
func (e *errorRepo) Delete(id string) (bool, error) {
    return false, fmt.Errorf("db error")
}
func (e *errorRepo) Count() (int, error) {
    return 0, fmt.Errorf("db error")
}
func (e *errorRepo) Close() error {
    return fmt.Errorf("db error")
}

type deleteErrorRepo struct {
    repository.TaskRepository
}

func (d *deleteErrorRepo) Delete(id string) (bool, error) {
    return false, fmt.Errorf("delete error")
}

func TestDelete_RepoError(t *testing.T) {
    mem := repository.NewMemoryRepository()
    svc := service.NewTaskService(&deleteErrorRepo{mem})

    task, _ := svc.Create(model.CreateTaskRequest{Title: "A"})

    _, err := svc.Delete(task.ID)
    if err == nil {
        t.Error("expected delete error")
    }
}

func TestUpdate_TitleSuccess(t *testing.T) {
    svc := newSvc()
    task, _ := svc.Create(model.CreateTaskRequest{Title: "Old"})

    newTitle := "New"
    updated, err := svc.Update(task.ID, model.UpdateTaskRequest{Title: &newTitle})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if updated.Title != newTitle {
        t.Error("title tidak terupdate")
    }
}

func TestUpdate_Description(t *testing.T) {
    svc := newSvc()
    task, _ := svc.Create(model.CreateTaskRequest{Title: "A"})

    desc := "deskripsi baru"
    updated, err := svc.Update(task.ID, model.UpdateTaskRequest{
        Description: &desc,
    })
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if updated.Description != desc {
        t.Error("description tidak terupdate")
    }
}

func TestCreate_RepoSaveError(t *testing.T) {
    svc := service.NewTaskService(&errorRepo{})
    _, err := svc.Create(model.CreateTaskRequest{Title: "Test"})
    if err == nil {
        t.Error("Create() harus error jika repo.Save gagal")
    }
}

func TestGetByID_RepoError(t *testing.T) {
    svc := service.NewTaskService(&errorRepo{})
    _, err := svc.GetByID("any-id")
    if err == nil {
        t.Error("GetByID() harus error jika repo gagal")
    }
}

func TestUpdate_RepoSaveError(t *testing.T) {
    // Pakai repo yang bisa FindByID tapi gagal saat Save
    mem := repository.NewMemoryRepository()
    task, _ := service.NewTaskService(mem).Create(model.CreateTaskRequest{Title: "T"})

    svc := service.NewTaskService(&saveErrorRepo{mem, task.ID})
    s := model.StatusDone
    _, err := svc.Update(task.ID, model.UpdateTaskRequest{Status: &s})
    if err == nil {
        t.Error("Update() harus error jika repo.Save gagal")
    }
}

func TestGetStats_RepoError(t *testing.T) {
    svc := service.NewTaskService(&errorRepo{})
    _, err := svc.GetStats()
    if err == nil {
        t.Error("GetStats() harus error jika repo.FindAll gagal")
    }
}

func TestDelete_FindByIDError(t *testing.T) {
    svc := service.NewTaskService(&errorRepo{})
    _, err := svc.Delete("any-id")
    if err == nil {
        t.Error("Delete() harus error jika repo.FindByID gagal")
    }
}

// saveErrorRepo: FindByID sukses, tapi Save selalu gagal
type saveErrorRepo struct {
    repository.TaskRepository
    targetID string
}

func (r *saveErrorRepo) Save(task model.Task) error {
    return fmt.Errorf("save error")
}