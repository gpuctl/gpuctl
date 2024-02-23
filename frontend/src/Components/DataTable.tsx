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

export const TableTab = ({ groups }: { groups: WorkStationGroup[] }) => {
  // default to show group, machine_name, gpu_name, isFree, brand, and memory_total
  const cols: Record<string, Boolean> = {
    Group: true,
    "Machine name": true,
    "GPU name": true,
    Free: true,
    "GPU brand": true,
    "Driver version": false,
    "Memory total": true,
    "Memory utilisation": false,
    "GPU utililisation": false,
    "Memory used": false,
    "Fan speed": false,
    "GPU temperature": false,
    "Memory temperature": false,
    "GPU voltage": false,
    "Power draw": false,
    "GPU clock": false,
    "Max GPU clock": false,
    "Memory clock": false,
    "Max memory clock": false,
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
                    <Tr key={id}>
                      {shownColumns["Group"] ? (
                        <Td key={id}> {group_name}</Td>
                      ) : null}
                      {shownColumns["Machine name"] ? (
                        <Td key={id + 1}> {workstation_name}</Td>
                      ) : null}
                      {shownColumns["GPU name"] ? (
                        <Td key={id + 2}> {gpu.gpu_name}</Td>
                      ) : null}
                      {shownColumns["Free"] ? (
                        <Td key={id + 3}> {gpu.gpu_name ? "❌" : "✅"}</Td>
                      ) : null}
                      {shownColumns["GPU brand"] ? (
                        <Td key={id + 4}> {gpu.gpu_brand}</Td>
                      ) : null}
                      {shownColumns["Driver version"] ? (
                        <Td key={id + 5}> {gpu.driver_ver}</Td>
                      ) : null}
                      {shownColumns["Memory total"] ? (
                        <Td key={id + 6}> {Math.round(gpu.memory_total)}</Td>
                      ) : null}
                      {shownColumns["Memory utilisation"] ? (
                        <Td key={id + 7}> {Math.round(gpu.memory_util)}</Td>
                      ) : null}
                      {shownColumns["GPU utilisation"] ? (
                        <Td key={id + 8}> {Math.round(gpu.gpu_util)}</Td>
                      ) : null}
                      {shownColumns["Memory used"] ? (
                        <Td key={id + 9}> {Math.round(gpu.memory_used)}</Td>
                      ) : null}
                      {shownColumns["Fan speed"] ? (
                        <Td key={id + 10}> {Math.round(gpu.fan_speed)}</Td>
                      ) : null}
                      {shownColumns["GPU temperature"] ? (
                        <Td key={id + 11}> {Math.round(gpu.gpu_temp)}</Td>
                      ) : null}
                      {shownColumns["Memory temperature"] ? (
                        <Td key={id + 12}> {Math.round(gpu.memory_temp)}</Td>
                      ) : null}
                      {shownColumns["GPU voltage"] ? (
                        <Td key={id + 13}>
                          {Math.round(gpu.graphics_voltage)}
                        </Td>
                      ) : null}
                      {shownColumns["Power draw"] ? (
                        <Td key={id + 14}> {Math.round(gpu.power_draw)}</Td>
                      ) : null}
                      {shownColumns["GPU clock"] ? (
                        <Td key={id + 15}> {Math.round(gpu.graphics_clock)}</Td>
                      ) : null}
                      {shownColumns["Max GPU clock"] ? (
                        <Td key={id + 16}>
                          {Math.round(gpu.max_graphics_clock)}
                        </Td>
                      ) : null}
                      {shownColumns["Memory clock"] ? (
                        <Td key={id + 17}> {Math.round(gpu.memory_clock)}</Td>
                      ) : null}
                      {shownColumns["Max memory clock"] ? (
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
