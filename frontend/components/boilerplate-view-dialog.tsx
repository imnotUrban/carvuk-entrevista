"use client";

import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { STATUS_LABELS } from "@/components/boilerplate-form-dialog";
import type { Boilerplate } from "@/lib/types";

interface BoilerplateViewDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  item: Boilerplate | null;
}

function formatDate(value: string) {
  return new Date(value).toLocaleString("es");
}

export function BoilerplateViewDialog({
  open,
  onOpenChange,
  item,
}: BoilerplateViewDialogProps) {
  const rows: [string, React.ReactNode][] = item
    ? [
        ["ID", item.id],
        ["Nombre", item.name],
        ["Código", item.code],
        ["Estado", <Badge key="status">{STATUS_LABELS[item.status]}</Badge>],
        ["Cantidad", item.quantity],
        ["Monto", `$${item.amount.toFixed(2)}`],
        ["Creado", formatDate(item.created_at)],
        ["Actualizado", formatDate(item.updated_at)],
      ]
    : [];

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Detalle del registro</DialogTitle>
        </DialogHeader>
        <dl className="grid grid-cols-[120px_1fr] gap-x-4 gap-y-3 py-2 text-sm">
          {rows.map(([label, value]) => (
            <div key={label} className="contents">
              <dt className="text-muted-foreground">{label}</dt>
              <dd>{value}</dd>
            </div>
          ))}
        </dl>
      </DialogContent>
    </Dialog>
  );
}
