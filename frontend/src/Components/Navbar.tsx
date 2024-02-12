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
} from "@chakra-ui/react";
import { ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { STATS_PATH } from "../Config/Paths";
import { AuthToken, VIEW_PAGE_INDEX, ViewPage } from "../App";
import { SignIn } from "./SignIn";
import { Validated, isSuccess } from "../Utils/Utils";

const URLS = [
  `${STATS_PATH}/card_view`,
  `${STATS_PATH}/table_view`,
  `${STATS_PATH}/admin_view`,
];

export const Navbar = ({
  children,
  signIn,
  signOut,
  authToken,
  initial,
}: {
  children: ReactNode[];
  signIn: (tok: AuthToken) => void;
  signOut: () => void;
  authToken: Validated<AuthToken>;
  initial: ViewPage;
}) => {
  const nav = useNavigate();

  return (
    <Tabs
      variant="soft-rounded"
      onChange={(i) => nav(URLS[i])}
      defaultIndex={VIEW_PAGE_INDEX[initial]}
    >
      <TabList>
        <Tab>Card View</Tab>
        <Tab>Table View</Tab>
        <Tab isDisabled={!isSuccess(authToken)}>Admin Panel</Tab>
        <Spacer />
        {isSuccess(authToken) ? (
          <Button onClick={signOut}>Sign Out</Button>
        ) : (
          <Popover>
            <PopoverTrigger>
              <Button>Admin Sign In</Button>
            </PopoverTrigger>
            <PopoverContent w="100%">
              <PopoverArrow />
              <PopoverCloseButton />
              <SignIn signIn={signIn} />
            </PopoverContent>
          </Popover>
        )}
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
