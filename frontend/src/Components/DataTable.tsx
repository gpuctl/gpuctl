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
import { isFree } from "../Utils/Utils";
import { useForceUpdate } from "framer-motion";

/*
  machine name - from workstation
  gpu_name: string;
  is_free: Boolean --derived value
  gpu_brand: string;
  driver_ver: string;
  memory_total: number;

  memory_util: number;
  gpu_util: number;
  memory_used: number;
  fan_speed: number;
  gpu_temp: number;

  */
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
  };
  const [shownColumns, setter] = useState(cols);
  const [refresh] = useForceUpdate();

  return (
    <div>
      <Menu closeOnSelect={false}>
        <MenuButton as={Button} colorScheme="blue">
          Columns
        </MenuButton>
        <MenuList>
          <MenuOptionGroup
            type="checkbox"
            defaultValue={Object.keys(shownColumns).filter(
              (key) => shownColumns[key]
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
            <MenuItemOption value="gpu_temp"> GPU temperature </MenuItemOption>
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
            {groups.map(({ name: group_name, workStations }, i) =>
              workStations.map(({ name: workstation_name, gpus }, j) =>
                gpus.map((gpu, k) => {
                  const id =
                    (i * gpus.length * workStations.length +
                      j * gpus.length +
                      k) *
                    12; //size of gpu
                  return (
                    <Tr key={id}>
                      {shownColumns.group ? (
                        <Td key={id}> {group_name}</Td>
                      ) : (
                        null
                      )}
                      {shownColumns.machine_name ? (
                        <Td key={id + 1}> {workstation_name}</Td>
                      ) : (
                        null
                      )}
                      {shownColumns.gpu_name ? (
                        <Td key={id + 2}> {gpu.gpu_name}</Td>
                      ) : (
                        null
                      )}
                      {shownColumns.is_free ? (
                        <Td key={id + 3}> {isFree(gpu).toString()}</Td>
                      ) : (
                        null
                      )}
                      {shownColumns.gpu_brand ? (
                        <Td key={id + 4}> {gpu.gpu_brand}</Td>
                      ) : (
                        null
                      )}
                      {shownColumns.driver_ver ? (
                        <Td key={id + 5}> {gpu.driver_ver}</Td>
                      ) : (
                        null
                      )}
                      {shownColumns.memory_total ? (
                        <Td key={id + 6}> {gpu.memory_total}</Td>
                      ) : (
                        null
                      )}
                      {shownColumns.memory_util ? (
                        <Td key={id + 7}> {gpu.memory_util}</Td>
                      ) : (
                        null
                      )}
                      {shownColumns.gpu_util ? (
                        <Td key={id + 8}> {gpu.gpu_util}</Td>
                      ) : (
                        null
                      )}
                      {shownColumns.memory_used ? (
                        <Td key={id + 9}> {gpu.memory_used}</Td>
                      ) : (
                        null
                      )}
                      {shownColumns.fan_speed ? (
                        <Td key={id + 10}> {gpu.fan_speed}</Td>
                      ) : (
                        null
                      )}
                      {shownColumns.gpu_temp ? (
                        <Td key={id + 11}> {gpu.gpu_temp}</Td>
                      ) : (
                        null
                      )}
                    </Tr>
                  );
                })
              )
            )}
          </Tbody>
        </Table>
      </TableContainer>
    </div>
  );
};
