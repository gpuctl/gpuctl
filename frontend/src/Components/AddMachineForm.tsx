import { useState } from "react";
import { Box, Button, Heading, HStack, Input, Spinner } from "@chakra-ui/react";
import { useAddMachine } from "../Hooks/Hooks";
import { WorkStationGroup } from "../Data";
import { GS } from "../Pages/AdminPanel";
import { useStats } from "../Providers/FetchProvider";
import { Validated, validatedElim } from "../Utils/Utils";

export const AddMachineForm = ({
  GroupSelect,
  groups,
}: {
  GroupSelect: GS;
  groups: WorkStationGroup[];
}) => {
  const [hostname, setHostname] = useState<string>("");
  const [group, setGroup] = useState<string>("");
  const [spinner, setSpinner] = useState<boolean>(false);
  const { onNextFetch } = useStats();
  const addMachine = useAddMachine(() => {
    setSpinner(false);
  });

  // Auto-complete code removed in
  // https://github.com/gpuctl/gpuctl/pull/220
  // Go back there if you want to re-implement it
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
        <Input
          w="40%"
          placeholder="Group Name (e.g. shared)"
          onChange={(e) => setGroup(e.target.value)}
          value={group}
        ></Input>
        <Button
          w="5%"
          onClick={() => {
            setSpinner(true);
            const cb = (stats: Validated<WorkStationGroup[]>) =>
              validatedElim(stats, {
                success: (s) => {
                  if (
                    !s
                      .flatMap((g) => g.workstations.map((w) => w.name))
                      .some((x) => x === hostname)
                  )
                    return false;

                  setSpinner(false);
                  return true;
                },
                failure: () => {
                  return false;
                },
              });

            onNextFetch(cb);
            addMachine(hostname, group === "" ? "Shared" : group);
          }}
        >
          Add
        </Button>
        {spinner ? <Spinner /> : <></>}
      </HStack>
    </Box>
  );
};
