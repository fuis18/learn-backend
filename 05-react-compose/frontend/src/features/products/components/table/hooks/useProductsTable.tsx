import { useTable } from "@tanstack/react-table";
import { CONFIG } from "@/config";
import type { Product } from "../../../types/products.types";
import { features } from "../products-features";
import { productsColumns } from "../products-columns";

const useProductsTable = (data: Product[]) => {
  return useTable({
    features,
    columns: productsColumns,
    data,
    getRowId: (row) => String(row.id),
    initialState: {
      pagination: {
        pageIndex: 0,
        pageSize: CONFIG.LIMIT,
      },
    },
  });
};

export default useProductsTable;
