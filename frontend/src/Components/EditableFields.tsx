import { Editable, EditableInput, EditablePreview, Td } from "@chakra-ui/react";
import { useModifyInfo } from "../Hooks/Hooks";
import { WorkStationData } from "../Data";

type EditableFieldProps = {
  workstation: WorkStationData;
  fieldKey: "cpu" | "motherboard" | "notes" | "group";
  placeholder: string;
};

export const EditableField = ({
  workstation,
  fieldKey,
  placeholder,
}: EditableFieldProps) => {
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

  return (
    <Td>
      <Editable
        defaultValue={workstation[fieldKey as keyof WorkStationData] as string}
        placeholder={placeholder}
        textColor={pickCol(
          workstation[fieldKey as keyof WorkStationData] as string,
        )}
        onSubmit={handleSubmit}
      >
        <EditablePreview />
        <EditableInput />
      </Editable>
    </Td>
  );
};
