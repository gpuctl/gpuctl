import {
  Button,
  LinkBox,
  LinkOverlay,
  Menu,
  MenuButton,
  MenuItemOption,
  MenuList,
  MenuOptionGroup,
  Table,
  TableContainer,
  Tbody,
  Td,
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
                  const paramsCopy = Object.assign({}, params);
                  paramsCopy.append("selected", workstation_name);
                  return (
                    <LinkBox>
                      <LinkOverlay
                        as={ReactRouterLink}
                        to={{ search: paramsCopy.toString() }}
                      />
                      <Tr key={id}>
                        {tablify(
                          shownColumns,
                          gpu,
                          group_name,
                          workstation_name,
                        ).map((s, i) =>
                          s === null ? null : <Td key={i}>{s}</Td>,
                        )}
                      </Tr>
                    </LinkBox>
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
