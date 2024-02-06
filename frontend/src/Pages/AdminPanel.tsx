import { API_URL } from "../App";

const ADD_MACHINE_URL = "/add_workstation";

const addMachine = (hostname: string) => {
  fetch(API_URL + ADD_MACHINE_URL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ hostname }),
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
  fetch(API_URL + ADD_MACHINE_URL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ hostname, column, data }),
  });
};

export const AdminPanel = () => {};
