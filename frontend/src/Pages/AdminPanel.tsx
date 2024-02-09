import { Input } from "@chakra-ui/input";
import { API_URL, AuthToken } from "../App";
import {
  Editable,
  EditableInput,
  EditableTextarea,
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
import { Box, Center, HStack, Heading, VStack } from "@chakra-ui/layout";
import { Button } from "@chakra-ui/button";
import { useState } from "react";
import { gracefulify } from "graceful-fs";

export const ADMIN_PATH = "/admin";
const ADD_MACHINE_URL = "/add_workstation";

const addMachine = async (
  hostname: string,
  group: string,
  token: AuthToken,
) => {
  const resp = await fetch(API_URL + ADMIN_PATH + ADD_MACHINE_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token.token}`,
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
  // We should probably await the response to give feedback on whether adding
  // the machine was successful...
};

type ModifyData = {
  group: string | null;
  motherboard: string | null;
  cpu: string | null;
  notes: string | null;
};

const modifyInfo = (hostname: string, mod: ModifyData, token: AuthToken) => {
  fetch(API_URL + ADMIN_PATH + STATS_PATH + "/modify", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token.token}`,
    },
    body: JSON.stringify({ hostname, ...mod }),
  });
};

export const AdminPanel = ({
  token,
  groups,
}: {
  token: AuthToken;
  groups: WorkStationGroup[];
}) => {
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
        <Input
          w="45%"
          placeholder="Group Name (e.g. shared)"
          onChange={(e) => setGroup(e.target.value)}
        ></Input>
        <Button
          w="5%"
          onClick={() => {
            addMachine(hostname, group, token);
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
                const id = (i * j + j) * 5;
                return (
                  <Tr>
                    <Th key={id}> {workStation.name} </Th>
                    <Th key={id + 1}>
                      <Editable
                        placeholder="Unknown"
                        defaultValue={group}
                        textColor={pickCol(group)}
                        onSubmit={(s) =>
                          modifyInfo(
                            workStation.name,
                            {
                              group: s,
                              cpu: null,
                              motherboard: null,
                              notes: null,
                            },
                            token,
                          )
                        }
                      >
                        <EditablePreview />
                        <EditableInput textColor={"gray.600"} />
                      </Editable>{" "}
                    </Th>
                    <Th key={id + 2}>
                      {" "}
                      <Editable
                        placeholder="Unknown"
                        textColor={pickCol(workStation.cpu)}
                        defaultValue={workStation.cpu}
                        onSubmit={(s) =>
                          modifyInfo(
                            workStation.name,
                            {
                              group: null,
                              cpu: s,
                              motherboard: null,
                              notes: null,
                            },
                            token,
                          )
                        }
                      >
                        <EditablePreview />
                        <EditableInput textColor={"gray.600"} />
                      </Editable>{" "}
                    </Th>
                    <Th key={id + 3}>
                      {" "}
                      <Editable
                        placeholder="Unknown"
                        defaultValue={workStation.motherboard}
                        textColor={pickCol(workStation.motherboard)}
                        onSubmit={(s) =>
                          modifyInfo(
                            workStation.name,
                            {
                              group: null,
                              cpu: null,
                              motherboard: s,
                              notes: null,
                            },
                            token,
                          )
                        }
                      >
                        <EditablePreview />
                        <EditableInput textColor={"gray.600"} />
                      </Editable>{" "}
                    </Th>
                    <Th key={id + 4}>
                      {" "}
                      <Editable
                        placeholder="None"
                        defaultValue={workStation.notes}
                        textColor={pickCol(workStation.notes)}
                        onSubmit={(s) =>
                          modifyInfo(
                            workStation.name,
                            {
                              group: null,
                              cpu: null,
                              motherboard: null,
                              notes: s,
                            },
                            token,
                          )
                        }
                      >
                        <EditablePreview />
                        <EditableInput textColor={"gray.600"} />
                      </Editable>
                    </Th>
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
