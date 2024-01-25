import { ReactNode, useRef } from "react";
import { useDims } from "../Utils/Hooks";
import { Center, HStack, VStack } from "@chakra-ui/react";
import { makeArr } from "../Utils/Utils";

export const ColumnGrid = ({
  minChildWidth,
  hMinSpacing,
  vSpacing,
  children,
}: {
  minChildWidth: number;
  hMinSpacing: number;
  vSpacing: number;
  children: ReactNode[];
}) => {
  const ref = useRef<HTMLHeadingElement>(null);
  const { w: width } = useDims(ref);
  const numCols = Math.min(
    Math.max(
      1,
      Math.floor((width - hMinSpacing) / (minChildWidth + hMinSpacing))
    ),
    children.length
  );
  const tempSpace = (width - numCols * minChildWidth) / (numCols + 1);
  const hspacing = numCols === 1 ? 0 : tempSpace;

  const grouped: ReactNode[][] = makeArr(numCols, () => []);
  let colNum = 0;
  children.forEach((c) => {
    grouped[colNum].push(c);
    colNum = (colNum + 1) % numCols;
  });
  return (
    <Center ref={ref}>
      <HStack align="top" spacing={`${hspacing}px`}>
        {grouped.map((cs, i) => (
          <VStack key={i} spacing={vSpacing}>
            {cs}
          </VStack>
        ))}
      </HStack>
    </Center>
  );
};
