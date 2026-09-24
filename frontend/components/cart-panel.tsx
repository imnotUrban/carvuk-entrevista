"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { boletasApi } from "@/lib/api";
import { maxQuantity, useCart } from "@/lib/cart";
import { formatCLP } from "@/lib/format";

export function CartPanel() {
  const { items, total, setQuantity, remove, clear } = useCart();
  const [buying, setBuying] = useState(false);
  const router = useRouter();

  // The backend prices the items and computes tax; the cart is only cleared on success.
  const handleBuy = async () => {
    setBuying(true);
    try {
      const boleta = await boletasApi.create({
        items: items.map((i) => ({ producto_id: i.producto_id, cantidad: i.cantidad })),
      });
      clear();
      toast.success(`Boleta #${boleta.id} generada`);
      router.push(`/boletas/${boleta.id}`);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "No se pudo generar la boleta");
    } finally {
      setBuying(false);
    }
  };

  return (
    <div className="space-y-4 rounded-md border p-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">Carrito</h2>
        <Button variant="ghost" size="sm" disabled={items.length === 0} onClick={clear}>
          Vaciar carrito
        </Button>
      </div>

      {items.length === 0 ? (
        <p className="text-sm text-muted-foreground">El carrito está vacío.</p>
      ) : (
        <ul className="divide-y">
          {items.map((i) => {
            const atMax = i.cantidad >= maxQuantity(i.stock);
            return (
              <li key={i.producto_id} className="space-y-2 py-3">
                <div className="flex justify-between gap-2 text-sm">
                  <span className="font-medium">{i.nombre}</span>
                  <span>{formatCLP(i.precio * i.cantidad)}</span>
                </div>
                <div className="flex items-center justify-between text-sm text-muted-foreground">
                  <span>{formatCLP(i.precio)} c/u</span>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      aria-label={`Disminuir ${i.nombre}`}
                      onClick={() => setQuantity(i.producto_id, i.cantidad - 1)}
                    >
                      −
                    </Button>
                    <span className="w-6 text-center text-foreground">{i.cantidad}</span>
                    <Button
                      variant="outline"
                      size="sm"
                      aria-label={`Aumentar ${i.nombre}`}
                      disabled={atMax}
                      onClick={() => setQuantity(i.producto_id, i.cantidad + 1)}
                    >
                      +
                    </Button>
                    <Button variant="ghost" size="sm" onClick={() => remove(i.producto_id)}>
                      Quitar
                    </Button>
                  </div>
                </div>
                {atMax && <p className="text-xs text-amber-600">Máximo disponible</p>}
              </li>
            );
          })}
        </ul>
      )}

      <div className="flex justify-between border-t pt-3 font-semibold">
        <span>Total</span>
        <span>{formatCLP(total)}</span>
      </div>
      <Button
        className="w-full"
        disabled={items.length === 0 || buying}
        onClick={handleBuy}
      >
        {buying ? "Procesando..." : "Generar boleta"}
      </Button>
    </div>
  );
}
