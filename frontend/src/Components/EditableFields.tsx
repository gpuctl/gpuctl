import {
  Editable,
  EditableInput,
  EditablePreview,
  HStack,
  Flex,
  IconButton,
  Td,
} from "@chakra-ui/react";
import { EditIcon } from "@chakra-ui/icons";
import { useModifyInfo, FieldKey } from "../Hooks/Hooks";
import { WorkStationData } from "../Data";

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
  isEven,
}: EditableFieldProps & {
  isEven: boolean;
}) => {
  const pickCol = (value: string) => (value ? "gray.600" : "gray.300");
  const modifyInfo = useModifyInfo();

  const handleSubmit = (newValue: string) => {
    modifyInfo(workstation.name, {
      group: fieldKey === "group" ? newValue : null,
      motherboard: fieldKey === "motherboard" ? newValue : null,
      cpu: fieldKey === "cpu" ? newValue : null,
      notes: fieldKey === "notes" ? newValue : null,
      owner: fieldKey === "owner" ? newValue : null,
    });
  };

  const EditableButton = ({
    isEditing,
    onEdit,
  }: {
    isEditing: boolean;
    onEdit: () => void;
  }) => {
    return isEditing ? (
      <> </>
    ) : (
      <IconButton
        size="sm"
        icon={<EditIcon />}
        aria-label="edit"
        onClick={onEdit}
        bgColor={isEven ? "gray.100" : "white"}
      />
    );
  };

  // Given we now are not using GroupSelect, there is a bit of code duplication
  // here, but be careful if you want to refactor it! 'group' (unlike the other
  // fields) is NOT stored in the workstation JSON!

  // Auto-complete code removed in
  // https://github.com/gpuctl/gpuctl/pull/220
  // Go back there if you want to re-implement it
  return (
    <Td pr="0">
      <Flex min-width="5rem">
        {fieldKey === "group" ? (
          <Editable
            defaultValue={group}
            placeholder={placeholder}
            textColor={pickCol(group)}
            onSubmit={(a) => {
              handleSubmit(a);
            }}
          >
            {(props) => (
              <HStack>
                <EditablePreview isTruncated />
                <EditableInput />
                <EditableButton {...props} />
              </HStack>
            )}
          </Editable>
        ) : (
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
            {(props) => (
              <HStack>
                <EditablePreview isTruncated />
                <EditableInput />
                <EditableButton {...props} />
              </HStack>
            )}
          </Editable>
        )}
      </Flex>
    </Td>
  );
};
