import { ButtonGroup } from "@/components/ui/button-group";
import ProductDialog from "@/components/products/ProductDialog";
import ProductsTable from "@/components/products/ProductsTable";
import useProducts from "@/hooks/useProducts";

const Products = () => {
  const { products, isLoading, error } = useProducts();

  return (
    <main className="container-dynamic">
      <div>
        <ProductsTable data={products} isLoading={isLoading} error={error} />
      </div>
      <ButtonGroup orientation="vertical" className="h-fit">
        <ProductDialog />
      </ButtonGroup>
    </main>
  );
};

export default Products;
