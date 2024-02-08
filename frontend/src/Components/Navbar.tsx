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
import { AUTH_TOKEN_ITEM, AuthToken, VIEW_PAGE_INDEX, ViewPage } from "../App";
import { SignIn } from "./SignIn";
import { Validated, failure, isSuccess } from "../Utils/Utils";

const URLS = [
  `${STATS_PATH}/card_view`,
  `${STATS_PATH}/table_view`,
  `${STATS_PATH}/admin_view`,
];

export const Navbar = ({
  children,
  setAuth,
  authToken,
  initial,
}: {
  children: ReactNode[];
  setAuth: (tok: Validated<AuthToken>) => void;
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
          <Button
            onClick={() => {
              setAuth(failure(Error("Signed Out")));
              localStorage.removeItem(AUTH_TOKEN_ITEM);
              window.location.reload();
            }}
          >
            Sign Out
          </Button>
        ) : (
          <Popover>
            <PopoverTrigger>
              <Button>Admin Sign In</Button>
            </PopoverTrigger>
            <PopoverContent w="100%">
              <PopoverArrow />
              <PopoverCloseButton />
              <SignIn setAuth={setAuth} />
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
