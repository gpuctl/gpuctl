import {
  VStack,
  Heading,
  UnorderedList,
  ListItem,
  Text,
} from "@chakra-ui/react";
import { useAggregateStats } from "../Hooks/Hooks";
import { validationElim } from "../Utils/Utils";

export const AggregateStats = () => {
  const stats = useAggregateStats();

  return (
    <VStack paddingLeft={5} spacing={5} align="left">
      <Heading size="lg">Statistics</Heading>
      {validationElim(stats, {
        success: (s) => (
          <UnorderedList spacing={3}>
            <ListItem>
              Total Energy: {Math.round(s.total_energy / (60 * 60))} Watt Hours
            </ListItem>
          </UnorderedList>
        ),
        failure: () => (
          <Text>Failed to fetch aggregate statistics! Retrying...</Text>
        ),
        loading: () => <Text>Fetching aggregate statistics...</Text>,
      })}
    </VStack>
  );
};
