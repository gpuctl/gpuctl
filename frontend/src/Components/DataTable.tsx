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
  machineName,
}: {
  params: URLSearchParams;
  row: (string | number | null)[];
  machineName: string;
}) => {
  const [hover, setHover] = useState(false);
  const newParams = new URLSearchParams(
    Object.fromEntries(Array.from(params.entries())),
  );

  newParams.append("selected", machineName);

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
        {row.map((s, i) =>
          s === null ? null : (
            <Td key={i}>
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

const invertDir = (dir: Direction) =>
  dir === "ascending" ? "descending" : "ascending";

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
  } | null>(null);

  const [rows, setRows] = useState<
    { name: string; rs: (string | number | null)[] }[]
  >([]);

  const sortedGroups = useMemo(() => {
    const sortableItems = [...rows];
    if (sortConfig !== null) {
      sortableItems.sort(({ rs: rs1 }, { rs: rs2 }) => {
        if (rs1[colToIdx(sortConfig.key)]! < rs2[colToIdx(sortConfig.key)]!) {
          return sortConfig.direction === "ascending" ? -1 : 1;
        }
        if (rs1[colToIdx(sortConfig.key)]! > rs2[colToIdx(sortConfig.key)]!) {
          return sortConfig.direction === "ascending" ? 1 : -1;
        }
        return 0;
      });
    }

    return sortableItems;
  }, [rows, sortConfig]);

  const requestSort = (key: TableViewCol) => {
    if (sortConfig?.key === key) {
      setSortConfig({ key, direction: invertDir(sortConfig.direction) });
    } else {
      setSortConfig({ key, direction: "ascending" });
    }
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
            return { name: workstation_name, rs: rows };
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

      <TableContainer overflowX="scroll">
        <Table variant="striped">
          <Thead>
            <Tr>
              {shown.map((col, i) =>
                shownColumns[col] ? (
                  <Th key={i} cursor="pointer" onClick={() => requestSort(col)}>
                    {
                      // We add a blank unicode character to prevent heading
                      // names from changing length causing the table columns to
                      // jump around
                      `${col} ${sortConfig?.key === col ? (sortConfig.direction === "ascending" ? "▲" : "▼") : "⠀"}`
                    }
                  </Th>
                ) : null,
              )}
            </Tr>
          </Thead>
          <Tbody>
            {sortedGroups.map(({ name, rs }, i) => (
              <Row machineName={name} key={i} params={params} row={rs} />
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
): (string | number | null)[] =>
  (group_name === undefined
    ? []
    : // Having to type annotate this `keepIf` is an incredible skill issue
      // Unification is cancelled ig
      [keepIf<string | number>(shownColumns["Group"], group_name)]
  )
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
      keepIf(shownColumns["Memory Total (MB)"], Math.round(gpu.memory_total)),
      keepIf(
        shownColumns["Memory Utilisation (%)"],
        Math.round(gpu.memory_util),
      ),
      keepIf(shownColumns["GPU Utilisation (%)"], Math.round(gpu.gpu_util)),
      keepIf(shownColumns["Memory Used (MB)"], Math.round(gpu.memory_used)),
      keepIf(shownColumns["Fan Speed (%)"], Math.round(gpu.fan_speed)),
      keepIf(shownColumns["GPU Temperature (°C)"], Math.round(gpu.gpu_temp)),
      keepIf(
        shownColumns["Memory Temperature (°C)"],
        Math.round(gpu.memory_temp),
      ),
      keepIf(
        shownColumns["GPU Voltage (mV)"],
        Math.round(gpu.graphics_voltage),
      ),
      keepIf(shownColumns["Power Draw (W)"], Math.round(gpu.power_draw)),
      keepIf(shownColumns["GPU Clock (MHz)"], Math.round(gpu.graphics_clock)),
      keepIf(
        shownColumns["Max GPU Clock (MHz)"],
        Math.round(gpu.max_graphics_clock),
      ),
      keepIf(shownColumns["Memory Clock (MHz)"], Math.round(gpu.memory_clock)),
      keepIf(
        shownColumns["Max Memory Clock (MHz)"],
        Math.round(gpu.max_memory_clock),
      ),
    ]);
