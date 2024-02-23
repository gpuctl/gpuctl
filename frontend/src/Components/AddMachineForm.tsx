import { useState } from "react";
import { Box, Button, Heading, HStack, Input } from "@chakra-ui/react";
import { AutoCompleteInput } from "@choc-ui/chakra-autocomplete";
import { useAddMachine } from "../Hooks/Hooks";
import { WorkStationGroup } from "../Data";
import { GS } from "../Pages/AdminPanel";

export const AddMachineForm = ({
  GroupSelect,
  groups,
}: {
  GroupSelect: GS;
  groups: WorkStationGroup[];
}) => {
  const [hostname, setHostname] = useState<string>("");
  const [group, setGroup] = useState<string>("");
  const addMachine = useAddMachine();

  return (
    <Box w="100%">
      <Heading size="lg">Add a Machine:</Heading>
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
    </Box>
  );
};
