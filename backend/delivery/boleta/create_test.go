package boleta_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	delivery "backend/delivery/boleta"
	productoentity "backend/entity/producto"
	"backend/pkg/testutil"
	boletarepo "backend/repository/boleta"
	productorepo "backend/repository/producto"
	service "backend/service/boleta"
)

// El request no declara campos de impuesto: enviarlos no debe alterar el resultado.
func TestCreateIgnoresClientTaxFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupDB(t)
	prepo := productorepo.NewRepository(db)
	svc := service.NewService(boletarepo.NewRepository(db), prepo)
	if err := svc.SeedUsuarioIfEmpty(); err != nil {
		t.Fatalf("seed usuario: %v", err)
	}
	p := &productoentity.Producto{Nombre: "Leche", Precio: 4200, Stock: 5}
	if err := db.Create(p).Error; err != nil {
		t.Fatalf("create producto: %v", err)
	}

	r := gin.New()
	delivery.NewHandler(svc).RegisterRoutes(r.Group("/api/v1"))

	body := `{"items":[{"producto_id":1,"cantidad":1}],
		"porcentaje_impuesto":19,"impuesto":1,"valor_neto":1,"valor_bruto":1}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/boletas", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	var got struct {
		ValorBruto         int `json:"valor_bruto"`
		Impuesto           int `json:"impuesto"`
		ValorNeto          int `json:"valor_neto"`
		PorcentajeImpuesto int `json:"porcentaje_impuesto"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ValorBruto != 4200 || got.Impuesto != 630 || got.ValorNeto != 3570 || got.PorcentajeImpuesto != 15 {
		t.Fatalf("client fields must be ignored: %+v", got)
	}
}
