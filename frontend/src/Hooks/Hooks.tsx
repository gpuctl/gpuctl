import { useToast } from "@chakra-ui/react";
import { STATS_PATH } from "../Config/Paths";
import { GPUStats } from "../Data";
import { useAuth } from "../Providers/AuthProvider";
import {
  fire,
  success,
  Validated,
  validatedElim,
  Validation,
} from "../Utils/Utils";
import { API_URL } from "../App";
import { useJarJar } from "../Utils/Hooks";
import { useInterval } from "@chakra-ui/react";

const ADD_MACHINE_URL = "/add_workstation";
const REMOVE_MACHINE_URL = "/rm_workstation";

const GET_ALL_FILES_URL = "/list_files";
const UPLOAD_FILE_URL = "/attach_file";
const REMOVE_FILE_URL = "/remove_file";
const GET_SPECIFIC_FILE_URL = "/get_file";

export const useGetAllFiles = (
  hostname: string,
  callback: (r: Validated<Response>) => void,
) => {
  const { useAuthFetch } = useAuth();
  const [, getAllFilesAuth] = useAuthFetch(
    GET_ALL_FILES_URL + "?hostname=" + hostname,
    callback,
  );

  return () => getAllFilesAuth({ method: "GET" });
};

const encodeFile = (file: Uint8Array) => {
  let binaryString = "";
  const len = file.byteLength;
  for (let i = 0; i < len; i++) {
    binaryString += String.fromCharCode(file[i]);
  }
  return btoa(binaryString);
};

export const useUploadFile = (callback: (r: Validated<Response>) => void) => {
  const { useAuthFetch } = useAuth();
  const [, uploadFileAuth] = useAuthFetch(UPLOAD_FILE_URL, callback);

  return (hostname: string, mime: string, filename: string, file: Uint8Array) =>
    uploadFileAuth({
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        hostname,
        mime,
        filename,
        file_enc: encodeFile(file),
      }),
    });
};

export const useRemoveFile = () => {
  const { useAuthFetch } = useAuth();
  const [, removeFileAuth] = useAuthFetch(REMOVE_FILE_URL);
  return (hostname: string, filename: string) =>
    removeFileAuth({
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ hostname, filename }),
    });
};

export const useGetSpecificFile = (
  hostname: string,
  filename: string,
  callback: (r: Validated<Response>) => void,
) => {
  const { useAuthFetch } = useAuth();
  const [, getSpecificFileAuth] = useAuthFetch(
    GET_SPECIFIC_FILE_URL + "?hostname=" + hostname + "&file=" + filename,
    callback,
  );
  return () =>
    getSpecificFileAuth({
      method: "GET",
    });
};

export const useAddMachine = (callback: () => void) => {
  const toast = useToast();
  const { useAuthFetch } = useAuth();
  const [, addMachineAuth] = useAuthFetch(ADD_MACHINE_URL, (r) => {
    callback();
    validatedElim(r, {
      success: (resp) => {
        if (resp.status === 200)
          toast({
            title: "Add machine successful!",
            description: "",
            status: "success",
            duration: 3000,
            isClosable: true,
          });
        else
          toast({
            title: "Add machine failed",
            description:
              "Please check that the hostname is correct and fully qualified",
            status: "error",
            duration: 9000,
            isClosable: true,
          });
      },
      failure: () => {},
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
            fire(async () => {
              const msg = await resp.text();
              toast({
                title: `Remove machine failed with error ${resp.status}`,
                description: `Please report this error to maintainers. ${msg}`,
                status: "error",
                duration: 9000,
                isClosable: true,
              });
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
  owner: string | null;
};

export const useModifyInfo = () => {
  const toast = useToast();
  const { useAuthFetch } = useAuth();
  const [, modifyAuth] = useAuthFetch(STATS_PATH + "/modify", (r) => {
    validatedElim(r, {
      success: (r) => {
        if (!r.ok) {
          toast({
            title: `Updating workstation information failed with error ${r.status}`,
            description: `Please inform the maintainers`,
            status: "error",
            duration: 9000,
            isClosable: true,
          });
        }
      },
      failure: () => {},
    });
  });
  return (hostname: string, modification: ModifyData) =>
    modifyAuth({
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ hostname, ...modification }),
    });
};

export type FieldKey = "cpu" | "motherboard" | "notes" | "group" | "owner";

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
    owner: fieldKey === "owner" ? newValue : null,
  });
};

export type HistorySample = {
  timestamp: number;
  sample: GPUStats;
};

const GRAPH_REFRESH_INTERVAL = 5000;

export const useHistoryStats = (
  hostname: string,
): Validation<HistorySample[][]> => {
  const [stats, updateStats] = useJarJar<HistorySample[][]>(async () =>
    success(
      await (
        await fetch(API_URL + `/stats/historical?hostname=${hostname}`)
      ).json(),
    ),
  );

  useInterval(() => {
    updateStats();
  }, GRAPH_REFRESH_INTERVAL);

  return stats;
};

export type AggregateData = { percent_used: number; total_energy: number };

export const useAggregateStats = (): Validation<AggregateData> => {
  const [agg, updateAgg] = useJarJar<AggregateData>(async () =>
    success(await (await fetch(API_URL + `/stats/aggregate`)).json()),
  );

  useInterval(() => {
    updateAgg();
  }, GRAPH_REFRESH_INTERVAL);

  return agg;
};
