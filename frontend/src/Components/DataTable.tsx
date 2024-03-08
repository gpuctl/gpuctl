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
import { GPUStats, WorkStationGroup, WorkStationData } from "../Data";
import { useEffect, useMemo, useState } from "react";
import { useForceUpdate } from "framer-motion";
import { keepIf } from "../Utils/Utils";

import { Link as ReactRouterLink, useSearchParams } from "react-router-dom";
import { sort } from "d3";

export const GPU_FIELDS = {
  "GPU Name": "gpu_name",
  Free: "free",
  "GPU Brand": "gpu_brand",
  "Driver Version": "driver_ver",
  "Memory Total": "memory_total",
  "Memory Utilisation": "memory_util",
  "GPU Utilisation": "gpu_util",
  "Memory Used": "memory_used",
  "Fan Speed": "fan_speed",
  "GPU Temperature": "gpu_temp",
  "Memory Temperature": "memory_temp",
  "GPU Voltage": "graphics_voltage",
  "Power Draw": "power_draw",
  "GPU Clock": "graphics_clock",
  "Max GPU Clock": "max_graphics_clock",
  "Memory Clock": "memory_clock",
  "Max Memory Clock": "max_memory_clock",
} as const;

const map_idx = (key: string) => {
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

  newParams.append("selected", row[map_idx("Machine Name")]!.toString());

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
        {" "}
        {row.map((s, i) =>
          s === null ? null : (
            <Td>
              <Text textDecoration={hover ? "underline" : ""}>{s}</Text>
            </Td>
          ),
        )}
      </Link>
    </Tr>
  );
};

export const TableTab = ({ groups }: { groups: WorkStationGroup[] }) => {
  // default to show group, machine_name, gpu_name, isFree, brand, and memory_total
  const SHOWN_COLS: {
    [_ in keyof typeof GPU_FIELDS | "Group" | "Machine Name"]: boolean;
  } = {
    Group: true,
    "Machine Name": true,
    "GPU Name": true,
    Free: true,
    "GPU Brand": true,
    "Driver Version": false,
    "Memory Total": true,
    "Memory Utilisation": false,
    "GPU Utilisation": false,
    "Memory Used": false,
    "Fan Speed": false,
    "GPU Temperature": false,
    "Memory Temperature": false,
    "GPU Voltage": false,
    "Power Draw": false,
    "GPU Clock": false,
    "Max GPU Clock": false,
    "Memory Clock": false,
    "Max Memory Clock": false,
  };

  const [shownColumns, setter] = useState<Record<string, boolean>>(SHOWN_COLS);
  const [refresh] = useForceUpdate();
  const [params] = useSearchParams();
  const [sortConfig, setSortConfig] = useState({
    key: "name",
    direction: "ascending",
  });

  const [rows, setRows] = useState<(string | number | null)[][]>([]);

  const sortedGroups = useMemo(() => {
    let sortableItems = [...rows];
    if (sortConfig !== null) {
      console.log(sortConfig);
      sortableItems.sort((a, b) => {
        if (a[map_idx(sortConfig.key)]! < b[map_idx(sortConfig.key)]!) {
          return sortConfig.direction === "ascending" ? -1 : 1;
        }
        if (a[map_idx(sortConfig.key)]! > b[map_idx(sortConfig.key)]!) {
          return sortConfig.direction === "ascending" ? 1 : -1;
        }
        return 0;
      });
    }

    return sortableItems;
  }, [groups, sortConfig]);

  const requestSort = (key: string) => {
    let direction = "ascending";
    if (sortConfig.key === key && sortConfig.direction === "ascending") {
      direction = "descending";
    }
    setSortConfig({ key, direction });
  };


  useEffect(() => {
    setRows(
      groups.flatMap(({ name: group_name, workstations }, i) =>
        workstations.flatMap(({ name: workstation_name, gpus }, j) =>
          gpus.map((gpu, k) => {
            const id =
              (i * gpus.length * workstations.length + j * gpus.length + k) *
              19; //size of gpu

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
  }, [groups]);

  return (
    <div>
      <Menu closeOnSelect={false}>
        <MenuButton as={Button} colorScheme="blue">
          Columns
        </MenuButton>
        <MenuList overflowY="scroll" maxHeight="200">
          <MenuOptionGroup
            type="checkbox"
            defaultValue={Object.keys(shownColumns).filter(
              (key) => shownColumns[key],
            )}
            onChange={(props) => {
              Object.keys(shownColumns).forEach((col) => {
                shownColumns[col] = props.includes(col);
              });
              setter(shownColumns);
              refresh();
            }}
          >
            {Object.keys(shownColumns).map((col, i) => {
              return (
                <MenuItemOption value={col} key={i}>
                  {" "}
                  {col}{" "}
                </MenuItemOption>
              );
            })}
          </MenuOptionGroup>
        </MenuList>
      </Menu>

      <TableContainer>
        <Table variant="striped">
          <Thead>
            <Tr>
              {Object.keys(shownColumns).map((col, i) => {
                if (shownColumns[col])
                  return (
                    <Th
                      key={i}
                      cursor="pointer"
                      onClick={() => requestSort(col)}
                    >
                      {col}
                    </Th>
                  );
                else return null;
              })}
            </Tr>
          </Thead>
          <Tbody key={1}>
            {sortedGroups.map((row) => (
              <Row params={params} row={row} />
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
        shownColumns["Memory Total"],
        Math.round(gpu.memory_total).toString(),
      ),
      keepIf(
        shownColumns["Memory Utilisation"],
        Math.round(gpu.memory_util).toString(),
      ),
      keepIf(
        shownColumns["GPU Utilisation"],
        Math.round(gpu.gpu_util).toString(),
      ),
      keepIf(
        shownColumns["Memory Used"],
        Math.round(gpu.memory_used).toString(),
      ),
      keepIf(shownColumns["Fan Speed"], Math.round(gpu.fan_speed).toString()),
      keepIf(
        shownColumns["GPU Temperature"],
        Math.round(gpu.gpu_temp).toString(),
      ),
      keepIf(
        shownColumns["Memory Temperature"],
        Math.round(gpu.memory_temp).toString(),
      ),
      keepIf(
        shownColumns["GPU Voltage"],
        Math.round(gpu.graphics_voltage).toString(),
      ),
      keepIf(shownColumns["Power Draw"], Math.round(gpu.power_draw).toString()),
      keepIf(
        shownColumns["GPU Clock"],
        Math.round(gpu.graphics_clock).toString(),
      ),
      keepIf(
        shownColumns["Max GPU Clock"],
        Math.round(gpu.max_graphics_clock).toString(),
      ),
      keepIf(
        shownColumns["Memory Clock"],
        Math.round(gpu.memory_clock).toString(),
      ),
      keepIf(
        shownColumns["Max Memory Clock"],
        Math.round(gpu.max_memory_clock).toString(),
      ),
    ]);
