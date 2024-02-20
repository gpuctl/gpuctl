import {
  Button,
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
import { WorkStationGroup } from "../Data";
import { useState } from "react";
import { useForceUpdate } from "framer-motion";
import { isFree } from "../Utils/Stats";

export const TableTab = ({ groups }: { groups: WorkStationGroup[] }) => {
  // default to show group, machine_name, gpu_name, isFree, brand, and memory_total
  const cols: Record<string, Boolean> = {
    group: true,
    machine_name: true,
    gpu_name: true,
    is_free: true,
    gpu_brand: true,
    driver_ver: false,
    memory_total: true,
    memory_util: false,
    gpu_util: false,
    memory_used: false,
    fan_speed: false,
    gpu_temp: false,
    memory_temp: false,
    graphics_voltage: false,
    power_draw: false,
    graphics_clock: false,
    max_graphics_clock: false,
    memory_clock: false,
    max_memory_clock: false,
  };
  const [shownColumns, setter] = useState(cols);
  const [refresh] = useForceUpdate();

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
            <MenuItemOption value="group"> Group </MenuItemOption>
            <MenuItemOption value="machine_name"> Machine name </MenuItemOption>
            <MenuItemOption value="gpu_name"> GPU Id </MenuItemOption>
            <MenuItemOption value="is_free"> Free </MenuItemOption>
            <MenuItemOption value="gpu_brand"> GPU Model </MenuItemOption>
            <MenuItemOption value="driver_ver"> Driver version </MenuItemOption>
            <MenuItemOption value="memory_total"> Memory total </MenuItemOption>
            <MenuItemOption value="memory_util">
              {" "}
              Memory utilisation{" "}
            </MenuItemOption>
            <MenuItemOption value="gpu_util"> GPU utilisation </MenuItemOption>
            <MenuItemOption value="memory_used"> Memory used </MenuItemOption>
            <MenuItemOption value="fan_speed"> Fan speed </MenuItemOption>
            <MenuItemOption value="gpu_temp">
              {" "}
              GPU temperature (C){" "}
            </MenuItemOption>
            <MenuItemOption value="memory_temp">
              {" "}
              Memory temperature (C){" "}
            </MenuItemOption>
            <MenuItemOption value="graphics_voltage"> Voltage </MenuItemOption>
            <MenuItemOption value="power_draw"> Power draw </MenuItemOption>
            <MenuItemOption value="graphics_clock">
              {" "}
              GPU Clock (Mhz){" "}
            </MenuItemOption>
            <MenuItemOption value="max_graphics_clock">
              {" "}
              Max GPU Clock (Mhz){" "}
            </MenuItemOption>
            <MenuItemOption value="memory_clock">
              {" "}
              Memory Clock (Mhz){" "}
            </MenuItemOption>
            <MenuItemOption value="max_memory_clock">
              {" "}
              Max Memory Clock (Mhz){" "}
            </MenuItemOption>
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
                    <Tr key={id}>
                      {shownColumns.group ? (
                        <Td key={id}> {group_name}</Td>
                      ) : null}
                      {shownColumns.machine_name ? (
                        <Td key={id + 1}> {workstation_name}</Td>
                      ) : null}
                      {shownColumns.gpu_name ? (
                        <Td key={id + 2}> {gpu.gpu_name}</Td>
                      ) : null}
                      {shownColumns.is_free ? (
                        <Td key={id + 3}> {isFree(gpu).toString()}</Td>
                      ) : null}
                      {shownColumns.gpu_brand ? (
                        <Td key={id + 4}> {gpu.gpu_brand}</Td>
                      ) : null}
                      {shownColumns.driver_ver ? (
                        <Td key={id + 5}> {gpu.driver_ver}</Td>
                      ) : null}
                      {shownColumns.memory_total ? (
                        <Td key={id + 6}> {Math.round(gpu.memory_total)}</Td>
                      ) : null}
                      {shownColumns.memory_util ? (
                        <Td key={id + 7}> {Math.round(gpu.memory_util)}</Td>
                      ) : null}
                      {shownColumns.gpu_util ? (
                        <Td key={id + 8}> {Math.round(gpu.gpu_util)}</Td>
                      ) : null}
                      {shownColumns.memory_used ? (
                        <Td key={id + 9}> {Math.round(gpu.memory_used)}</Td>
                      ) : null}
                      {shownColumns.fan_speed ? (
                        <Td key={id + 10}> {Math.round(gpu.fan_speed)}</Td>
                      ) : null}
                      {shownColumns.gpu_temp ? (
                        <Td key={id + 11}> {Math.round(gpu.gpu_temp)}</Td>
                      ) : null}
                      {shownColumns.memory_temp ? (
                        <Td key={id + 12}> {Math.round(gpu.memory_temp)}</Td>
                      ) : null}
                      {shownColumns.graphics_voltage ? (
                        <Td key={id + 13}>
                          {" "}
                          {Math.round(gpu.graphics_voltage)}
                        </Td>
                      ) : null}
                      {shownColumns.graphics_clock ? (
                        <Td key={id + 14}> {Math.round(gpu.graphics_clock)}</Td>
                      ) : null}
                      {shownColumns.graphics_clock ? (
                        <Td key={id + 15}> {Math.round(gpu.graphics_clock)}</Td>
                      ) : null}
                      {shownColumns.max_graphics_clock ? (
                        <Td key={id + 16}>
                          {" "}
                          {Math.round(gpu.max_graphics_clock)}
                        </Td>
                      ) : null}
                      {shownColumns.memory_clock ? (
                        <Td key={id + 17}> {Math.round(gpu.memory_clock)}</Td>
                      ) : null}
                      {shownColumns.max_memory_clock ? (
                        <Td key={id + 18}>
                          {" "}
                          {Math.round(gpu.max_memory_clock)}
                        </Td>
                      ) : null}
                    </Tr>
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
