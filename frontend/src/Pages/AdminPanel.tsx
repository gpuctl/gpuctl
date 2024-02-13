import { Input } from "@chakra-ui/input";
import { API_URL } from "../App";
import {
  Editable,
  EditableInput,
  EditablePreview,
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
import { useState } from "react";
import {
  AutoComplete,
  AutoCompleteInput,
  AutoCompleteItem,
  AutoCompleteList,
} from "@choc-ui/chakra-autocomplete";

export const ADMIN_PATH = "/admin";
const ADD_MACHINE_URL = "/add_workstation";
const REMOVE_MACHINE_URL = "/rm_workstation";

const addMachine = async (hostname: string, group: string) => {
  const resp = await fetch(API_URL + ADMIN_PATH + ADD_MACHINE_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ hostname, group }),
  });
  if (resp.ok) {
    console.log("Success!");
  } else if (resp.status === 401) {
    console.log("Auth Error!");
  } else {
    console.log("Unknown Error!");
  }
};

const removeMachine = async (hostname: string) => {
  const resp = await fetch(API_URL + ADMIN_PATH + REMOVE_MACHINE_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ hostname }),
  });
  if (resp.ok) {
    console.log("Success!");
  } else if (resp.status === 401) {
    console.log("Auth Error!");
  } else {
    console.log("Unknown Error!");
  }
};

type ModifyData = {
  group: string | null;
  motherboard: string | null;
  cpu: string | null;
  notes: string | null;
};

const modifyInfo = (hostname: string, mod: ModifyData) => {
  fetch(API_URL + ADMIN_PATH + STATS_PATH + "/modify", {
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
        <AutoComplete
          openOnFocus
          creatable
          onChange={(e) => setGroup(e.target.value)}
          values={groups.map((g) => g.name)}
        >
          <AutoCompleteInput
            w="45%"
            placeholder="Group Name (e.g. shared)"
          ></AutoCompleteInput>
          <AutoCompleteList maxHeight="200px" overflow="scroll">
            {groups.map((g, i) => (
              <AutoCompleteItem key={i} value={g.name}></AutoCompleteItem>
            ))}
          </AutoCompleteList>
        </AutoComplete>
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
              <Th key={0}> Hostname </Th>
              <Th key={1}> Group </Th>
              <Th key={2}> CPU </Th>
              <Th key={3}> Motherboard </Th>
              <Th key={4}> Notes </Th>
            </Tr>
          </Thead>
          <Tbody>
            {groups.map(({ name: group, workStations }, i) => {
              return workStations.map((workStation, j) => {
                const id = i * j + j;
                return (
                  <Tr key={id}>
                    <Td> {workStation.name} </Td>
                    <Td>
                      <Editable
                        placeholder="Unknown"
                        defaultValue={group}
                        textColor={pickCol(group)}
                        onSubmit={(s) =>
                          modifyInfo(workStation.name, {
                            group: s,
                            cpu: null,
                            motherboard: null,
                            notes: null,
                          })
                        }
                      >
                        <EditablePreview />
                        <EditableInput textColor={"gray.600"} />
                      </Editable>
                    </Td>
                    <Td>
                      {" "}
                      {/* Do we need this empty string? */}
                      <Editable
                        placeholder="Unknown"
                        textColor={pickCol(workStation.cpu)}
                        defaultValue={workStation.cpu}
                        onSubmit={(s) =>
                          modifyInfo(workStation.name, {
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
                        defaultValue={workStation.motherboard}
                        textColor={pickCol(workStation.motherboard)}
                        onSubmit={(s) =>
                          modifyInfo(workStation.name, {
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
                        defaultValue={workStation.notes}
                        textColor={pickCol(workStation.notes)}
                        onSubmit={(s) =>
                          modifyInfo(workStation.name, {
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
                );
              });
            })}
          </Tbody>
        </Table>
      </TableContainer>
    </VStack>
  );
};
