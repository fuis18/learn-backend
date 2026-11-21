import { z } from "zod";

/* ---------- Entity ---------- */
export interface Product {
  id: number;
  code: string;
  name: string;
  price: number;
}

export interface CreateProduct {
  code?: string;
  name: string;
  price: number;
}

/* ---------- Validation ---------- */
const optionalCode = z.preprocess(
  (value) => (value === "" || value === null ? undefined : value),
  z.string().trim().regex(/^\d+$/, "Invalid code").optional(),
);

export const productSchema = z.object({
  code: optionalCode,
  name: z.string().min(3, "Invalid name").max(30, "Máximo 30 caracteres"),
  price: z.coerce.number().positive("Invalid price"),
});

export type ProductFormInput = z.input<typeof productSchema>;
export type ProductFormValues = z.output<typeof productSchema>;
