import { useState } from "react";
import { Box, Button, Heading, HStack, Input } from "@chakra-ui/react";
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
          w="55%"
          placeholder="Hostname (e.g. mira05.doc.ic.ac.uk)"
          onChange={(e) => setHostname(e.target.value)}
          value={hostname}
        ></Input>
        {/* <GroupSelect
          w="100%"
          onChange={(a) => {
            setGroup(a);
          }}
        >
          <AutoCompleteInput
            ref={ref}
            w="100%"
            onChange={(e) => setGroup(e.target.value)}
            placeholder="Group Name (e.g. shared)"
            // onSelect={() => {
            //   setTimeout(() => {
            //     ref.current?.focus();
            //   }, 1);
            //   setShouldFetch(false);
            // }}
            // onBlur={() => {
            //   setShouldFetch(true);
            // }}
            value={group}
          ></AutoCompleteInput>
        </GroupSelect> */}
        <Input
          w="40%"
          placeholder="Group Name (e.g. shared)"
          onChange={(e) => setGroup(e.target.value)}
          value={group}
        ></Input>
        <Button
          w="5%"
          onClick={() => {
            addMachine(hostname, group === "" ? "Shared" : group);
          }}
        >
          Add
        </Button>
      </HStack>
    </Box>
  );
};
