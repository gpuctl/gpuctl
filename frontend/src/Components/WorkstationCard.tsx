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

export const WorkstationTab = ({ name, gpus }: WorkStationData) => {
  const textCol = useColorModeValue("black", "white");
  return (
    <Center>
      <Box w="100%" rounded={"md"} bg={useColorModeValue("white", "gray.900")}>
        <Heading size="lg" color={textCol}>
          {name}
        </Heading>
        {gpus.map((s, i) => {
          return (
            <Box key={i}>
              <Table variant="striped">
                <Tbody>
                  <Tr>
                    <Td>
                      <Heading size="md">GPU Model:</Heading>
                    </Td>
                    <Td>
                      <Heading size="md">{s.gpu_name}</Heading>
                    </Td>
                  </Tr>
                  <Tr>
                    <Td>
                      <Heading size="md">Free?</Heading>
                    </Td>
                    <Td>
                      <Heading size="md">
                        {s.gpu_util < 5
                          ? `Yes ✅ (${s.gpu_util}% used)`
                          : "No ❌"}
                      </Heading>
                    </Td>
                  </Tr>
                  <Tr>
                    <Td>
                      <Heading size="md">GPU Memory:</Heading>
                    </Td>
                    <Td>
                      <Heading size="md">
                        {(s.memory_total / 1000).toFixed(0)} GB
                      </Heading>
                    </Td>
                  </Tr>
                  <Tr>
                    <Td>
                      <Heading size="md">GPU Temperature:</Heading>
                    </Td>
                    <Td>
                      <Heading size="md">{s.gpu_temp} °C</Heading>
                    </Td>
                  </Tr>
                </Tbody>
              </Table>
            </Box>
          );
        })}
      </Box>
    </Center>
  );
};
