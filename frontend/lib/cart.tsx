"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import type { Producto } from "./types";

export const CART_STORAGE_KEY = "carvuk_cart";
const MAX_PER_LINE = 99;

export interface CartItem {
  producto_id: number;
  nombre: string;
  precio: number;
  stock: number;
  cantidad: number;
}

interface CartContextValue {
  items: CartItem[];
  total: number;
  count: number;
  add: (producto: Producto) => void;
  setQuantity: (id: number, n: number) => void;
  remove: (id: number) => void;
  clear: () => void;
}

const CartContext = createContext<CartContextValue | null>(null);

export const maxQuantity = (stock: number) => Math.min(MAX_PER_LINE, stock);

function isCartItem(v: unknown): v is CartItem {
  if (typeof v !== "object" || v === null) return false;
  const i = v as Record<string, unknown>;
  return (
    Number.isInteger(i.producto_id) &&
    typeof i.nombre === "string" &&
    typeof i.precio === "number" &&
    Number.isInteger(i.stock) &&
    Number.isInteger(i.cantidad) &&
    (i.cantidad as number) >= 1
  );
}

function loadCart(): CartItem[] {
  try {
    const raw = window.localStorage.getItem(CART_STORAGE_KEY);
    if (!raw) return [];
    const parsed: unknown = JSON.parse(raw);
    if (!Array.isArray(parsed)) return [];
    return parsed
      .filter(isCartItem)
      .map((i) => ({ ...i, cantidad: Math.min(i.cantidad, maxQuantity(i.stock)) }))
      .filter((i) => i.cantidad >= 1);
  } catch {
    return [];
  }
}

export function CartProvider({ children }: { children: ReactNode }) {
  const [items, setItems] = useState<CartItem[]>([]);
  const [hydrated, setHydrated] = useState(false);

  useEffect(() => {
    // Hydrate from localStorage after mount (avoids SSR mismatch).
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setItems(loadCart());
    setHydrated(true);
  }, []);

  useEffect(() => {
    if (!hydrated) return;
    try {
      window.localStorage.setItem(CART_STORAGE_KEY, JSON.stringify(items));
    } catch {
      // storage unavailable: cart just won't persist
    }
  }, [items, hydrated]);

  const add = useCallback((p: Producto) => {
    setItems((prev) => {
      const max = maxQuantity(p.stock);
      if (max < 1) return prev;
      const existing = prev.find((i) => i.producto_id === p.id);
      if (!existing) {
        return [
          ...prev,
          { producto_id: p.id, nombre: p.nombre, precio: p.precio, stock: p.stock, cantidad: 1 },
        ];
      }
      return prev.map((i) =>
        i.producto_id === p.id
          ? { ...i, stock: p.stock, cantidad: Math.min(i.cantidad + 1, max) }
          : i
      );
    });
  }, []);

  const setQuantity = useCallback((id: number, n: number) => {
    setItems((prev) =>
      prev
        .map((i) =>
          i.producto_id === id
            ? { ...i, cantidad: Math.min(Math.floor(n), maxQuantity(i.stock)) }
            : i
        )
        .filter((i) => i.cantidad >= 1)
    );
  }, []);

  const remove = useCallback(
    (id: number) => setItems((prev) => prev.filter((i) => i.producto_id !== id)),
    []
  );
  const clear = useCallback(() => setItems([]), []);

  const value = useMemo<CartContextValue>(
    () => ({
      items,
      total: items.reduce((s, i) => s + i.precio * i.cantidad, 0),
      count: items.reduce((s, i) => s + i.cantidad, 0),
      add,
      setQuantity,
      remove,
      clear,
    }),
    [items, add, setQuantity, remove, clear]
  );

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
}

export function useCart(): CartContextValue {
  const ctx = useContext(CartContext);
  if (!ctx) throw new Error("useCart must be used within CartProvider");
  return ctx;
}
