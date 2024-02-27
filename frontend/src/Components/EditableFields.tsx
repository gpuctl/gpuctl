import { Editable, EditableInput, EditablePreview, Td } from "@chakra-ui/react";
import { useModifyInfo, FieldKey } from "../Hooks/Hooks";
import { WorkStationData } from "../Data";
import { GS } from "../Pages/AdminPanel";

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

  // Given we now are not useing GroupSelect, there is a bit of code duplication
  // here, but be careful if you want to refactor it! 'group' (unlike the other
  // fields) is NOT stored in the workstation JSON!
  return (
    <Td>
      {fieldKey === "group" ? (
        <Editable
          defaultValue={group}
          placeholder={placeholder}
          textColor={pickCol(group)}
          onSubmit={(a) => {
            handleSubmit(a);
          }}
        >
          <EditablePreview />
          <EditableInput />
        </Editable>
      ) : (
        // <GroupSelect onChange={handleSubmit}>
        //   <AutoCompleteInput
        //     ref={ref}
        //     placeholder={"Unknown"}
        //     onSelect={() => {
        //       setTimeout(() => {
        //         ref.current?.focus();
        //       }, 1);
        //       setShouldFetch(false);
        //     }}
        //     onBlur={() => {
        //       setShouldFetch(true);
        //     }}
        //     onChange={(a) => setGroupBacking(a.target.value)}
        //     value={groupBacking}
        //   ></AutoCompleteInput>
        // </GroupSelect>
        <Editable
          defaultValue={
            workstation[fieldKey as keyof WorkStationData] as string
          }
          placeholder={placeholder}
          textColor={pickCol(
            workstation[fieldKey as keyof WorkStationData] as string,
          )}
          onSubmit={(a) => {
            handleSubmit(a);
          }}
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
