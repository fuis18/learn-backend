import { ButtonGroup } from "@/components/ui/button-group";
import ProductDialog from "./ProductDialog";

const ProductsOptions = () => {
  return (
    <ButtonGroup orientation="vertical" className="h-fit">
      <ProductDialog />
    </ButtonGroup>
  );
};

export default ProductsOptions;
