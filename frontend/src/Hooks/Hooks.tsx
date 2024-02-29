import { useToast } from "@chakra-ui/react";
import { STATS_PATH } from "../Config/Paths";
import { useAuth } from "../Providers/AuthProvider";
import { validatedElim } from "../Utils/Utils";

const ADD_MACHINE_URL = "/add_workstation";
const REMOVE_MACHINE_URL = "/rm_workstation";

export const useAddMachine = (callback: () => void) => {
  const toast = useToast();
  const { useAuthFetch } = useAuth();
  const [, addMachineAuth] = useAuthFetch(ADD_MACHINE_URL, (r) => {
    callback();
    validatedElim(r, {
      success: () => {
        toast({
          title: "Add machine successful!",
          description: "",
          status: "success",
          duration: 3000,
          isClosable: true,
        });
      },
      failure: () => {
        toast({
          title: "Add machine failed",
          description:
            "Please check that the hostname is correct and fully qualified",
          status: "error",
          duration: 9000,
          isClosable: true,
        });
      },
    });
  });

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
  const toast = useToast();
  const { useAuthFetch } = useAuth();
  const [, addMachineAuth] = useAuthFetch(REMOVE_MACHINE_URL, (r) => {
    validatedElim(r, {
      success: (resp) => {
        switch (resp.status) {
          case 200: {
            toast({
              title: "Remove machine successful!",
              description: "",
              status: "success",
              duration: 3000,
              isClosable: true,
            });
            break;
          }
          case 512: {
            toast({
              title: "SSH to machine failed",
              description: `Ensure the monitor account details are authorised for that machine`,
              status: "error",
              duration: 9000,
              isClosable: true,
            });
            break;
          }
          default:
            toast({
              title: "Remove machine failed with error ${resp.status}",
              description: `Please report this error to maintainers. ${resp.text}`,
              status: "error",
              duration: 9000,
              isClosable: true,
            });
        }
      },
      failure: () => {},
    });
  });
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
