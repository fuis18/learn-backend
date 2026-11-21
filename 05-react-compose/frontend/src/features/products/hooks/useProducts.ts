import { useQuery } from "@tanstack/react-query";
import { listProducts } from "../api/products.api";
import type { Product } from "../types/products.types";

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
