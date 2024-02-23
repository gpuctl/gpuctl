import { Editable, EditableInput, EditablePreview, Td } from "@chakra-ui/react";
import { useModifyInfo, FieldKey } from "../Hooks/Hooks";
import { WorkStationData } from "../Data";
import { AutoCompleteInput } from "@choc-ui/chakra-autocomplete";
import { GS } from "../Pages/AdminPanel";
import { useState } from "react";
import { useStats } from "../Providers/FetchProvider";

type EditableFieldProps = {
  group: string;
  workstation: WorkStationData;
  fieldKey: FieldKey;
  placeholder: string;
};

export const EditableField = ({
  group,
  workstation,
  fieldKey,
  placeholder,
  GroupSelect,
}: EditableFieldProps & {
  GroupSelect: GS;
}) => {
  const pickCol = (value: string) => (value ? "gray.600" : "gray.300");
  const modifyInfo = useModifyInfo();

  const handleSubmit = (newValue: string) => {
    modifyInfo(workstation.name, {
      group: fieldKey === "group" ? newValue : null,
      motherboard: fieldKey === "motherboard" ? newValue : null,
      cpu: fieldKey === "cpu" ? newValue : null,
      notes: fieldKey === "notes" ? newValue : null,
    });
  };

  const { setShouldFetch } = useStats();

  const [groupBacking, setGroupBacking] = useState(group);

  return (
    <Td>
      {fieldKey === "group" ? (
        <GroupSelect onChange={handleSubmit}>
          <AutoCompleteInput
            placeholder={"Unknown"}
            onSelect={() => {
              setShouldFetch(false);
            }}
            onBlur={() => {
              setShouldFetch(true);
            }}
            onChange={(a) => setGroupBacking(a.target.value)}
            value={groupBacking}
          ></AutoCompleteInput>
        </GroupSelect>
      ) : (
        <Editable
          defaultValue={
            workstation[fieldKey as keyof WorkStationData] as string
          }
          placeholder={placeholder}
          textColor={pickCol(
            workstation[fieldKey as keyof WorkStationData] as string,
          )}
          // onCancel={(a) => {}}
          onSubmit={(a) => {
            handleSubmit(a);
          }}
          // onEdit={() => setShouldFetch(false)}
        >
          <EditablePreview />
          <EditableInput />
        </Editable>
      )}
    </Td>
  );
};

export const EditableFieldGroupSelect = ({
  workstation,
  fieldKey,
  placeholder,
}: EditableFieldProps) => {};
