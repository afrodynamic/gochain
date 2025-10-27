'use client';
import {
  Box,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from '@mui/material';
import type { ReactNode } from 'react';

export type MRT_ColumnDef<TData extends Record<string, unknown>> = {
  accessorKey?: Extract<keyof TData, string>;
  accessorFn?: (row: TData) => ReactNode;
  header: string;
  Cell?: (ctx: { row: TData; value: unknown }) => ReactNode;
  size?: number;
  align?: 'inherit' | 'left' | 'center' | 'right' | 'justify';
};

type MaterialReactTableProps<TData extends Record<string, unknown>> = {
  columns: MRT_ColumnDef<TData>[];
  data: TData[];
  emptyContent?: ReactNode;
};

const toNode = (value: unknown): ReactNode => {
  if (value == null) return null;
  if (
    typeof value === 'string' ||
    typeof value === 'number' ||
    typeof value === 'boolean'
  ) {
    return String(value);
  }
  return value as ReactNode;
};

export function MaterialReactTable<TData extends Record<string, unknown>>({
  columns,
  data,
  emptyContent,
}: MaterialReactTableProps<TData>) {
  return (
    <TableContainer
      component={Paper}
      elevation={6}
      className="border"
      sx={{ borderColor: 'primary.main' }}
    >
      <Table size="small">
        <TableHead>
          <TableRow>
            {columns.map((column, i) => (
              <TableCell
                key={column.accessorKey ?? `col-${i}`}
                sx={{ fontWeight: 700, fontSize: 14, width: column.size }}
                align={column.align}
              >
                {column.header}
              </TableCell>
            ))}
          </TableRow>
        </TableHead>

        <TableBody>
          {data.length === 0 ? (
            <TableRow>
              <TableCell colSpan={columns.length}>
                <Box py={4} textAlign="center">
                  {emptyContent ?? (
                    <Typography variant="body2" color="text.secondary">
                      No results yet.
                    </Typography>
                  )}
                </Box>
              </TableCell>
            </TableRow>
          ) : (
            data.map((row, idx) => (
              <TableRow key={idx} hover>
                {columns.map((column, ci) => {
                  const raw =
                    column.accessorFn?.(row) ??
                    (column.accessorKey
                      ? (row[column.accessorKey] as unknown)
                      : undefined);
                  const rendered = column.Cell
                    ? column.Cell({ row, value: raw })
                    : toNode(raw);

                  return (
                    <TableCell
                      key={column.accessorKey ?? `cell-${idx}-${ci}`}
                      align={column.align}
                      sx={{ fontSize: 13 }}
                    >
                      {rendered}
                    </TableCell>
                  );
                })}
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </TableContainer>
  );
}
