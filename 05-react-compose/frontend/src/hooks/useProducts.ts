import { useQuery } from "@tanstack/react-query";
import { listProducts } from "@/services/products";
import type { Product } from "@/types/products";

const EMPTY_PRODUCTS: Product[] = [];

const useProducts = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["products"],
    queryFn: listProducts,
  });

  return {
    // data
    products: data ?? EMPTY_PRODUCTS,
    isLoading,
    error,
    // actions
  };
};
export default useProducts;
