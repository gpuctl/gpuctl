import {
  Button,
  Heading,
  Spacer,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
} from "@chakra-ui/react";
import { ReactNode } from "react";
import { JsxElement } from "typescript";

export const Navbar = ({ children }: { children: ReactNode[] }) => {
  return (
    <Tabs variant="soft-rounded">
      <TabList>
        <Tab>Card View</Tab>
        <Tab>Table View</Tab>
        <Spacer />
        <Button mr={5}> Admin Sign In </Button>
      </TabList>
      <Heading size="2xl">Welcome to the GPU Control Room!</Heading>
      <TabPanels>
        {children.map((c, i) => (
          <TabPanel key={i}>{c}</TabPanel>
        ))}
      </TabPanels>
    </Tabs>
  );
};
