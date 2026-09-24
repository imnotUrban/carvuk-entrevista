"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { CartPanel } from "@/components/cart-panel";
import { ApiError, productosApi } from "@/lib/api";
import { maxQuantity, useCart } from "@/lib/cart";
import { formatCLP } from "@/lib/format";
import type { Producto } from "@/lib/types";

export default function TiendaPage() {
  const { items: cart, add } = useCart();
  const [productos, setProductos] = useState<Producto[]>([]);
  const [loading, setLoading] = useState(false);

  const loadProductos = useCallback(async () => {
    setLoading(true);
    try {
      const res = await productosApi.list({ limit: 100, sort_by: "id", order: "asc" });
      setProductos(res.data);
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Error al cargar productos");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    // Fetch-on-mount: syncing with the external API.
    // eslint-disable-next-line react-hooks/set-state-in-effect
    loadProductos();
  }, [loadProductos]);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">Tienda</h1>
      <div className="grid gap-6 md:grid-cols-[1fr_360px]">
        <div className="space-y-3">
          {loading ? (
            <p className="text-sm text-muted-foreground">Cargando...</p>
          ) : productos.length === 0 ? (
            <p className="text-sm text-muted-foreground">Sin productos</p>
          ) : (
            productos.map((p) => {
              const inCart = cart.find((i) => i.producto_id === p.id)?.cantidad ?? 0;
              const noStock = p.stock <= 0;
              const atMax = !noStock && inCart >= maxQuantity(p.stock);
              return (
                <div
                  key={p.id}
                  className="flex items-center justify-between rounded-md border p-4"
                >
                  <div>
                    <div className="font-medium">{p.nombre}</div>
                    <div className="text-sm text-muted-foreground">
                      {formatCLP(p.precio)} ·{" "}
                      {noStock ? "Sin stock" : `Stock: ${p.stock}`}
                    </div>
                    {atMax && <div className="text-xs text-amber-600">Máximo disponible</div>}
                  </div>
                  <Button disabled={noStock || atMax} onClick={() => add(p)}>
                    {noStock ? "Sin stock" : "Agregar"}
                  </Button>
                </div>
              );
            })
          )}
        </div>
        <CartPanel />
      </div>
    </div>
  );
}
