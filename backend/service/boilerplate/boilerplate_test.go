package boilerplate_test

import (
	"errors"
	"testing"

	"gorm.io/gorm"

	entity "backend/entity/boilerplate"
	"backend/pkg/apperr"
	"backend/pkg/testutil"
	repo "backend/repository/boilerplate"
	service "backend/service/boilerplate"
)

func newService(t *testing.T) (service.Service, *gorm.DB) {
	t.Helper()
	db := testutil.SetupDB(t)
	return service.NewService(repo.NewRepository(db)), db
}

func mustCreate(t *testing.T, svc service.Service, name, code string, amount float64) *entity.Boilerplate {
	t.Helper()
	item, err := svc.Create(service.CreateInput{Name: name, Code: code, Amount: amount})
	if err != nil {
		t.Fatalf("create %q: %v", name, err)
	}
	return item
}

func TestCreate(t *testing.T) {
	svc, _ := newService(t)

	item, err := svc.Create(service.CreateInput{Name: "  Widget ", Code: "W-1", Quantity: 3, Amount: 9.5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID == 0 || item.Name != "Widget" || item.Code != "W-1" {
		t.Fatalf("unexpected item: %+v", item)
	}
	if item.Status != entity.StatusDraft {
		t.Fatalf("expected default status draft, got %q", item.Status)
	}
}

func TestCreateValidation(t *testing.T) {
	svc, _ := newService(t)

	cases := map[string]service.CreateInput{
		"empty name":        {Code: "A"},
		"empty code":        {Name: "A"},
		"invalid status":    {Name: "A", Code: "A", Status: "nope"},
		"negative quantity": {Name: "A", Code: "A", Quantity: -1},
		"negative amount":   {Name: "A", Code: "A", Amount: -1},
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := svc.Create(input); !errors.Is(err, apperr.ErrValidation) {
				t.Fatalf("expected ErrValidation, got %v", err)
			}
		})
	}
}

func TestCreateDuplicateCode(t *testing.T) {
	svc, _ := newService(t)
	mustCreate(t, svc, "One", "DUP", 1)

	_, err := svc.Create(service.CreateInput{Name: "Two", Code: "DUP"})
	if !errors.Is(err, apperr.ErrDuplicateCode) {
		t.Fatalf("expected ErrDuplicateCode, got %v", err)
	}
}

func TestGet(t *testing.T) {
	svc, _ := newService(t)
	created := mustCreate(t, svc, "One", "C1", 1)

	got, err := svc.GetByID(created.ID)
	if err != nil || got.ID != created.ID {
		t.Fatalf("unexpected result: %+v, %v", got, err)
	}
}

func TestGetNotFound(t *testing.T) {
	svc, _ := newService(t)

	if _, err := svc.GetByID(999); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestListFilters(t *testing.T) {
	svc, _ := newService(t)
	mustCreate(t, svc, "Alpha Widget", "A-1", 10)
	mustCreate(t, svc, "Beta Widget", "B-1", 50)
	if _, err := svc.Create(service.CreateInput{Name: "Gamma", Code: "G-1", Status: entity.StatusActive, Amount: 30}); err != nil {
		t.Fatal(err)
	}

	minAmount := 20.0
	maxAmount := 40.0

	tests := []struct {
		name  string
		input service.ListInput
		want  int64
	}{
		{"all", service.ListInput{}, 3},
		{"by name case-insensitive", service.ListInput{Name: "widget"}, 2},
		{"by code", service.ListInput{Code: "g-"}, 1},
		{"by status", service.ListInput{Status: entity.StatusActive}, 1},
		{"amount range", service.ListInput{MinAmount: &minAmount, MaxAmount: &maxAmount}, 1},
		{"combined", service.ListInput{Name: "widget", Status: entity.StatusDraft}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, total, err := svc.List(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if total != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, total)
			}
		})
	}
}

func TestListPagination(t *testing.T) {
	svc, _ := newService(t)
	for _, code := range []string{"P1", "P2", "P3"} {
		mustCreate(t, svc, "Item "+code, code, 1)
	}

	items, total, err := svc.List(service.ListInput{Page: 2, Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if total != 3 || len(items) != 1 {
		t.Fatalf("expected total 3 and 1 item on page 2, got total %d, %d items", total, len(items))
	}
}

func TestUpdate(t *testing.T) {
	svc, _ := newService(t)
	created := mustCreate(t, svc, "Old", "U1", 1)

	name, status, qty, amount := "New", entity.StatusArchived, 7, 12.5
	updated, err := svc.Update(created.ID, service.UpdateInput{Name: &name, Status: &status, Quantity: &qty, Amount: &amount})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != "New" || updated.Status != entity.StatusArchived || updated.Quantity != 7 || updated.Amount != 12.5 {
		t.Fatalf("unexpected update result: %+v", updated)
	}
	if updated.Code != "U1" {
		t.Fatalf("partial update must keep code, got %q", updated.Code)
	}
}

func TestUpdateNotFound(t *testing.T) {
	svc, _ := newService(t)

	name := "x"
	if _, err := svc.Update(999, service.UpdateInput{Name: &name}); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdateRejectsDuplicateCode(t *testing.T) {
	svc, _ := newService(t)
	mustCreate(t, svc, "One", "TAKEN", 1)
	other := mustCreate(t, svc, "Two", "FREE", 1)

	code := "TAKEN"
	if _, err := svc.Update(other.ID, service.UpdateInput{Code: &code}); !errors.Is(err, apperr.ErrDuplicateCode) {
		t.Fatalf("expected ErrDuplicateCode, got %v", err)
	}
}

func TestSoftDelete(t *testing.T) {
	svc, db := newService(t)
	created := mustCreate(t, svc, "Gone", "D1", 1)

	if err := svc.Delete(created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.GetByID(created.ID); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
	if _, total, _ := svc.List(service.ListInput{}); total != 0 {
		t.Fatalf("expected empty list after delete, got %d", total)
	}

	// Row still exists physically (soft delete).
	var count int64
	db.Unscoped().Model(&entity.Boilerplate{}).Where("id = ?", created.ID).Count(&count)
	if count != 1 {
		t.Fatalf("expected row to remain in DB, count=%d", count)
	}
}

func TestDeleteNotFound(t *testing.T) {
	svc, _ := newService(t)

	if err := svc.Delete(999); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
