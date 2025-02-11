import { Stack } from "@mui/material";
import {
  DataGrid,
  DataGridProps,
  GridCellParams,
  GridColumns,
  GridSelectionModel,
  GridValueGetterParams,
} from "@mui/x-data-grid";
import { isFunction } from "lodash";
import { memo } from "react";
import { DestinationTypeCell } from "./cells";

import styles from "./cells.module.scss";

export enum DestinationsTableField {
  NAME = "name",
  TYPE = "type",
}

interface DestinationsDataGridProps extends Omit<DataGridProps, "columns"> {
  setSelectionModel?: (names: GridSelectionModel) => void;
  onEditDestination: (name: string) => void;
  loading: boolean;
  columnFields?: DestinationsTableField[];
  minHeight?: string;
  selectionModel?: GridSelectionModel;
  destinationsPage?: boolean;
}

export const DestinationsDataGrid: React.FC<DestinationsDataGridProps> = memo(
  ({
    setSelectionModel,
    onEditDestination,
    columnFields,
    minHeight,
    selectionModel,
    destinationsPage,
    ...dataGridProps
  }) => {
    function renderNameCell(cellParams: GridCellParams<string>): JSX.Element {
      if (cellParams.row.kind === "Destination") {
        return (
          <button
            onClick={() => onEditDestination(cellParams.value!)}
            className={styles.link}
          >
            {cellParams.value}
          </button>
        );
      }

      return renderStringCell(cellParams);
    }

    const columns: GridColumns = (columnFields || []).map((field) => {
      switch (field) {
        case DestinationsTableField.NAME:
          return {
            field: DestinationsTableField.NAME,
            width: 300,

            headerName: "Name",
            valueGetter: (params: GridValueGetterParams) =>
              params.row.metadata.name,
            renderCell: renderNameCell,
          };
        case DestinationsTableField.TYPE:
          return {
            field: DestinationsTableField.TYPE,
            flex: 1,
            headerName: "Type",
            valueGetter: (params: GridValueGetterParams) =>
              params.row.spec.type,
            renderCell: renderTypeCell,
          };
        default:
          return { field: DestinationsTableField.TYPE };
      }
    });

    return (
      <DataGrid
        {...dataGridProps}
        checkboxSelection={isFunction(setSelectionModel)}
        onSelectionModelChange={setSelectionModel}
        components={{
          NoRowsOverlay: () => (
            <Stack height="100%" alignItems="center" justifyContent="center">
              No Destinations
            </Stack>
          ),
        }}
        style={{ minHeight }}
        disableSelectionOnClick
        getRowId={(row) => `${row.kind}|${row.metadata.name}`}
        columns={columns}
        selectionModel={selectionModel}
      />
    );
  }
);

function renderTypeCell(cellParams: GridCellParams<string>): JSX.Element {
  return <DestinationTypeCell type={cellParams.value ?? ""} />;
}

function renderStringCell(cellParams: GridCellParams<string>): JSX.Element {
  return <>{cellParams.value}</>;
}

DestinationsDataGrid.defaultProps = {
  minHeight: "calc(100vh - 250px)",
  columnFields: [DestinationsTableField.NAME, DestinationsTableField.TYPE],
};
