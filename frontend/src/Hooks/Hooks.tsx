import { STATS_PATH } from "../Config/Paths";
import { useAuth } from "../Providers/AuthProvider";

const ADD_MACHINE_URL = "/add_workstation";
const REMOVE_MACHINE_URL = "/rm_workstation";

export const useAddMachine = () => {
  const { useAuthFetch } = useAuth();
  const [, addMachineAuth] = useAuthFetch(ADD_MACHINE_URL);
  return (hostname: string, group: string) =>
    addMachineAuth({
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ hostname, group }),
    });
};

export const useRemoveMachine = () => {
  const { useAuthFetch } = useAuth();
  const [, addMachineAuth] = useAuthFetch(REMOVE_MACHINE_URL);
  return (hostname: string) =>
    addMachineAuth({
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ hostname }),
    });
};

export type ModifyData = {
  group: string | null;
  motherboard: string | null;
  cpu: string | null;
  notes: string | null;
};

export const useModifyInfo = () => {
  const { useAuthFetch } = useAuth();
  const [, addMachineAuth] = useAuthFetch(STATS_PATH + "/modify");
  return (hostname: string, modification: ModifyData) =>
    addMachineAuth({
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ hostname, ...modification }),
    });
};

export type FieldKey = "cpu" | "motherboard" | "notes" | "group";

export const useHandleSubmit = (
  name: string,
  fieldKey: FieldKey,
  newValue: string,
) => {
  const modifyInfo = useModifyInfo();

  modifyInfo(name, {
    group: fieldKey === "group" ? newValue : null,
    motherboard: fieldKey === "motherboard" ? newValue : null,
    cpu: fieldKey === "cpu" ? newValue : null,
    notes: fieldKey === "notes" ? newValue : null,
  });
};
