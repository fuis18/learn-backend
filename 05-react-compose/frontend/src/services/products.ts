import { apiFetch } from "@/lib/api";
import type { CreateProduct, Product } from "@/types/products";

const RESOURCE = "/products";

// GET /api/products
export const listProducts = () => apiFetch<Product[]>(RESOURCE);

// GET /api/products/{id}
export const getProduct = (id: number) =>
  apiFetch<Product>(`${RESOURCE}/${id}`);

// POST /api/products → 201
export const createProduct = (data: CreateProduct) =>
  apiFetch<Product>(RESOURCE, { method: "POST", body: JSON.stringify(data) });

// PATCH /api/products/{id} → updates only the fields sent
export const updateProduct = (id: number, data: Partial<CreateProduct>) =>
  apiFetch<Product>(`${RESOURCE}/${id}`, {
    method: "PATCH",
    body: JSON.stringify(data),
  });

// Duplicate checks (client-side): the API has no GET /products?code= endpoint,
// so we filter the full list. Returns undefined when no match.
export const getProductByCode = async (code: string) => {
  const products = await listProducts();
  return products.find((p) => p.code === code);
};

export const getProductByName = async (name: string) => {
  const products = await listProducts();
  return products.find((p) => p.name === name);
};

// DELETE /api/products/{id} → 204
export const deleteProduct = (id: number) =>
  apiFetch<void>(`${RESOURCE}/${id}`, { method: "DELETE" });
