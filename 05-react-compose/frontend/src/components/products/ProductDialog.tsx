import { Button } from "@/components/ui/button";
import FormDialog from "@/components/FormDialog";
import ProductForm from "./ProductForm";
import type { Product } from "@/types/products";

interface ProductDialogProps {
  product?: Product;
}

const ProductDialog = ({ product }: ProductDialogProps) => {
  const isEdit = !!product;

  return (
    <FormDialog
      title={isEdit ? "Edit product" : "Create product"}
      trigger={
        isEdit ? (
          <button type="button" className="text-success hover:text-success/80">
            Edit
          </button>
        ) : (
          <Button type="button">Create</Button>
        )
      }
    >
      <ProductForm product={product} />
    </FormDialog>
  );
};

export default ProductDialog;
