"use client";

import Link from "next/link";
import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { ApiError, boletasApi } from "@/lib/api";
import { formatCLP, formatDateTime } from "@/lib/format";
import type { BoletaListItem } from "@/lib/types";

export default function BoletasPage() {
  const [items, setItems] = useState<BoletaListItem[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);

  const limit = 10;

  const loadItems = useCallback(async () => {
    setLoading(true);
    try {
      const res = await boletasApi.list({ page, limit, order: "desc" });
      setItems(res.data);
      setTotal(res.total);
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Error al cargar boletas");
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => {
    // Fetch-on-mount/page-change: syncing with the external API.
    // eslint-disable-next-line react-hooks/set-state-in-effect
    loadItems();
  }, [loadItems]);

  const totalPages = Math.max(1, Math.ceil(total / limit));

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">Boletas</h1>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>N°</TableHead>
              <TableHead>Fecha</TableHead>
              <TableHead>Usuario</TableHead>
              <TableHead>Productos</TableHead>
              <TableHead>Total</TableHead>
              <TableHead className="text-right">Acciones</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center text-muted-foreground">
                  Cargando...
                </TableCell>
              </TableRow>
            ) : items.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center text-muted-foreground">
                  Aún no hay boletas
                </TableCell>
              </TableRow>
            ) : (
              items.map((item) => (
                <TableRow key={item.id}>
                  <TableCell>{item.id}</TableCell>
                  <TableCell>{formatDateTime(item.created_at)}</TableCell>
                  <TableCell>
                    <div>{item.usuario.nombre}</div>
                    <div className="text-xs text-muted-foreground">{item.usuario.correo}</div>
                  </TableCell>
                  <TableCell>{item.cantidad_items}</TableCell>
                  <TableCell>{formatCLP(item.valor_bruto)}</TableCell>
                  <TableCell className="text-right">
                    <Link
                      href={`/boletas/${item.id}`}
                      className={buttonVariants({ variant: "outline", size: "sm" })}
                    >
                      Ver detalle
                    </Link>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <div className="flex items-center justify-between text-sm text-muted-foreground">
        <span>
          Página {page} de {totalPages} ({total} boletas)
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
    </div>
  );
}
