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
              {Object.keys(shownColumns).map((col) => {
                if (shownColumns[col]) return <Th key={col}>{col}</Th>;
                else return <></>;
              })}
            </Tr>
          </Thead>
          <Tbody>
            {groups.map(({ name: group_name, workStations }) =>
              workStations.map(({ name: workstation_name, gpus }) =>
                gpus.map((gpu) => (
                  <Tr key={gpu.gpu_name}>
                    {shownColumns.group ? <Td> {group_name}</Td> : <></>}
                    {shownColumns.machine_name ? (
                      <Td> {workstation_name}</Td>
                    ) : (
                      <></>
                    )}
                    {shownColumns.gpu_name ? <Td> {gpu.gpu_name}</Td> : <></>}
                    {shownColumns.is_free ? (
                      <Td> {isFree(gpu).toString()}</Td>
                    ) : (
                      <></>
                    )}
                    {shownColumns.gpu_brand ? <Td> {gpu.gpu_brand}</Td> : <></>}
                    {shownColumns.driver_ver ? (
                      <Td> {gpu.driver_ver}</Td>
                    ) : (
                      <></>
                    )}
                    {shownColumns.memory_total ? (
                      <Td> {gpu.memory_total}</Td>
                    ) : (
                      <></>
                    )}
                    {shownColumns.memory_util ? (
                      <Td> {gpu.memory_util}</Td>
                    ) : (
                      <></>
                    )}
                    {shownColumns.gpu_util ? <Td> {gpu.gpu_util}</Td> : <></>}
                    {shownColumns.memory_used ? (
                      <Td> {gpu.memory_used}</Td>
                    ) : (
                      <></>
                    )}
                    {shownColumns.fan_speed ? <Td> {gpu.fan_speed}</Td> : <></>}
                    {shownColumns.gpu_temp ? <Td> {gpu.gpu_temp}</Td> : <></>}
                  </Tr>
                ))
              )
            )}
          </Tbody>
        </Table>
      </TableContainer>
    </div>
  );
};
