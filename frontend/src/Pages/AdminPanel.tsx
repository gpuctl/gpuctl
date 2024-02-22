import { Input } from "@chakra-ui/input";
import {
  Editable,
  EditableInput,
  EditablePreview,
  StyleProps,
  Table,
  TableContainer,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import { WorkStationGroup } from "../Data";

import { STATS_PATH } from "../Config/Paths";
import { Box, HStack, Heading, VStack } from "@chakra-ui/layout";
import { Button } from "@chakra-ui/button";
import { PropsWithChildren, useState } from "react";
import {
  AutoComplete,
  AutoCompleteInput,
  AutoCompleteItem,
  AutoCompleteList,
} from "@choc-ui/chakra-autocomplete";
import { instKeys } from "../Utils/Utils";
import { useAuth } from "../Providers/AuthProvider";

export const ADMIN_PATH = "/admin";
const ADD_MACHINE_URL = "/add_workstation";
const REMOVE_MACHINE_URL = "/rm_workstation";

const useAddMachine = () => {
  const { useAuthFetch } = useAuth();
  const [, addMachineAuth] = useAuthFetch(ADD_MACHINE_URL);
  return (hostname: string, group: string) =>
    addMachineAuth({
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ hostname, group }),
    });
};

const useRemoveMachine = () => {
  const { useAuthFetch } = useAuth();
  const [, addMachineAuth] = useAuthFetch(REMOVE_MACHINE_URL);
  return (hostname: string) =>
    addMachineAuth({
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ hostname }),
    });
};

type ModifyData = {
  group: string | null;
  motherboard: string | null;
  cpu: string | null;
  notes: string | null;
};

const useModifyInfo = () => {
  const { useAuthFetch } = useAuth();
  const [, addMachineAuth] = useAuthFetch(STATS_PATH + "/modify");
  return (hostname: string, mod: ModifyData) =>
    addMachineAuth({
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ hostname, ...mod }),
    });
};

export const AdminPanel = ({ groups }: { groups: WorkStationGroup[] }) => {
  const [hostname, setHostname] = useState("");
  const [group, setGroup] = useState("");

  const pickCol = (value: string) => (value ? "gray.600" : "gray.300");

  const addMachine = useAddMachine();
  const removeMachine = useRemoveMachine();
  const modifyInfo = useModifyInfo();

  const GroupSelect = ({
    onChange,
    ...props
  }: PropsWithChildren & StyleProps & { onChange: (v: string) => void }) => (
    <AutoComplete
      openOnFocus
      creatable
      onChange={(a) => onChange(a)}
      values={groups.map((g) => g.name)}
    >
      <Box {...props} />
      <AutoCompleteList maxHeight="200px" overflow="scroll">
        {groups.map((g, i) => (
          <AutoCompleteItem key={i} value={g.name}></AutoCompleteItem>
        ))}
      </AutoCompleteList>
    </AutoComplete>
  );

  return (
    <VStack padding={10} spacing={10}>
      <Box w="100%" textAlign="left">
        <Heading size="lg">Add a Machine:</Heading>
      </Box>
      <HStack w="100%">
        <Input
          w="50%"
          placeholder="Hostname (e.g. mira05.doc.ic.ac.uk)"
          onChange={(e) => setHostname(e.target.value)}
        ></Input>
        <GroupSelect w="100%" onChange={setGroup}>
          <AutoCompleteInput
            w="100%"
            placeholder="Group Name (e.g. shared)"
          ></AutoCompleteInput>
        </GroupSelect>

        <Button
          w="5%"
          onClick={() => {
            addMachine(hostname, group);
          }}
        >
          Add
        </Button>
      </HStack>
      <Box w="100%" textAlign="left">
        <Heading size="lg">Group & Info Management:</Heading>
      </Box>
      <TableContainer>
        <Table variant="striped">
          <Thead>
            <Tr>
              <Th> Hostname </Th>
              <Th> Group </Th>
              <Th> CPU </Th>
              <Th> Motherboard </Th>
              <Th> Notes </Th>
            </Tr>
          </Thead>
          <Tbody>
            {instKeys(
              groups.flatMap(({ name: group, workstations }, i) =>
                workstations.map((workstation, j) => (k: number) => (
                  <Tr key={k}>
                    <Td> {workstation.name} </Td>
                    <Td>
                      <GroupSelect
                        onChange={(s) => {
                          modifyInfo(workstation.name, {
                            group: s,
                            cpu: null,
                            motherboard: null,
                            notes: null,
                          });
                        }}
                      >
                        <AutoCompleteInput placeholder="Unknown"></AutoCompleteInput>
                      </GroupSelect>
                    </Td>
                    <Td>
                      {" "}
                      {/* Do we need this empty string? */}
                      <Editable
                        placeholder="Unknown"
                        textColor={pickCol(workstation.cpu)}
                        defaultValue={workstation.cpu}
                        onSubmit={(s) =>
                          modifyInfo(workstation.name, {
                            group: null,
                            cpu: s,
                            motherboard: null,
                            notes: null,
                          })
                        }
                      >
                        <EditablePreview />
                        <EditableInput textColor={"gray.600"} />
                      </Editable>{" "}
                    </Td>
                    <Td>
                      {" "}
                      <Editable
                        placeholder="Unknown"
                        defaultValue={workstation.motherboard}
                        textColor={pickCol(workstation.motherboard)}
                        onSubmit={(s) =>
                          modifyInfo(workstation.name, {
                            group: null,
                            cpu: null,
                            motherboard: s,
                            notes: null,
                          })
                        }
                      >
                        <EditablePreview />
                        <EditableInput textColor={"gray.600"} />
                      </Editable>{" "}
                    </Td>
                    <Td>
                      {" "}
                      <Editable
                        placeholder="None"
                        defaultValue={workstation.notes}
                        textColor={pickCol(workstation.notes)}
                        onSubmit={(s) =>
                          modifyInfo(workstation.name, {
                            group: null,
                            cpu: null,
                            motherboard: null,
                            notes: s,
                          })
                        }
                      >
                        <EditablePreview />
                        <EditableInput textColor={"gray.600"} />
                      </Editable>
                    </Td>
                    <Td>
                      <Button
                        bgColor={"red.300"}
                        onClick={() => {
                          removeMachine(hostname);
                        }}
                      >
                        Kill Satellite
                      </Button>
                    </Td>
                  </Tr>
                )),
              ),
            )}
          </Tbody>
        </Table>
      </TableContainer>
    </VStack>
  );
};
