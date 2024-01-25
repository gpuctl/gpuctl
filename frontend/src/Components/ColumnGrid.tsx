import { ReactNode, useRef } from "react";
import { useContainerDim } from "../Utils/Hooks";
import { Center, HStack, Spacer, VStack } from "@chakra-ui/react";
import { makeArr } from "../Utils/Utils";

export const ColumnGrid = ({
  minChildWidth,
  spacing,
  children,
}: {
  minChildWidth: number;
  spacing: number;
  children: ReactNode[];
}) => {
  const ref = useRef<HTMLHeadingElement>(null);
  const { w: width } = useContainerDim(ref);
  const numCols = Math.min(
    Math.max(1, Math.floor(width / minChildWidth)),
    children.length
  );
  const tempSpace = (width + 100 - numCols * minChildWidth) / (numCols + 1);
  const spacingaa = numCols === 1 ? 0 : tempSpace;
  console.log(`Huh ${width} ${minChildWidth}`);

  const grouped: ReactNode[][] = makeArr(() => [], numCols);
  let colNum = 0;
  children.forEach((c) => {
    grouped[colNum].push(c);
    colNum = (colNum + 1) % numCols;
  });
  return (
    <Center ref={ref}>
      <HStack align="top" spacing={spacingaa}>
        {grouped.map((cs) => (
          <VStack>{cs}</VStack>
        ))}
      </HStack>
    </Center>
  );
};
