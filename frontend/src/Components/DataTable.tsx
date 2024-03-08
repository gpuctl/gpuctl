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
import { useState } from "react";
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

const Row = ({
  params,
  workstation_name,
  shownColumns,
  gpu,
  group_name,
}: {
  params: URLSearchParams;
  workstation_name: string;
  shownColumns: Record<string, boolean>;
  gpu: GPUStats;
  group_name: string;
}) => {
  const [hover, setHover] = useState(false);
  const newParams = new URLSearchParams(
    Object.fromEntries(Array.from(params.entries())),
  );
  newParams.append("selected", workstation_name);

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
        {tablify(shownColumns, gpu, group_name, workstation_name).map((s, i) =>
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
  const [refresh] = useForceUpdate();
  const [params] = useSearchParams();

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
                if (shownColumns[col]) return <Th key={i}>{col}</Th>;
                else return null;
              })}
            </Tr>
          </Thead>
          <Tbody key={1}>
            {groups.map(({ name: group_name, workstations }, i) =>
              workstations.map(({ name: workstation_name, gpus }, j) =>
                gpus.map((gpu, k) => {
                  const id =
                    (i * gpus.length * workstations.length +
                      j * gpus.length +
                      k) *
                    19; //size of gpu
                  return (
                    <Row
                      key={id}
                      params={params}
                      workstation_name={workstation_name}
                      shownColumns={shownColumns}
                      gpu={gpu}
                      group_name={group_name}
                    />
                  );
                }),
              ),
            )}
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
