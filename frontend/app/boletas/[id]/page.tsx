"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";
import { buttonVariants } from "@/components/ui/button";
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
import type { Boleta } from "@/lib/types";

export default function BoletaDetailPage() {
  const params = useParams<{ id: string }>();
  const [boleta, setBoleta] = useState<Boleta | null>(null);
  const [status, setStatus] = useState<"loading" | "ok" | "notfound" | "error">("loading");

  useEffect(() => {
    let cancelled = false;
    boletasApi
      .get(Number(params.id))
      .then((b) => {
        if (cancelled) return;
        setBoleta(b);
        setStatus("ok");
      })
      .catch((err) => {
        if (cancelled) return;
        setStatus(err instanceof ApiError && err.status === 404 ? "notfound" : "error");
      });
    return () => {
      cancelled = true;
    };
  }, [params.id]);

  const back = (
    <Link href="/boletas" className={buttonVariants({ variant: "outline", size: "sm" })}>
      Volver a boletas
    </Link>
  );

  if (status === "loading") {
    return <p className="text-muted-foreground">Cargando...</p>;
  }
  if (status !== "ok" || !boleta) {
    return (
      <div className="space-y-4">
        <h1 className="text-2xl font-semibold">
          {status === "notfound" ? "Boleta no encontrada" : "Error al cargar la boleta"}
        </h1>
        {back}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Boleta N° {boleta.id}</h1>
        {back}
      </div>

      <div className="text-sm">
        <div>{formatDateTime(boleta.created_at)}</div>
        <div>
          {boleta.usuario.nombre} —{" "}
          <span className="text-muted-foreground">{boleta.usuario.correo}</span>
        </div>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Producto</TableHead>
              <TableHead>Precio unitario</TableHead>
              <TableHead>Cantidad</TableHead>
              <TableHead className="text-right">Subtotal</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {boleta.items.map((it) => (
              <TableRow key={it.producto_id}>
                <TableCell>{it.nombre}</TableCell>
                <TableCell>{formatCLP(it.precio_unitario)}</TableCell>
                <TableCell>{it.cantidad}</TableCell>
                <TableCell className="text-right">{formatCLP(it.subtotal)}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      <dl className="ml-auto max-w-xs space-y-1 text-sm">
        <div className="flex justify-between">
          <dt>Neto</dt>
          <dd>{formatCLP(boleta.valor_neto)}</dd>
        </div>
        <div className="flex justify-between">
          <dt>Impuesto ({boleta.porcentaje_impuesto}%)</dt>
          <dd>{formatCLP(boleta.impuesto)}</dd>
        </div>
        <div className="flex justify-between text-base font-semibold">
          <dt>Total</dt>
          <dd>{formatCLP(boleta.valor_bruto)}</dd>
        </div>
      </dl>
    </div>
  );
}
