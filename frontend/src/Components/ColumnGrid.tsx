import { ReactNode, useRef } from "react";
import { useContainerDim } from "../Utils/Hooks";
import { HStack, VStack } from "@chakra-ui/react";

export const ColumnGrid = ({
  minWidth,
  children,
}: {
  minWidth: number;
  children: ReactNode[];
}) => {
  const { w: width } = useContainerDim(useRef<HTMLHeadingElement>());
  const numCols = Math.floor(width / minWidth);
  const grouped: ReactNode[][] = Array.from(Array(5)).map(() => []);
  let colNum = 0;
  children.forEach((c) => {
    grouped[colNum].push(c);
    colNum = (colNum + 1) % numCols;
  });
  return (
    <HStack>
      {grouped.map((cs) => (
        <VStack>{cs}</VStack>
      ))}
    </HStack>
  );
};
