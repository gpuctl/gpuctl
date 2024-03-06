import { VStack, Heading, UnorderedList, ListItem } from "@chakra-ui/react";

export const AggregateStats = ({
  data,
}: {
  data: { average_use: number; uptime: string };
}) => {
  return (
    <VStack paddingLeft={5} spacing={5} align="left">
      <Heading size="lg">Statistics</Heading>
      <UnorderedList spacing={3}>
        <ListItem>Average Usage: {data.average_use}%</ListItem>
        <ListItem>Uptime: {data.uptime}</ListItem>
      </UnorderedList>
    </VStack>
  );
};
