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

enum ModCol {
  MOTHERBOARD = "motherboard",
  CPU = "cpu",
  NOTES = "notes",
}

const modifyInfo = (hostname: string, column: string, data: string) => {
  fetch(API_URL + ADMIN_PATH + STATS_PATH + "/modify", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ hostname, column, data }),
  });
};

export const AdminPanel = () => {
  // TODO
};
