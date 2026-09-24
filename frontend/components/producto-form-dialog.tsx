"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { Producto, ProductoInput } from "@/lib/types";

interface ProductoFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  item?: Producto | null;
  onSubmit: (input: ProductoInput) => Promise<void>;
}

const emptyForm: ProductoInput = {
  nombre: "",
  precio: 0,
  stock: 0,
};

export function ProductoFormDialog({
  open,
  onOpenChange,
  item,
  onSubmit,
}: ProductoFormDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        {open && (
          <ProductoForm
            key={item?.id ?? "new"}
            item={item}
            onSubmit={onSubmit}
            onDone={() => onOpenChange(false)}
          />
        )}
      </DialogContent>
    </Dialog>
  );
}

interface ProductoFormProps {
  item?: Producto | null;
  onSubmit: (input: ProductoInput) => Promise<void>;
  onDone: () => void;
}

function ProductoForm({ item, onSubmit, onDone }: ProductoFormProps) {
  const [form, setForm] = useState<ProductoInput>(
    item
      ? { nombre: item.nombre, precio: item.precio, stock: item.stock }
      : emptyForm
  );
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      await onSubmit(form);
      onDone();
    } catch {
      // el error ya se muestra en la página; el dialog queda abierto
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <DialogHeader>
        <DialogTitle>{item ? "Editar producto" : "Nuevo producto"}</DialogTitle>
      </DialogHeader>
      <div className="grid gap-4 py-4">
        <div className="grid gap-2">
          <Label htmlFor="nombre">Nombre</Label>
          <Input
            id="nombre"
            value={form.nombre}
            onChange={(e) => setForm({ ...form, nombre: e.target.value })}
            required
          />
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div className="grid gap-2">
            <Label htmlFor="precio">Precio (CLP, con impuesto)</Label>
            <Input
              id="precio"
              type="number"
              step="1"
              min="1"
              value={form.precio}
              onChange={(e) =>
                setForm({ ...form, precio: Number(e.target.value) })
              }
              required
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="stock">Stock</Label>
            <Input
              id="stock"
              type="number"
              step="1"
              min="0"
              value={form.stock}
              onChange={(e) =>
                setForm({ ...form, stock: Number(e.target.value) })
              }
              required
            />
          </div>
        </div>
      </div>
      <DialogFooter>
        <Button type="submit" disabled={submitting}>
          {item ? "Guardar cambios" : "Crear producto"}
        </Button>
      </DialogFooter>
    </form>
  );
}
