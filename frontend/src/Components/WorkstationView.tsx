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

export const WorkstationView = () => {
  const [params, setPs] = useSearchParams();
  const selected = params.get("selected");

  return (
    <Modal
      size="xl"
      isOpen={selected !== null}
      onClose={() => {
        setPs((ps) => {
          ps.delete("selected");
          return ps;
        });
      }}
    >
      <ModalOverlay bg="blackAlpha.300" backdropFilter="blur(15px)" />
      <ModalContent>
        <ModalHeader>{`${selected}`}</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          Body!
          <Graph
            width={200}
            height={200}
            data={[
              { x: 1, y: 90 },
              { x: 2, y: 12 },
              { x: 3, y: 34 },
              { x: 4, y: 53 },
              { x: 5, y: 98 },
            ]}
          ></Graph>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
};
