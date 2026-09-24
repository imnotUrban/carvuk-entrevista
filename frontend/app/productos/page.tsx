"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { ProductoFormDialog } from "@/components/producto-form-dialog";
import { ApiError, productosApi } from "@/lib/api";
import { formatCLP } from "@/lib/format";
import type { Producto,ProductoInput } from "@/lib/types";

export default function ProductosPage() {
  const [items, setItems] = useState<Producto[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<Producto | null>(null);
  const [nombreFilter, setNombreFilter] = useState("");

  const limit = 10;

  const loadItems = useCallback(async () => {
    setLoading(true);
    try {
      const res = await productosApi.list({ nombre: nombreFilter, page, limit });
      setItems(res.data);
      setTotal(res.total);
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Error al cargar productos");
    } finally {
      setLoading(false);
    }
  }, [nombreFilter, page]);

  useEffect(() => {
    // Fetch-on-mount/filter-change: syncing with the external API.
    // eslint-disable-next-line react-hooks/set-state-in-effect
    loadItems();
  }, [loadItems]);

  const handleCreateOrUpdate = async (input: ProductoInput) => {
    try {
      if (editingItem) {
        await productosApi.update(editingItem.id, input);
        toast.success("Producto actualizado");
      } else {
        await productosApi.create(input);
        toast.success("Producto creado");
      }
      await loadItems();
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Error al guardar producto");
      throw err;
    }
  };

  const handleDelete = async (item: Producto) => {
    if (!confirm(`¿Eliminar el producto "${item.nombre}"?`)) return;
    try {
      await productosApi.remove(item.id);
      toast.success("Producto eliminado");
      await loadItems();
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Error al eliminar producto");
    }
  };

  const totalPages = Math.max(1, Math.ceil(total / limit));

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Productos</h1>
        <Button
          onClick={() => {
            setEditingItem(null);
            setDialogOpen(true);
          }}
        >
          Nuevo producto
        </Button>
      </div>

      <div className="flex flex-wrap gap-2">
        <Input
          placeholder="Filtrar por nombre..."
          value={nombreFilter}
          onChange={(e) => {
            setPage(1);
            setNombreFilter(e.target.value);
          }}
          className="max-w-[240px]"
        />
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Nombre</TableHead>
              <TableHead>Precio</TableHead>
              <TableHead>Stock</TableHead>
              <TableHead className="text-right">Acciones</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={5} className="text-center text-muted-foreground">
                  Cargando...
                </TableCell>
              </TableRow>
            ) : items.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="text-center text-muted-foreground">
                  Sin productos
                </TableCell>
              </TableRow>
            ) : (
              items.map((item) => (
                <TableRow key={item.id}>
                  <TableCell>{item.id}</TableCell>
                  <TableCell>{item.nombre}</TableCell>
                  <TableCell>{formatCLP(item.precio)}</TableCell>
                  <TableCell>{item.stock}</TableCell>
                  <TableCell className="text-right space-x-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => {
                        setEditingItem(item);
                        setDialogOpen(true);
                      }}
                    >
                      Editar
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => handleDelete(item)}
                    >
                      Eliminar
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <div className="flex items-center justify-between text-sm text-muted-foreground">
        <span>
          Página {page} de {totalPages} ({total} productos)
        </span>
        <div className="flex gap-2">
          <Button
            variant="outline"
            size="sm"
            disabled={page <= 1}
            onClick={() => setPage((p) => Math.max(1, p - 1))}
          >
            Anterior
          </Button>
          <Button
            variant="outline"
            size="sm"
            disabled={page >= totalPages}
            onClick={() => setPage((p) => p + 1)}
          >
            Siguiente
          </Button>
        </div>
      </div>

      <ProductoFormDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        item={editingItem}
        onSubmit={handleCreateOrUpdate}
      />
    </div>
  );
}
