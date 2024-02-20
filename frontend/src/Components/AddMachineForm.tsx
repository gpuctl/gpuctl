import React, { useState } from "react";
import { Box, Button, Heading, HStack, Input } from "@chakra-ui/react";
import {
  AutoComplete,
  AutoCompleteInput,
  AutoCompleteItem,
  AutoCompleteList,
} from "@choc-ui/chakra-autocomplete";
import { useAddMachine } from "../Hooks/hooks";
import { WorkStationGroup } from "../Data";

type AddMachineFormProps = {
  groups: WorkStationGroup[];
};

export const AddMachineForm: React.FC<AddMachineFormProps> = ({ groups }) => {
  const [hostname, setHostname] = useState<string>("");
  const [group, setGroup] = useState<string>("");
  const addMachine = useAddMachine();

  return (
    <Box w="100%">
      <Heading size="lg">Add a Machine:</Heading>
      <HStack w="100%" mt={4}>
        <Input
          w="50%"
          placeholder="Hostname (e.g. mira05.doc.ic.ac.uk)"
          onChange={(e) => setHostname(e.target.value)}
        />
        <AutoComplete
          openOnFocus
          creatable
          onChange={(e) => setGroup(e.target.value)}
        >
          <AutoCompleteInput placeholder="Group Name (e.g. shared)" />
          <AutoCompleteList>
            {groups.map((g, i) => (
              <AutoCompleteItem key={i} value={g.name} label={g.name} />
            ))}
          </AutoCompleteList>
        </AutoComplete>
        <Button onClick={() => addMachine(hostname, group)}>Add</Button>
      </HStack>
    </Box>
  );
};
