import { useEffect } from "react";
import { ArrowDown, ArrowUp, ChevronsUpDown } from "lucide-react";
import type { TableFeatures, TableMeta } from "@tanstack/react-table";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import PagTable from "@/components/PaginationTable";
import useProductsTable from "./useProductsTable";
import type { Product } from "@/types/products";

interface ProductsTableProps {
  data: Product[];
  isLoading?: boolean;
  error?: Error | null;
  meta?: TableMeta<TableFeatures, Product>;
  onSelectionChange?: (rows: Product[]) => void;
}

const ProductsTable = ({ data, isLoading, error }: ProductsTableProps) => {
  const table = useProductsTable(data);
  const rows = table.getRowModel().rows;

  const pageCount = table.getPageCount();
  const pageIndex = table.state.pagination.pageIndex;

  useEffect(() => {
    if (pageCount === 0) {
      if (pageIndex !== 0) table.setPageIndex(0);
      return;
    }
    if (pageIndex > pageCount - 1) {
      table.setPageIndex(pageCount - 1);
    }
  }, [pageCount, pageIndex, table]);

  return (
    <div>
      <Table className="products-table">
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                const isSorted = header.column.getIsSorted();

                return (
                  <TableHead
                    key={header.id}
                    onClick={header.column.getToggleSortingHandler()}
                    className={
                      header.column.getCanSort()
                        ? "cursor-pointer select-none"
                        : undefined
                    }
                  >
                    <span className="inline-flex items-center gap-1">
                      {header.isPlaceholder ? null : (
                        <table.FlexRender header={header} />
                      )}
                      {header.column.getCanSort() &&
                        (isSorted === "asc" ? (
                          <ArrowUp className="h-3 w-3" />
                        ) : isSorted === "desc" ? (
                          <ArrowDown className="h-3 w-3" />
                        ) : (
                          <ChevronsUpDown className="h-3 w-3 opacity-50" />
                        ))}
                    </span>
                  </TableHead>
                );
              })}
            </TableRow>
          ))}
        </TableHeader>

        <TableBody>
          {isLoading ? (
            <TableRow>
              <TableCell
                colSpan={table.getAllLeafColumns().length}
                className="h-24 text-center text-muted-foreground"
              >
                Loading products…
              </TableCell>
            </TableRow>
          ) : error ? (
            <TableRow>
              <TableCell
                colSpan={table.getAllLeafColumns().length}
                className="h-24 text-center text-destructive"
              >
                {error.message || "Error loading products"}
              </TableCell>
            </TableRow>
          ) : rows.length ? (
            rows.map((row) => (
              <TableRow key={row.id}>
                {row.getAllCells().map((cell) => (
                  <TableCell key={cell.id}>
                    <table.FlexRender cell={cell} />
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell
                colSpan={table.getAllLeafColumns().length}
                className="h-24 text-center text-muted-foreground"
              >
                No products
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>

      {!isLoading && !error && (
        <PagTable
          page={table.state.pagination.pageIndex + 1}
          setPage={(page) => table.setPageIndex(page - 1)}
          totalPages={table.getPageCount()}
        />
      )}
    </div>
  );
};

export default ProductsTable;
