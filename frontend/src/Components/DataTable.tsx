import {
  Button,
  Link,
  Menu,
  MenuButton,
  MenuItemOption,
  MenuList,
  MenuOptionGroup,
  Table,
  TableContainer,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import { GPUStats, WorkStationGroup } from "../Data";
import { useEffect, useMemo, useState } from "react";
import { useForceUpdate } from "framer-motion";
import { keepIf } from "../Utils/Utils";

import { Link as ReactRouterLink, useSearchParams } from "react-router-dom";

export const GPU_FIELDS = {
  "GPU Name": "gpu_name",
  Free: "free",
  "GPU Brand": "gpu_brand",
  "Driver Version": "driver_ver",
  "Memory Total (MB)": "memory_total",
  "Memory Utilisation (%)": "memory_util",
  "GPU Utilisation (%)": "gpu_util",
  "Memory Used (MB)": "memory_used",
  "Fan Speed (%)": "fan_speed",
  "GPU Temperature (°C)": "gpu_temp",
  "Memory Temperature (°C)": "memory_temp",
  "GPU Voltage (mV)": "graphics_voltage",
  "Power Draw (W)": "power_draw",
  "GPU Clock (MHz)": "graphics_clock",
  "Max GPU Clock (MHz)": "max_graphics_clock",
  "Memory Clock (MHz)": "memory_clock",
  "Max Memory Clock (MHz)": "max_memory_clock",
} as const;

const colToIdx = (key: TableViewCol) => {
  if (key === "Group") return 0;
  if (key === "Machine Name") return 1;

  return Object.keys(GPU_FIELDS).indexOf(key) + 2;
};

const Row = ({
  params,
  row,
}: {
  params: URLSearchParams;
  row: (string | number | null)[];
}) => {
  const [hover, setHover] = useState(false);
  const newParams = new URLSearchParams(
    Object.fromEntries(Array.from(params.entries())),
  );

  newParams.append("selected", row[colToIdx("Machine Name")]!.toString());

  return (
    <Tr>
      <Link
        as={ReactRouterLink}
        to={{ search: newParams.toString() }}
        display="contents"
        onMouseEnter={() => {
          setHover(true);
        }}
        onMouseLeave={() => {
          setHover(false);
        }}
      >
        {row.map((s) =>
          s === null ? null : (
            <Td key={s}>
              <Text textDecoration={hover ? "underline" : ""}>{s}</Text>
            </Td>
          ),
        )}
      </Link>
    </Tr>
  );
};

type TableViewCol = keyof typeof GPU_FIELDS | "Group" | "Machine Name";
type Direction = "ascending" | "descending";

export const TableTab = ({ groups }: { groups: WorkStationGroup[] }) => {
  // default to show group, machine_name, gpu_name, isFree, brand, and memory_total
  const SHOWN_COLS: {
    [_ in TableViewCol]: boolean;
  } = {
    Group: true,
    "Machine Name": true,
    "GPU Name": true,
    Free: true,
    "GPU Brand": true,
    "Driver Version": false,
    "Memory Total (MB)": true,
    "Memory Utilisation (%)": false,
    "GPU Utilisation (%)": false,
    "Memory Used (MB)": false,
    "Fan Speed (%)": false,
    "GPU Temperature (°C)": false,
    "Memory Temperature (°C)": false,
    "GPU Voltage (mV)": false,
    "Power Draw (W)": false,
    "GPU Clock (MHz)": false,
    "Max GPU Clock (MHz)": false,
    "Memory Clock (MHz)": false,
    "Max Memory Clock (MHz)": false,
  };

  const [shownColumns, setter] = useState<Record<string, boolean>>(SHOWN_COLS);
  const [refresh, subscribe] = useForceUpdate();
  const [params] = useSearchParams();
  const [sortConfig, setSortConfig] = useState<{
    key: TableViewCol;
    direction: Direction;
  }>({
    key: "GPU Name",
    direction: "ascending",
  });

  const [rows, setRows] = useState<(string | number | null)[][]>([]);

  const sortedGroups = useMemo(() => {
    const sortableItems = [...rows];
    if (sortConfig !== null) {
      console.log(sortConfig);
      sortableItems.sort((a, b) => {
        if (a[colToIdx(sortConfig.key)]! < b[colToIdx(sortConfig.key)]!) {
          return sortConfig.direction === "ascending" ? -1 : 1;
        }
        if (a[colToIdx(sortConfig.key)]! > b[colToIdx(sortConfig.key)]!) {
          return sortConfig.direction === "ascending" ? 1 : -1;
        }
        return 0;
      });
    }

    return sortableItems;
  }, [rows, sortConfig]);

  const requestSort = (key: TableViewCol) => {
    const direction =
      sortConfig.key === key && sortConfig.direction === "ascending"
        ? "ascending"
        : "descending";
    setSortConfig({ key, direction });
  };

  useEffect(() => {
    setRows(
      groups.flatMap(({ name: group_name, workstations }) =>
        workstations.flatMap(({ name: workstation_name, gpus }) =>
          gpus.map((gpu) => {
            const rows: (string | number | null)[] = tablify(
              shownColumns,
              gpu,
              group_name,
              workstation_name,
            );
            return rows;
          }),
        ),
      ),
    );
  }, [groups, shownColumns, subscribe]);

  const shown = Object.keys(shownColumns) as TableViewCol[];

  return (
    <div>
      <Menu closeOnSelect={false}>
        <MenuButton as={Button} colorScheme="blue">
          Columns
        </MenuButton>
        <MenuList overflowY="scroll" maxHeight="200">
          <MenuOptionGroup
            type="checkbox"
            defaultValue={shown.filter((key) => shownColumns[key])}
            onChange={(props) => {
              shown.forEach((col) => {
                shownColumns[col] = props.includes(col);
              });
              setter(shownColumns);
              refresh();
            }}
          >
            {shown.map((col, i) => (
              <MenuItemOption value={col} key={i}>
                {col}
              </MenuItemOption>
            ))}
          </MenuOptionGroup>
        </MenuList>
      </Menu>

      <TableContainer>
        <Table variant="striped">
          <Thead>
            <Tr>
              {shown.map((col, i) =>
                shownColumns[col] ? (
                  <Th key={i} cursor="pointer" onClick={() => requestSort(col)}>
                    {col}
                  </Th>
                ) : null,
              )}
            </Tr>
          </Thead>
          <Tbody>
            {sortedGroups.map((row, i) => (
              <Row key={i} params={params} row={row} />
            ))}
          </Tbody>
        </Table>
      </TableContainer>
    </div>
  );
};

export const tablify = (
  shownColumns: Record<string, boolean>,
  gpu: GPUStats,
  group_name?: string,
  workstation_name?: string,
): (string | null)[] =>
  (group_name === undefined ? [] : [keepIf(shownColumns["Group"], group_name)])
    .concat(
      workstation_name === undefined
        ? []
        : [keepIf(shownColumns["Machine Name"], workstation_name)],
    )
    .concat([
      keepIf(shownColumns["GPU Name"], gpu.gpu_name),
      keepIf(shownColumns["Free"], gpu.in_use ? "❌" : "✅"),
      keepIf(shownColumns["GPU Brand"], gpu.gpu_brand),
      keepIf(shownColumns["Driver Version"], gpu.driver_ver),
      keepIf(
        shownColumns["Memory Total (MB)"],
        Math.round(gpu.memory_total).toString(),
      ),
      keepIf(
        shownColumns["Memory Utilisation (%)"],
        Math.round(gpu.memory_util).toString(),
      ),
      keepIf(
        shownColumns["GPU Utilisation (%)"],
        Math.round(gpu.gpu_util).toString(),
      ),
      keepIf(
        shownColumns["Memory Used (MB)"],
        Math.round(gpu.memory_used).toString(),
      ),
      keepIf(
        shownColumns["Fan Speed (%)"],
        Math.round(gpu.fan_speed).toString(),
      ),
      keepIf(
        shownColumns["GPU Temperature (°C)"],
        Math.round(gpu.gpu_temp).toString(),
      ),
      keepIf(
        shownColumns["Memory Temperature (°C)"],
        Math.round(gpu.memory_temp).toString(),
      ),
      keepIf(
        shownColumns["GPU Voltage (mV)"],
        Math.round(gpu.graphics_voltage).toString(),
      ),
      keepIf(
        shownColumns["Power Draw (W)"],
        Math.round(gpu.power_draw).toString(),
      ),
      keepIf(
        shownColumns["GPU Clock (MHz)"],
        Math.round(gpu.graphics_clock).toString(),
      ),
      keepIf(
        shownColumns["Max GPU Clock (MHz)"],
        Math.round(gpu.max_graphics_clock).toString(),
      ),
      keepIf(
        shownColumns["Memory Clock (MHz)"],
        Math.round(gpu.memory_clock).toString(),
      ),
      keepIf(
        shownColumns["Max Memory Clock (MHz)"],
        Math.round(gpu.max_memory_clock).toString(),
      ),
    ]);
