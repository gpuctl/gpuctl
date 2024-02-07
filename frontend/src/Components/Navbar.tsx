import {
  Button,
  Heading,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverCloseButton,
  PopoverContent,
  PopoverHeader,
  PopoverTrigger,
  Spacer,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
} from "@chakra-ui/react";
import { ReactNode, useState } from "react";
import { useNavigate } from "react-router-dom";
import { STATS_PATH } from "../Config/Paths";
import { VIEW_PAGE_INDEX, ViewPage } from "../App";
import { SignIn } from "./SignIn";

const URLS = [
  `${STATS_PATH}/card_view`,
  `${STATS_PATH}/table_view`,
  "/admin_panel",
];

export const Navbar = ({
  children,
  initial,
}: {
  children: ReactNode[];
  initial: ViewPage;
}) => {
  const nav = useNavigate();

  const [signingIn, setSigningIn] = useState<boolean>(false);

  return (
    <Tabs
      variant="soft-rounded"
      onChange={(i) => nav(URLS[i])}
      defaultIndex={VIEW_PAGE_INDEX[initial]}
    >
      <TabList>
        <Tab>Card View</Tab>
        <Tab>Table View</Tab>
        <Tab isDisabled>Admin Panel</Tab>
        <Spacer />
        <Popover>
          <PopoverTrigger>
            <Button>Admin Sign In</Button>
          </PopoverTrigger>
          <PopoverContent w="100%">
            <PopoverArrow />
            <PopoverCloseButton />
            <SignIn />
          </PopoverContent>
        </Popover>
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
