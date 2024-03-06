import {
  Button,
  Heading,
  Popover,
  PopoverArrow,
  PopoverCloseButton,
  PopoverContent,
  PopoverTrigger,
  Spacer,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
} from "@chakra-ui/react";
import { useState, ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { STATS_PATH } from "../Config/Paths";
import { VIEW_PAGE_INDEX, ViewPage } from "../App";
import { SignIn } from "./SignIn";
import { useAuth } from "../Providers/AuthProvider";

const URLS = [
  `${STATS_PATH}/card_view`,
  `${STATS_PATH}/table_view`,
  `${STATS_PATH}/admin_view`,
];

export const Navbar = ({
  children,
  initial,
}: {
  children: ReactNode[];
  initial: ViewPage;
}) => {
  const { isSignedIn, logout } = useAuth();
  const nav = useNavigate();
  const [tabs] = useState({ index: VIEW_PAGE_INDEX[initial] });
  const admin_panel = 2;

  return (
    <Tabs
      align="center"
      onChange={(i) => {
        nav(URLS[i]);
        tabs.index = i;
      }}
      defaultIndex={VIEW_PAGE_INDEX[initial]}
      index={tabs.index}
    >
      <Heading size="xl">GPU Control Room</Heading>
      <TabList>
        <Tab>Card View</Tab>
        <Tab>Table View</Tab>
        {isSignedIn() ? (
          <Tab>Admin Panel</Tab>
        ) : (
          <Popover>
            <PopoverTrigger>
              <Button variant="ghost" _hover={{ background: "white" }}>
                <Text fontWeight="normal">Admin Panel</Text>
              </Button>
            </PopoverTrigger>
            <PopoverContent w="100%">
              <PopoverArrow />
              <PopoverCloseButton />
              <SignIn
                panelCallback={() => {
                  tabs.index = admin_panel;
                }}
              />
            </PopoverContent>
          </Popover>
        )}
      </TabList>
      <TabPanels>
        {children.map((c, i) => (
          <Text align="left">
            <TabPanel key={i}>{c}</TabPanel>
          </Text>
        ))}
      </TabPanels>
    </Tabs>
  );
};
