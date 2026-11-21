import { useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useFormDialog } from "@/components/FormDialog";
import { ApiError } from "@/lib/api";

import {
  productSchema,
  type Product,
  type ProductFormInput,
  type ProductFormValues,
} from "@/features/products/types/products.types";

import {
  getProductByCode,
  getProductByName,
  createProduct,
  updateProduct,
} from "@/features/products/api/products.api";

interface ProductFormProps {
  product?: Product;
}

const ProductForm = ({ product }: ProductFormProps) => {
  const { close, show } = useFormDialog();
  const queryClient = useQueryClient();
  const isEdit = !!product;

  const form = useForm<ProductFormInput, unknown, ProductFormValues>({
    resolver: zodResolver(productSchema),
    defaultValues: {
      code: product?.code ?? "",
      name: product?.name ?? "",
      price: product?.price ?? 0,
    },
  });

  const isDisabled =
    isEdit && Object.keys(form.formState.dirtyFields).length === 0;

  const onSubmit = async (data: ProductFormValues) => {
    try {
      const parsed = productSchema.parse(data);

      if (parsed.code !== undefined) {
        const existingByCode = await getProductByCode(parsed.code);
        if (existingByCode && (!isEdit || existingByCode.id !== product?.id)) {
          form.setError("code", {
            type: "manual",
            message: "Existing code",
          });
          return;
        }
      }

      const existingByName = await getProductByName(parsed.name);
      if (existingByName && (!isEdit || existingByName.id !== product?.id)) {
        form.setError("name", { type: "manual", message: "Nombre existente" });
        return;
      }

      if (isEdit && product) {
        const { dirtyFields } = form.formState;

        const changes: Partial<ProductFormValues> = {};

        if (dirtyFields.code && parsed.code !== undefined)
          changes.code = parsed.code;
        if (dirtyFields.name) changes.name = parsed.name;
        if (dirtyFields.price) changes.price = parsed.price;

        if (Object.keys(changes).length === 0) {
          close();
          return;
        }

        await updateProduct(product.id, changes);

        show({
          variant: "success",
          title: "Product updated",
          description: "The changes were saved successfully.",
        });
      } else {
        await createProduct(parsed);

        show({
          variant: "success",
          title: "Product created",
          description: "The product was added successfully.",
        });
      }

      close();
      queryClient.invalidateQueries({ queryKey: ["products"] });
    } catch (err) {
      show({
        variant: "error",
        title: isEdit ? "Error updating" : "Error during creation",
        description:
          err instanceof ApiError
            ? err.message
            : "The product could not be saved",
      });
    }
  };

  return (
    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 mt-2">
      <div className="flex gap-1 flex-col">
        <Label htmlFor="code">Code</Label>
        <Input
          id="code"
          type="text"
          inputMode="numeric"
          placeholder="Automatic"
          {...form.register("code")}
        />
        {form.formState.errors.code && (
          <p className="text-destructive text-sm">
            {form.formState.errors.code.message}
          </p>
        )}
      </div>

      <div className="flex gap-1 flex-col">
        <Label htmlFor="name">Name</Label>
        <Input id="name" {...form.register("name")} />
        {form.formState.errors.name && (
          <p className="text-destructive text-sm">
            {form.formState.errors.name.message}
          </p>
        )}
      </div>

      <div className="flex gap-1 flex-col">
        <Label htmlFor="price">Precio</Label>
        <Input id="price" type="text" {...form.register("price")} />
        {form.formState.errors.price && (
          <p className="text-destructive text-sm">
            {form.formState.errors.price.message}
          </p>
        )}
      </div>

      <Button type="submit" disabled={isDisabled}>
        {isEdit ? "Save changes" : "Create"}
      </Button>
    </form>
  );
};

export default ProductForm;
