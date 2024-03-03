import {
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
} from "@chakra-ui/react";
import { useSearchParams } from "react-router-dom";
import { Graph } from "./Graph";
import { useHistoryStats } from "../Hooks/Hooks";
import { useState } from "react";
import { GraphField } from "../Data";
import { mapSuccess, success, validationElim } from "../Utils/Utils";
import { Text } from "@chakra-ui/react";

const USE_FAKE_STATS = true;

const FAKE_STATS = [
  { x: 1, y: 90 },
  { x: 2, y: 12 },
  { x: 3, y: 34 },
  { x: 4, y: 53 },
  { x: 5, y: 98 },
];

export const WorkstationView = ({ hostname }: { hostname: string }) => {
  const [, setPs] = useSearchParams();
  const [field] = useState<GraphField>(GraphField.GPU_UTIL);

  const historyStats = useHistoryStats(hostname);

  const statsToDisplay = USE_FAKE_STATS
    ? success(FAKE_STATS)
    : mapSuccess(historyStats, (s) =>
        s.map(({ timestamp, sample }) => ({
          x: timestamp,
          // TODO: Support machines with multiple GPUs/no GPUs
          y: sample[0][field],
        })),
      );

  return (
    <Modal
      size="xl"
      isOpen={hostname !== null}
      onClose={() => {
        setPs((ps) => {
          ps.delete("selected");
          return ps;
        });
      }}
    >
      <ModalOverlay bg="blackAlpha.300" backdropFilter="blur(15px)" />
      <ModalContent>
        <ModalHeader>{`${hostname}`}</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          Body!
          {validationElim(statsToDisplay, {
            success: (s) => <Graph width={200} height={200} data={s}></Graph>,
            failure: () => (
              <Text>Failed to fetch historical data for graph!</Text>
            ),
            loading: () => <Text>Fetching data...</Text>,
          })}
        </ModalBody>
      </ModalContent>
    </Modal>
  );
};
