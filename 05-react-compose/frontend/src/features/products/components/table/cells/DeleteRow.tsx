import { useMutation, useQueryClient } from "@tanstack/react-query";
import AlertDialog from "@/components/AlertDialog";
import { useAlert } from "@/hooks/useAlert";
import { ApiError } from "@/lib/api";
import { deleteProduct } from "../../../api/products.api";
import type { Product } from "../../../types/products.types";

interface DeleteRowProps {
  product: Product;
}

const DeleteRow = ({ product }: DeleteRowProps) => {
  const queryClient = useQueryClient();
  const { alert, show, hide } = useAlert();

  const mutation = useMutation({
    mutationFn: () => deleteProduct(product.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["products"] });
      show({
        variant: "success",
        title: "Product removed",
        description: `"${product.name}" was successfully removed.`,
      });
    },
    onError: (error) => {
      show({
        variant: "error",
        title: "Error while deleting",
        description:
          error instanceof ApiError
            ? error.message
            : "The product could not be deleted",
      });
    },
  });

  return (
    <>
      <button
        type="button"
        className="text-destructive hover:text-destructive/80"
        disabled={mutation.isPending}
        onClick={() => mutation.mutate()}
      >
        Delete
      </button>

      <AlertDialog
        open={alert.open}
        onOpenChange={(next) => {
          if (!next) hide();
        }}
        variant={alert.variant}
        title={alert.title}
        description={alert.description}
      />
    </>
  );
};

export default DeleteRow;
