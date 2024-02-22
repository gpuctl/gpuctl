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
import { WorkStationGroup } from "../Data";
import { instKeys } from "../Utils/Utils";
import { useRemoveMachine } from "../Hooks/Hooks";
import { GS } from "../Pages/AdminPanel";

export const GroupInfoManagement = ({
  GroupSelect,
  groups,
}: {
  GroupSelect: GS;
  groups: WorkStationGroup[];
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
            {instKeys(
              groups.flatMap((group) =>
                group.workstations.map((workstation) => (k) => (
                  <Tr key={k}>
                    <Td>{workstation.name}</Td>
                    <EditableField
                      GroupSelect={GroupSelect}
                      workstation={workstation}
                      fieldKey="group"
                      placeholder="unknown"
                    />
                    <EditableField
                      GroupSelect={GroupSelect}
                      workstation={workstation}
                      fieldKey="cpu"
                      placeholder="unknown"
                    />
                    <EditableField
                      GroupSelect={GroupSelect}
                      workstation={workstation}
                      fieldKey="motherboard"
                      placeholder="unknown"
                    />
                    <EditableField
                      GroupSelect={GroupSelect}
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
              ),
            )}
          </Tbody>
        </Table>
      </TableContainer>
    </Box>
  );
};
