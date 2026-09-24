package producto_test

import (
	"errors"
	"testing"

	"gorm.io/gorm"

	entity "backend/entity/producto"
	"backend/pkg/apperr"
	"backend/pkg/testutil"
	repo "backend/repository/producto"
	service "backend/service/producto"
)

func newService(t *testing.T) (service.Service, *gorm.DB) {
	t.Helper()
	db := testutil.SetupDB(t)
	return service.NewService(repo.NewRepository(db)), db
}

func mustCreate(t *testing.T, svc service.Service, nombre string, precio, stock int) *entity.Producto {
	t.Helper()
	item, err := svc.Create(service.CreateInput{Nombre: nombre, Precio: precio, Stock: stock})
	if err != nil {
		t.Fatalf("create %q: %v", nombre, err)
	}
	return item
}

func TestCreate(t *testing.T) {
	svc, _ := newService(t)

	item, err := svc.Create(service.CreateInput{Nombre: "  Leche ", Precio: 1100, Stock: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID == 0 || item.Nombre != "Leche" || item.Precio != 1100 || item.Stock != 5 {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestCreateStockDefaultsToZero(t *testing.T) {
	svc, _ := newService(t)
	item := mustCreate(t, svc, "Pan", 1800, 0)
	if item.Stock != 0 {
		t.Fatalf("expected stock 0, got %d", item.Stock)
	}
}

func TestCreateValidation(t *testing.T) {
	svc, _ := newService(t)

	cases := map[string]service.CreateInput{
		"empty nombre":    {Precio: 100},
		"blank nombre":    {Nombre: "   ", Precio: 100},
		"zero precio":     {Nombre: "A", Precio: 0},
		"negative precio": {Nombre: "A", Precio: -5},
		"negative stock":  {Nombre: "A", Precio: 100, Stock: -1},
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := svc.Create(input); !errors.Is(err, apperr.ErrValidation) {
				t.Fatalf("expected ErrValidation, got %v", err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	svc, _ := newService(t)
	created := mustCreate(t, svc, "Leche", 1100, 1)

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

func TestListFilterPaginationSort(t *testing.T) {
	svc, _ := newService(t)
	mustCreate(t, svc, "Leche", 1100, 1)
	mustCreate(t, svc, "Leche Sin Lactosa", 1500, 2)
	mustCreate(t, svc, "Pan", 1800, 3)

	items, total, err := svc.List(service.ListInput{Nombre: "LECHE"})
	if err != nil || total != 2 || len(items) != 2 {
		t.Fatalf("filter nombre: total=%d len=%d err=%v", total, len(items), err)
	}

	items, total, err = svc.List(service.ListInput{Page: 2, Limit: 2})
	if err != nil || total != 3 || len(items) != 1 {
		t.Fatalf("pagination: total=%d len=%d err=%v", total, len(items), err)
	}

	items, _, err = svc.List(service.ListInput{SortBy: "precio", Order: "desc"})
	if err != nil || items[0].Nombre != "Pan" {
		t.Fatalf("sort: %+v err=%v", items, err)
	}
}

func TestListInvalidSortBy(t *testing.T) {
	svc, _ := newService(t)
	if _, _, err := svc.List(service.ListInput{SortBy: "deleted_at; DROP TABLE productos"}); !errors.Is(err, apperr.ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestUpdatePartial(t *testing.T) {
	svc, _ := newService(t)
	created := mustCreate(t, svc, "Leche", 1100, 10)

	stock := 42
	updated, err := svc.Update(created.ID, service.UpdateInput{Stock: &stock})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Stock != 42 || updated.Nombre != "Leche" || updated.Precio != 1100 {
		t.Fatalf("unexpected item: %+v", updated)
	}

	nombre := "Leche Entera"
	updated, err = svc.Update(created.ID, service.UpdateInput{Nombre: &nombre})
	if err != nil || updated.Nombre != "Leche Entera" || updated.Stock != 42 {
		t.Fatalf("unexpected item: %+v, %v", updated, err)
	}
}

func TestUpdateValidation(t *testing.T) {
	svc, _ := newService(t)
	created := mustCreate(t, svc, "Leche", 1100, 10)

	empty, zero, neg := "  ", 0, -1
	cases := map[string]service.UpdateInput{
		"empty nombre":    {Nombre: &empty},
		"zero precio":     {Precio: &zero},
		"negative stock":  {Stock: &neg},
		"negative precio": {Precio: &neg},
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := svc.Update(created.ID, input); !errors.Is(err, apperr.ErrValidation) {
				t.Fatalf("expected ErrValidation, got %v", err)
			}
		})
	}
}

func TestUpdateNotFound(t *testing.T) {
	svc, _ := newService(t)
	stock := 1
	if _, err := svc.Update(999, service.UpdateInput{Stock: &stock}); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSoftDelete(t *testing.T) {
	svc, db := newService(t)
	created := mustCreate(t, svc, "Leche", 1100, 1)

	if err := svc.Delete(created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := svc.GetByID(created.ID); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
	if _, total, _ := svc.List(service.ListInput{}); total != 0 {
		t.Fatalf("expected empty list, got total %d", total)
	}

	var count int64
	db.Unscoped().Model(&entity.Producto{}).Where("id = ?", created.ID).Count(&count)
	if count != 1 {
		t.Fatalf("expected row to remain in DB, got %d", count)
	}
}

func TestDeleteNotFound(t *testing.T) {
	svc, _ := newService(t)
	if err := svc.Delete(999); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSeedIfEmpty(t *testing.T) {
	svc, _ := newService(t)

	if err := svc.SeedIfEmpty(); err != nil {
		t.Fatalf("seed: %v", err)
	}
	items, total, _ := svc.List(service.ListInput{})
	if total != 3 || items[0].Nombre != "Leche" || items[0].Stock != 50 {
		t.Fatalf("unexpected seed: %+v", items)
	}

	if err := svc.SeedIfEmpty(); err != nil {
		t.Fatalf("second seed: %v", err)
	}
	if _, total, _ := svc.List(service.ListInput{}); total != 3 {
		t.Fatalf("seed duplicated rows: total %d", total)
	}
}
