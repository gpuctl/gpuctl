import { VStack, Heading, UnorderedList, ListItem } from "@chakra-ui/react";

export const AggregateStats = ({
  data,
}: {
  data: { percent_used: number; total_energy: number };
}) => {
  return (
    <VStack paddingLeft={5} spacing={5} align="left">
      <Heading size="lg">Statistics</Heading>
      <UnorderedList spacing={3}>
        <ListItem>Percent Used: {data.percent_used}%</ListItem>
        <ListItem>Total Energy: {data.total_energy} Joules</ListItem>
      </UnorderedList>
    </VStack>
  );
};
