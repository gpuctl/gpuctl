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
        <ListItem>Average Usage: {data.percent_used}%</ListItem>
        <ListItem>Uptime: {data.total_energy}</ListItem>
      </UnorderedList>
    </VStack>
  );
};
