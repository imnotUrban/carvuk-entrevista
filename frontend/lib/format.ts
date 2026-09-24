const clp = new Intl.NumberFormat("es-CL", {
  style: "currency",
  currency: "CLP",
  maximumFractionDigits: 0,
});

export const formatCLP = (n: number) => clp.format(n);

export const formatDateTime = (iso: string) =>
  new Date(iso).toLocaleString("es-CL", { dateStyle: "short", timeStyle: "short" });
