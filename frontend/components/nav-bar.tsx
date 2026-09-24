"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useCart } from "@/lib/cart";
import { cn } from "@/lib/utils";

const links = [
  { href: "/tienda", label: "Tienda" },
  { href: "/productos", label: "Productos" },
  { href: "/boletas", label: "Boletas" },
  { href: "/boilerplate", label: "Boilerplate" },
];

export function NavBar() {
  const pathname = usePathname();
  const { count } = useCart();

  return (
    <header className="border-b">
      <nav className="mx-auto flex max-w-5xl items-center gap-6 px-6 py-4">
        <span className="font-semibold">Carvuk CRUD</span>
        <div className="flex gap-4">
          {links.map((link) => (
            <Link
              key={link.href}
              href={link.href}
              className={cn(
                "text-sm text-muted-foreground hover:text-foreground transition-colors",
                pathname?.startsWith(link.href) && "text-foreground font-medium"
              )}
            >
              {link.label}
              {link.href === "/tienda" && count > 0 && (
                <span className="ml-1 rounded-full bg-primary px-1.5 text-xs text-primary-foreground">
                  {count}
                </span>
              )}
            </Link>
          ))}
        </div>
      </nav>
    </header>
  );
}
