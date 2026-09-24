import Link from "next/link";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function Home() {
  return (
    <div className="grid gap-6 sm:grid-cols-2">
      <Link href="/productos">
        <Card className="transition-colors hover:border-primary">
          <CardHeader>
            <CardTitle>Productos</CardTitle>
          </CardHeader>
          <CardContent className="text-sm text-muted-foreground">
            Administrar el catálogo: precio, stock, crear, editar y eliminar.
          </CardContent>
        </Card>
      </Link>
      <Link href="/boilerplate">
        <Card className="transition-colors hover:border-primary">
          <CardHeader>
            <CardTitle>Boilerplate</CardTitle>
          </CardHeader>
          <CardContent className="text-sm text-muted-foreground">
            Crear, listar, filtrar, editar y eliminar (soft delete) registros.
          </CardContent>
        </Card>
      </Link>
    </div>
  );
}
