import { useTable } from "@tanstack/react-table";
import { CONFIG } from "@/config";
import { dataTableFeatures } from "@/lib/table-features";
import type { Product } from "@/types/products";
import { productsColumns } from "./products-columns";

const useProductsTable = (data: Product[]) => {
  return useTable({
    features: dataTableFeatures,
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
