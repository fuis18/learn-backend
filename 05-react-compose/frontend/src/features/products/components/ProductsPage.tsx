import ProductsTable from "./ProductsTable";
import ProductsOptions from "./ProductsOptions";
import useProducts from "../hooks/useProducts";

const ProductsPage = () => {
  const { products, isLoading, error } = useProducts();

  return (
    <main className="container-dynamic">
      <div>
        <ProductsTable data={products} isLoading={isLoading} error={error} />
      </div>
      <ProductsOptions />
    </main>
  );
};

export default ProductsPage;
