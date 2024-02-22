import React from "react";
import {
  Box,
  Button,
  Heading,
  Table,
  TableContainer,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import { EditableField } from "./EditableFields";
import { useRemoveMachine } from "../Hooks/Hooks";
import { WorkStationGroup } from "../Data";

type GroupInfoManagementProps = {
  groups: WorkStationGroup[];
};

export const GroupInfoManagement: React.FC<GroupInfoManagementProps> = ({
  groups,
}) => {
  const removeMachine = useRemoveMachine();

  return (
    <Box w="100%">
      <Heading size="lg">Group & Info Management:</Heading>
      <TableContainer mt={4}>
        <Table variant="striped">
          <Thead>
            <Tr>
              <Th>Hostname</Th>
              <Th>Group</Th>
              <Th>CPU</Th>
              <Th>Motherboard</Th>
              <Th>Notes</Th>
              <Th>Action</Th>
            </Tr>
          </Thead>
          <Tbody>
            {groups.flatMap((group) =>
              group.workstations.map((workstation) => (
                <Tr key={workstation.name}>
                  <Td>{workstation.name}</Td>
                  <EditableField
                    workstation={workstation}
                    fieldKey="group"
                    placeholder="unknown"
                  />
                  <EditableField
                    workstation={workstation}
                    fieldKey="cpu"
                    placeholder="unknown"
                  />
                  <EditableField
                    workstation={workstation}
                    fieldKey="motherboard"
                    placeholder="unknown"
                  />
                  <EditableField
                    workstation={workstation}
                    fieldKey="notes"
                    placeholder="none"
                  />
                  <Td>
                    <Button
                      colorScheme="red"
                      onClick={() => removeMachine(workstation.name)}
                    >
                      Remove
                    </Button>
                  </Td>
                </Tr>
              )),
            )}
          </Tbody>
        </Table>
      </TableContainer>
    </Box>
  );
};
