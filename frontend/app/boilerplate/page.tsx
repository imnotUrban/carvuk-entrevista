"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  BoilerplateFormDialog,
  STATUS_LABELS,
} from "@/components/boilerplate-form-dialog";
import { BoilerplateViewDialog } from "@/components/boilerplate-view-dialog";
import { ApiError, boilerplateApi } from "@/lib/api";
import type {
  Boilerplate,
  BoilerplateInput,
  BoilerplateStatus,
} from "@/lib/types";

const STATUS_ANY = "any";

const STATUS_VARIANT: Record<BoilerplateStatus, "default" | "secondary" | "outline"> = {
  draft: "secondary",
  active: "default",
  archived: "outline",
};

export default function BoilerplatePage() {
  const [items, setItems] = useState<Boilerplate[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<Boilerplate | null>(null);
  const [viewingItem, setViewingItem] = useState<Boilerplate | null>(null);

  const [nameFilter, setNameFilter] = useState("");
  const [codeFilter, setCodeFilter] = useState("");
  const [statusFilter, setStatusFilter] = useState(STATUS_ANY);

  const limit = 10;

  const loadItems = useCallback(async () => {
    setLoading(true);
    try {
      const res = await boilerplateApi.list({
        name: nameFilter,
        code: codeFilter,
        status:
          statusFilter === STATUS_ANY ? undefined : (statusFilter as BoilerplateStatus),
        page,
        limit,
      });
      setItems(res.data);
      setTotal(res.total);
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Error al cargar registros");
    } finally {
      setLoading(false);
    }
  }, [nameFilter, codeFilter, statusFilter, page]);

  useEffect(() => {
    // Fetch-on-mount/filter-change: syncing with the external API.
    // eslint-disable-next-line react-hooks/set-state-in-effect
    loadItems();
  }, [loadItems]);

  const handleCreateOrUpdate = async (input: BoilerplateInput) => {
    try {
      if (editingItem) {
        await boilerplateApi.update(editingItem.id, input);
        toast.success("Registro actualizado");
      } else {
        await boilerplateApi.create(input);
        toast.success("Registro creado");
      }
      await loadItems();
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Error al guardar registro");
      throw err;
    }
  };

  const handleDelete = async (item: Boilerplate) => {
    if (!confirm(`¿Eliminar el registro "${item.name}"? (soft delete)`)) return;
    try {
      await boilerplateApi.remove(item.id);
      toast.success("Registro eliminado");
      await loadItems();
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Error al eliminar registro");
    }
  };

  const totalPages = Math.max(1, Math.ceil(total / limit));

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Boilerplate</h1>
        <Button
          onClick={() => {
            setEditingItem(null);
            setDialogOpen(true);
          }}
        >
          Nuevo registro
        </Button>
      </div>

      <div className="flex flex-wrap gap-2">
        <Input
          placeholder="Filtrar por nombre..."
          value={nameFilter}
          onChange={(e) => {
            setPage(1);
            setNameFilter(e.target.value);
          }}
          className="max-w-[200px]"
        />
        <Input
          placeholder="Filtrar por código..."
          value={codeFilter}
          onChange={(e) => {
            setPage(1);
            setCodeFilter(e.target.value);
          }}
          className="max-w-[160px]"
        />
        <Select
          value={statusFilter}
          onValueChange={(v) => {
            setPage(1);
            setStatusFilter(v ?? STATUS_ANY);
          }}
        >
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Estado" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value={STATUS_ANY}>Todos los estados</SelectItem>
            {(Object.keys(STATUS_LABELS) as BoilerplateStatus[]).map((status) => (
              <SelectItem key={status} value={status}>
                {STATUS_LABELS[status]}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Nombre</TableHead>
              <TableHead>Código</TableHead>
              <TableHead>Estado</TableHead>
              <TableHead>Cantidad</TableHead>
              <TableHead>Monto</TableHead>
              <TableHead className="text-right">Acciones</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={7} className="text-center text-muted-foreground">
                  Cargando...
                </TableCell>
              </TableRow>
            ) : items.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} className="text-center text-muted-foreground">
                  Sin registros
                </TableCell>
              </TableRow>
            ) : (
              items.map((item) => (
                <TableRow key={item.id}>
                  <TableCell>{item.id}</TableCell>
                  <TableCell>{item.name}</TableCell>
                  <TableCell>{item.code}</TableCell>
                  <TableCell>
                    <Badge variant={STATUS_VARIANT[item.status]}>
                      {STATUS_LABELS[item.status]}
                    </Badge>
                  </TableCell>
                  <TableCell>{item.quantity}</TableCell>
                  <TableCell>${item.amount.toFixed(2)}</TableCell>
                  <TableCell className="text-right space-x-2">
                    <Button
                      variant="secondary"
                      size="sm"
                      onClick={() => setViewingItem(item)}
                    >
                      Ver
                    </Button>
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
          Página {page} de {totalPages} ({total} registros)
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

      <BoilerplateFormDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        item={editingItem}
        onSubmit={handleCreateOrUpdate}
      />

      <BoilerplateViewDialog
        open={viewingItem !== null}
        onOpenChange={(open) => {
          if (!open) setViewingItem(null);
        }}
        item={viewingItem}
      />
    </div>
  );
}
