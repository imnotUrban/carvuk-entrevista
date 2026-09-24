package boleta_test

import (
	"errors"
	"testing"

	"gorm.io/gorm"

	boletaentity "backend/entity/boleta"
	compraentity "backend/entity/compra"
	productoentity "backend/entity/producto"
	"backend/pkg/apperr"
	"backend/pkg/tax"
	"backend/pkg/testutil"
	boletarepo "backend/repository/boleta"
	productorepo "backend/repository/producto"
	service "backend/service/boleta"
	productoservice "backend/service/producto"
)

type env struct {
	svc       service.Service
	productos productoservice.Service
	db        *gorm.DB
}

func newEnv(t *testing.T) *env {
	t.Helper()
	db := testutil.SetupDB(t)
	prepo := productorepo.NewRepository(db)
	svc := service.NewService(boletarepo.NewRepository(db), prepo)
	if err := svc.SeedUsuarioIfEmpty(); err != nil {
		t.Fatalf("seed usuario: %v", err)
	}
	return &env{svc: svc, productos: productoservice.NewService(prepo), db: db}
}

func (e *env) producto(t *testing.T, nombre string, precio, stock int) *productoentity.Producto {
	t.Helper()
	p, err := e.productos.Create(productoservice.CreateInput{Nombre: nombre, Precio: precio, Stock: stock})
	if err != nil {
		t.Fatalf("create producto: %v", err)
	}
	return p
}

func (e *env) stock(t *testing.T, id uint) int {
	t.Helper()
	p, err := e.productos.GetByID(id)
	if err != nil {
		t.Fatalf("get producto: %v", err)
	}
	return p.Stock
}

func (e *env) count(t *testing.T, model any) int64 {
	t.Helper()
	var n int64
	if err := e.db.Model(model).Count(&n).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

func (e *env) assertNothingCreated(t *testing.T) {
	t.Helper()
	if e.count(t, &compraentity.Compra{}) != 0 || e.count(t, &compraentity.CompraItem{}) != 0 || e.count(t, &boletaentity.Boleta{}) != 0 {
		t.Fatal("expected no compra/items/boleta to be created")
	}
}

func input(pairs ...int) service.CreateInput {
	in := service.CreateInput{}
	for i := 0; i < len(pairs); i += 2 {
		in.Items = append(in.Items, service.ItemInput{ProductoID: uint(pairs[i]), Cantidad: pairs[i+1]})
	}
	return in
}

func TestCreateExample(t *testing.T) {
	e := newEnv(t)
	leche := e.producto(t, "Leche", 1100, 10)
	aceite := e.producto(t, "Aceite", 2000, 10)

	b, err := e.svc.Create(input(int(leche.ID), 2, int(aceite.ID), 1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.ValorBruto != 4200 || b.Impuesto != 630 || b.ValorNeto != 3570 || b.PorcentajeImpuesto != 15 {
		t.Fatalf("unexpected totals: %+v", b)
	}
	if b.Usuario == nil || b.Usuario.Correo != "cliente.demo@example.com" || len(b.Items) != 2 {
		t.Fatalf("unexpected detail: %+v", b)
	}

	var stored boletaentity.Boleta
	if err := e.db.First(&stored, b.ID).Error; err != nil {
		t.Fatalf("boleta not stored: %v", err)
	}
	if stored.CompraID != b.CompraID || stored.CompraID == 0 {
		t.Fatalf("boleta.id_compra mismatch: %+v", stored)
	}
	var items []compraentity.CompraItem
	e.db.Where("compra_id = ?", stored.CompraID).Find(&items)
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
}

func TestCreateRounding(t *testing.T) {
	e := newEnv(t)
	p := e.producto(t, "Cosa", 1001, 5)

	b, err := e.svc.Create(input(int(p.ID), 1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.ValorBruto != 1001 || b.Impuesto != 150 || b.ValorNeto != 851 {
		t.Fatalf("unexpected totals: %+v", b)
	}
}

func TestCreateSameBrutoSameTax(t *testing.T) {
	e := newEnv(t)
	a := e.producto(t, "A", 1000, 10)
	b := e.producto(t, "B", 400, 10)
	c := e.producto(t, "C", 2000, 10)

	b1, err := e.svc.Create(input(int(a.ID), 2, int(b.ID), 1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b2, err := e.svc.Create(input(int(c.ID), 1, int(b.ID), 1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b1.ValorBruto != b2.ValorBruto || b1.Impuesto != b2.Impuesto || b1.ValorNeto != b2.ValorNeto {
		t.Fatalf("same bruto must give same tax: %+v vs %+v", b1, b2)
	}
	if b1.PorcentajeImpuesto != tax.TaxRatePercent || b2.PorcentajeImpuesto != tax.TaxRatePercent {
		t.Fatalf("porcentaje_impuesto must be 15: %d, %d", b1.PorcentajeImpuesto, b2.PorcentajeImpuesto)
	}
}

func TestCreateValidation(t *testing.T) {
	e := newEnv(t)
	p := e.producto(t, "Leche", 1100, 200)
	id := int(p.ID)

	cases := map[string]service.CreateInput{
		"empty":        {},
		"zero qty":     input(id, 0),
		"negative qty": input(id, -1),
		"qty over 99":  input(id, 100),
		"duplicated":   input(id, 1, id, 2),
	}
	for name, in := range cases {
		if _, err := e.svc.Create(in); !errors.Is(err, apperr.ErrValidation) {
			t.Errorf("%s: expected ErrValidation, got %v", name, err)
		}
	}

	var tooMany service.CreateInput
	for i := 1; i <= 101; i++ {
		tooMany.Items = append(tooMany.Items, service.ItemInput{ProductoID: uint(i), Cantidad: 1})
	}
	if _, err := e.svc.Create(tooMany); !errors.Is(err, apperr.ErrValidation) {
		t.Errorf("101 items: expected ErrValidation, got %v", err)
	}

	e.assertNothingCreated(t)
	if got := e.stock(t, p.ID); got != 200 {
		t.Fatalf("stock changed to %d", got)
	}
}

func TestCreateProductoNotFound(t *testing.T) {
	e := newEnv(t)
	ok := e.producto(t, "Leche", 1100, 10)
	deleted := e.producto(t, "Viejo", 500, 10)
	if err := e.productos.Delete(deleted.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	for name, in := range map[string]service.CreateInput{
		"missing": input(int(ok.ID), 1, 9999, 1),
		"deleted": input(int(ok.ID), 1, int(deleted.ID), 1),
	} {
		if _, err := e.svc.Create(in); !errors.Is(err, apperr.ErrNotFound) {
			t.Errorf("%s: expected ErrNotFound, got %v", name, err)
		}
	}

	e.assertNothingCreated(t)
	if got := e.stock(t, ok.ID); got != 10 {
		t.Fatalf("stock changed to %d", got)
	}
}

func TestCreateDecrementsStock(t *testing.T) {
	e := newEnv(t)
	a := e.producto(t, "A", 100, 10)
	b := e.producto(t, "B", 200, 5)

	if _, err := e.svc.Create(input(int(a.ID), 3, int(b.ID), 2)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := e.stock(t, a.ID); got != 7 {
		t.Errorf("A stock = %d, want 7", got)
	}
	if got := e.stock(t, b.ID); got != 3 {
		t.Errorf("B stock = %d, want 3", got)
	}
}

func TestCreateInsufficientStockRollsBack(t *testing.T) {
	e := newEnv(t)
	a := e.producto(t, "A", 100, 10)
	b := e.producto(t, "B", 200, 1)

	// A alcanza, B no: nada debe cambiar.
	_, err := e.svc.Create(input(int(a.ID), 3, int(b.ID), 2))
	if !errors.Is(err, apperr.ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}
	e.assertNothingCreated(t)
	if got := e.stock(t, a.ID); got != 10 {
		t.Errorf("A stock = %d, want 10", got)
	}
	if got := e.stock(t, b.ID); got != 1 {
		t.Errorf("B stock = %d, want 1", got)
	}
}

func TestCreateRepositoryGuardRollsBack(t *testing.T) {
	// Simula una venta concurrente: el stock baja entre la lectura del service y el descuento.
	e := newEnv(t)
	a := e.producto(t, "A", 100, 10)
	b := e.producto(t, "B", 200, 5)
	repo := boletarepo.NewRepository(e.db)

	e.db.Model(&productoentity.Producto{}).Where("id = ?", b.ID).Update("stock", 0)

	compra := &compraentity.Compra{UserID: 1, Items: []compraentity.CompraItem{
		{ProductoID: a.ID, Nombre: "A", PrecioUnitario: 100, Cantidad: 3},
		{ProductoID: b.ID, Nombre: "B", PrecioUnitario: 200, Cantidad: 2},
	}}
	err := repo.Create(compra, &boletaentity.Boleta{})
	if !errors.Is(err, apperr.ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}
	e.assertNothingCreated(t)
	if got := e.stock(t, a.ID); got != 10 {
		t.Errorf("A stock = %d, want 10 (rolled back)", got)
	}
}

func TestCreateExactStock(t *testing.T) {
	e := newEnv(t)
	p := e.producto(t, "Leche", 1100, 3)

	if _, err := e.svc.Create(input(int(p.ID), 3)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := e.stock(t, p.ID); got != 0 {
		t.Fatalf("stock = %d, want 0", got)
	}
	if _, err := e.svc.Create(input(int(p.ID), 1)); !errors.Is(err, apperr.ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}
}

func TestSnapshotSurvivesProductChanges(t *testing.T) {
	e := newEnv(t)
	p := e.producto(t, "Leche", 1100, 10)

	b, err := e.svc.Create(input(int(p.ID), 2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newName, newPrice := "Leche Premium", 5000
	if _, err := e.productos.Update(p.ID, productoservice.UpdateInput{Nombre: &newName, Precio: &newPrice}); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := e.productos.Delete(p.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	var items []compraentity.CompraItem
	e.db.Where("compra_id = ?", b.CompraID).Find(&items)
	if len(items) != 1 || items[0].Nombre != "Leche" || items[0].PrecioUnitario != 1100 {
		t.Fatalf("snapshot changed: %+v", items)
	}
}

func TestSeedUsuarioIsIdempotent(t *testing.T) {
	e := newEnv(t)
	if err := e.svc.SeedUsuarioIfEmpty(); err != nil {
		t.Fatalf("second seed: %v", err)
	}
	var n int64
	e.db.Table("usuarios").Count(&n)
	if n != 1 {
		t.Fatalf("expected 1 usuario, got %d", n)
	}
}

func TestListNewestFirstWithUsuario(t *testing.T) {
	e := newEnv(t)
	p := e.producto(t, "Leche", 1100, 50)
	first, _ := e.svc.Create(input(int(p.ID), 1))
	second, err := e.svc.Create(input(int(p.ID), 3))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, total, err := e.svc.List(service.ListInput{})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 2 || len(items) != 2 || items[0].ID != second.ID || items[1].ID != first.ID {
		t.Fatalf("expected newest first, got %+v (total %d)", items, total)
	}
	if items[0].CantidadItems != 3 || items[0].ValorBruto != 3300 {
		t.Fatalf("unexpected summary: %+v", items[0])
	}
	u := items[0].Usuario
	if u == nil || u.ID == 0 || u.Nombre != "Cliente Demo" || u.Correo != "cliente.demo@example.com" {
		t.Fatalf("expected usuario in list item, got %+v", u)
	}
}

func TestListPagination(t *testing.T) {
	e := newEnv(t)
	p := e.producto(t, "Leche", 1100, 50)
	for i := 0; i < 5; i++ {
		if _, err := e.svc.Create(input(int(p.ID), 1)); err != nil {
			t.Fatalf("create: %v", err)
		}
	}

	items, total, err := e.svc.List(service.ListInput{Page: 2, Limit: 2})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 5 || len(items) != 2 || items[0].ID != 3 || items[1].ID != 2 {
		t.Fatalf("unexpected page: %+v (total %d)", items, total)
	}

	items, _, err = e.svc.List(service.ListInput{Limit: 1000})
	if err != nil || len(items) != 5 {
		t.Fatalf("limit >100 should be capped, not fail: %v", err)
	}

	if _, _, err := e.svc.List(service.ListInput{Order: "sideways"}); !errors.Is(err, apperr.ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestGetByIDExample(t *testing.T) {
	e := newEnv(t)
	a := e.producto(t, "A", 1000, 10)
	b := e.producto(t, "B", 2200, 10)
	created, err := e.svc.Create(input(int(a.ID), 2, int(b.ID), 1))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	d, err := e.svc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if d.ValorNeto != 3570 || d.Impuesto != 630 || d.ValorBruto != 4200 || d.PorcentajeImpuesto != tax.TaxRatePercent {
		t.Fatalf("unexpected totals: %+v", d)
	}
	if d.ValorNeto+d.Impuesto != d.ValorBruto {
		t.Fatal("neto + impuesto != bruto")
	}
	if d.Usuario == nil || d.Usuario.Correo != "cliente.demo@example.com" {
		t.Fatalf("expected usuario, got %+v", d.Usuario)
	}
	sum := 0
	for _, it := range d.Items {
		if it.Subtotal != it.PrecioUnitario*it.Cantidad || it.Nombre == "" {
			t.Fatalf("bad item: %+v", it)
		}
		sum += it.Subtotal
	}
	if len(d.Items) != 2 || sum != d.ValorBruto {
		t.Fatalf("items %+v sum %d", d.Items, sum)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	e := newEnv(t)
	if _, err := e.svc.GetByID(999); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetByIDKeepsSnapshotAfterProductChanges(t *testing.T) {
	e := newEnv(t)
	p := e.producto(t, "Leche", 1100, 10)
	created, _ := e.svc.Create(input(int(p.ID), 2))

	newName, newPrice := "Leche Premium", 5000
	if _, err := e.productos.Update(p.ID, productoservice.UpdateInput{Nombre: &newName, Precio: &newPrice}); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := e.productos.Delete(p.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	d, err := e.svc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if d.ValorBruto != 2200 || d.Items[0].Nombre != "Leche" || d.Items[0].PrecioUnitario != 1100 {
		t.Fatalf("snapshot changed: %+v", d)
	}
}
