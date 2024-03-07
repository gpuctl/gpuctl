import {
  IconButton,
  Popover,
  PopoverArrow,
  PopoverContent,
  PopoverTrigger,
  PopoverCloseButton,
  PopoverHeader,
  PopoverBody,
  Textarea,
} from "@chakra-ui/react";
import { useModifyInfo } from "../Hooks/Hooks";
import { EditIcon } from "@chakra-ui/icons";
import { useState } from "react";

export const NotesPopout = ({
  wname,
  notes,
}: {
  wname: string;
  notes: string;
}) => {
  const modifyInfo = useModifyInfo();
  const [newNotes, setNewNotes] = useState(notes);

  const handleSubmit = () => {
    setNewNotes(newNotes);
    modifyInfo(wname, {
      group: null,
      motherboard: null,
      cpu: null,
      notes: newNotes,
      owner: null,
    });
  };

  return (
    <Popover>
      <PopoverTrigger>
        <IconButton size="sm" icon={<EditIcon />} aria-label="edit" />
      </PopoverTrigger>
      <PopoverContent>
        <PopoverArrow />
        <PopoverCloseButton />
        <PopoverHeader>Notes: {wname}</PopoverHeader>
        <PopoverBody>
          <Textarea
            value={newNotes}
            onChange={(s) => {
              setNewNotes(s.target.value);
            }}
            onBlur={() => {
              handleSubmit();
            }}
          />
        </PopoverBody>
      </PopoverContent>
    </Popover>
  );
};
