import {
  Box,
  Center,
  Heading,
  Text,
  Table,
  useColorModeValue,
  Thead,
  Tbody,
  Td,
  Tr,
} from "@chakra-ui/react";
import { WorkStationData } from "../App";
import { useRef } from "react";

export const WorkstationCardNew = ({ name, gpus }: WorkStationData) => {
  const textCol = useColorModeValue("black", "white");
  const ref = useRef<HTMLHeadingElement>(null);
  return (
    <Center>
      <Box w="100%" rounded={"md"} bg={useColorModeValue("white", "gray.900")}>
        <Heading size="md" color={textCol}>
          {name}
        </Heading>
        {gpus.map((s, i) => {
          return (
            <Box key={i}>
              <p>{`${s.gpu_name} (${(s.memory_total / 1000).toFixed(
                0
              )} GB)`}</p>
              <p>{true ? "🟢 Free" : "🔴 In-use"}</p>
            </Box>
          );
        })}
      </Box>
    </Center>
  );
};
