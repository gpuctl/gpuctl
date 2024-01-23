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

export const WorkstationCardNew = ({ name, gpus }: WorkStationData) => {
  const textCol = useColorModeValue("black", "white");
  return (
    <Center>
      <Box w="100%" rounded={"md"} bg={useColorModeValue("white", "gray.900")}>
        <Heading size="md" color={textCol}>
          {name}
        </Heading>
        {gpus.map((s, i) => {
          return <Box key={i}>{s.gpu_name}</Box>;
        })}
      </Box>
    </Center>
  );
};
