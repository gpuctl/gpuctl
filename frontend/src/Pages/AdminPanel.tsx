import { API_URL } from "../App";
import { STATS_PATH } from "../Config/Paths";

export const ADMIN_PATH = "/admin";
const ADD_MACHINE_URL = "/add_workstation";

const addMachine = (hostname: string, group: string) => {
  fetch(API_URL + ADMIN_PATH + ADD_MACHINE_URL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ hostname, group }),
  });
  // We should probably await the response to give feedback on whether adding
  // the machine was successful...
};

type ModifyData = {
  group: string | null;
  motherboard: string | null;
  cpu: string | null;
  notes: string | null;
};

const modifyInfo = (hostname: string, mod: ModifyData) => {
  fetch(API_URL + ADMIN_PATH + STATS_PATH + "/modify", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ hostname, ...mod }),
  });
};

export const AdminPanel = () => {
  // TODO
  return <></>;
};
