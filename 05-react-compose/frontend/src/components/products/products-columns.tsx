import { createColumnHelper } from "@tanstack/react-table";
import type { DataTableFeatures } from "@/lib/table-features";
import type { Product } from "@/types/products";
import DeleteRow from "./DeleteRow";
import ProductDialog from "./ProductDialog";

const columnHelper = createColumnHelper<DataTableFeatures, Product>();

const priceFormatter = new Intl.NumberFormat("en-US", {
  style: "currency",
  currency: "USD",
});

export const productsColumns = columnHelper.columns([
  columnHelper.accessor("code", {
    header: "Code",
  }),
  columnHelper.accessor("name", {
    header: "Name",
  }),
  columnHelper.accessor("price", {
    header: "Price",
    cell: (info) => priceFormatter.format(info.getValue()),
  }),
  columnHelper.display({
    id: "actions",
    header: "Action",
    cell: (info) => (
      <div className="flex gap-2">
        <ProductDialog product={info.row.original} />
        <DeleteRow product={info.row.original} />
      </div>
    ),
  }),
]);
