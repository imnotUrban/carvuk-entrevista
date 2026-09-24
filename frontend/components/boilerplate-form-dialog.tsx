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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type {
  Boilerplate,
  BoilerplateInput,
  BoilerplateStatus,
} from "@/lib/types";

interface BoilerplateFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  item?: Boilerplate | null;
  onSubmit: (input: BoilerplateInput) => Promise<void>;
}

export const STATUS_LABELS: Record<BoilerplateStatus, string> = {
  draft: "Borrador",
  active: "Activo",
  archived: "Archivado",
};

const emptyForm: BoilerplateInput = {
  name: "",
  code: "",
  status: "draft",
  quantity: 0,
  amount: 0,
};

export function BoilerplateFormDialog({
  open,
  onOpenChange,
  item,
  onSubmit,
}: BoilerplateFormDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        {open && (
          <BoilerplateForm
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

interface BoilerplateFormProps {
  item?: Boilerplate | null;
  onSubmit: (input: BoilerplateInput) => Promise<void>;
  onDone: () => void;
}

function BoilerplateForm({ item, onSubmit, onDone }: BoilerplateFormProps) {
  const [form, setForm] = useState<BoilerplateInput>(
    item
      ? {
          name: item.name,
          code: item.code,
          status: item.status,
          quantity: item.quantity,
          amount: item.amount,
        }
      : emptyForm
  );
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      await onSubmit(form);
      onDone();
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <DialogHeader>
        <DialogTitle>{item ? "Editar registro" : "Nuevo registro"}</DialogTitle>
      </DialogHeader>
      <div className="grid gap-4 py-4">
        <div className="grid gap-2">
          <Label htmlFor="name">Nombre</Label>
          <Input
            id="name"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            required
          />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="code">Código</Label>
          <Input
            id="code"
            value={form.code}
            onChange={(e) => setForm({ ...form, code: e.target.value })}
            required
          />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="status">Estado</Label>
          <Select
            value={form.status}
            onValueChange={(value) =>
              setForm({ ...form, status: (value ?? "draft") as BoilerplateStatus })
            }
          >
            <SelectTrigger id="status">
              <SelectValue placeholder="Selecciona un estado" />
            </SelectTrigger>
            <SelectContent>
              {(Object.keys(STATUS_LABELS) as BoilerplateStatus[]).map((status) => (
                <SelectItem key={status} value={status}>
                  {STATUS_LABELS[status]}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div className="grid gap-2">
            <Label htmlFor="quantity">Cantidad</Label>
            <Input
              id="quantity"
              type="number"
              min="0"
              value={form.quantity}
              onChange={(e) =>
                setForm({ ...form, quantity: Number(e.target.value) })
              }
              required
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="amount">Monto</Label>
            <Input
              id="amount"
              type="number"
              step="0.01"
              min="0"
              value={form.amount}
              onChange={(e) =>
                setForm({ ...form, amount: Number(e.target.value) })
              }
              required
            />
          </div>
        </div>
      </div>
      <DialogFooter>
        <Button type="submit" disabled={submitting}>
          {item ? "Guardar cambios" : "Crear registro"}
        </Button>
      </DialogFooter>
    </form>
  );
}
