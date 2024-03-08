import * as d3 from "d3";
import { ScaleLinear } from "d3";
import { useMemo, useRef } from "react";
import { useDims } from "../Utils/Hooks";
import { Box } from "@chakra-ui/react";
import { catNulls, chunks, mapNotNulls } from "../Utils/Utils";
import { GRAPH_COLS } from "./WorkstationView";

const AXIS_MARGIN = { x: 20, y: 20 };

export const Graph = ({
  data,
  xlabel,
  xdivisor,
  ydivisor,
  maxPoints,
}: {
  data: ({ off: number; line: { x: number; y: number }[] } | null)[];
  xlabel: string;
  xdivisor?: number;
  ydivisor?: number;
  maxPoints: number;
}) => {
  const minX = Math.min(
    ...catNulls(data.flatMap((d) => d?.line?.map(({ x }) => x))),
  );
  const maxX = Math.max(
    ...catNulls(data.flatMap((d) => d?.line?.map(({ x }) => x))),
  );
  const minY = Math.min(
    ...catNulls(data.flatMap((d) => d?.line?.map(({ y }) => y))),
  );
  const maxY = Math.max(
    ...catNulls(data.flatMap((d) => d?.line?.map(({ y }) => y))),
  );

  const ref = useRef<HTMLHeadingElement>(null);

  const { w: width, h: height } = useDims(ref);

  const innerWidth = width - AXIS_MARGIN.x * 2;
  const innerHeight = height - AXIS_MARGIN.y * 2;

  const xScale: ScaleLinear<number, number> = d3
    .scaleLinear()
    .domain([minX, maxX])
    .range([0, innerWidth]);

  const yScale: ScaleLinear<number, number> = d3
    .scaleLinear()
    .domain([maxY, minY])
    .range([0, innerHeight]);

  const lineBuilder = d3
    .line()
    .x(([x]) => xScale(x))
    .y(([, y]) => yScale(y));

  const downsampled = mapNotNulls(data, ({ off, line }, i) => {
    const lineMinX = Math.min(...line.map(({ x }) => x));
    const lineMaxX = Math.max(...line.map(({ x }) => x));
    const lineFrac = (lineMaxX - lineMinX) / (maxX - minX);

    const chunkSize = Math.ceil(line.length / (maxPoints * lineFrac));
    return chunks(line, chunkSize, off).map((c) => {
      const { x, y } = c.reduce(({ x: x1, y: y1 }, { x: x2, y: y2 }) => ({
        x: x1,
        y: y1 + y2,
      }));
      return { x, y: y / c.length };
    });
  });

  const linePaths = mapNotNulls(downsampled, (d) =>
    lineBuilder(d.map(({ x, y }) => [x, y])),
  );

  return (
    <Box minWidth={200} minHeight={400} ref={ref}>
      <svg width={width} height={height}>
        <g
          width={innerWidth}
          height={innerHeight}
          transform={`translate(${AXIS_MARGIN.x * 2 - AXIS_MARGIN.x / 4}, ${AXIS_MARGIN.y / 4})`}
        >
          <g
            transform={`translate(0, 0)`}
            shapeRendering={"geometricPrecision"}
          >
            <Axis scale={yScale} pixelsPerTick={40} vertical={true} divisor={ydivisor ?? 1} />
          </g>

          <g
            transform={`translate(${innerWidth / 2}, ${height - 5})`}
            shapeRendering={"geometricPrecision"}
          >
            <text
              style={{
                fontSize: "15px",
                textAnchor: "middle",
              }}
            >
              {xlabel}
            </text>
          </g>

          <g
            transform={`translate(0, ${innerHeight})`}
            shapeRendering={"geometricPrecision"}
          >
            <Axis scale={xScale} pixelsPerTick={40} vertical={false} divisor={xdivisor ?? 1} />
          </g>

          {linePaths.map((p, i) => (
            <path
              d={p ?? undefined}
              stroke={GRAPH_COLS[i % GRAPH_COLS.length]}
              fill="none"
              strokeWidth={3}
            />
          ))}
        </g>
      </svg>
    </Box>
  );
};

// From https://www.react-graph-gallery.com/scatter-plot

const TICK_LENGTH = 6;

export const Axis = ({
  scale,
  pixelsPerTick,
  vertical,
  divisor,
}: {
  scale: ScaleLinear<number, number>;
  pixelsPerTick: number;
  vertical: boolean;
  divisor: number;
}) => {
  const range = scale.range();

  const ticks = useMemo(() => {
    // const diff = Math.abs(range[1] - range[0]);
    // const numberOfTicksTarget = Math.round(diff / pixelsPerTick);
    const numberOfTicksTarget = 10;

    return scale.ticks(numberOfTicksTarget).map((value) => ({
      value,
      offset: scale(value),
    }));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [scale]);

  return (
    <>
      {/* Main axis line */}
      <path
        d={(vertical
          ? ["M", 0, range[0], "L", 0, range[1]]
          : ["M", range[0], 0, "L", range[1], 0]
        ).join(" ")}
        fill="none"
        stroke="currentColor"
      />

      {/* Ticks and labels */}
      {ticks.map(({ value, offset }) => (
        <g
          key={value}
          transform={
            vertical ? `translate(0, ${offset})` : `translate(${offset}, 0)`
          }
        >
          {vertical ? (
            <line
              x1={0}
              x2={-TICK_LENGTH}
              stroke="currentColor"
              strokeWidth={1}
            />
          ) : (
            <line
              y1={0}
              y2={TICK_LENGTH}
              stroke="currentColor"
              strokeWidth={1}
            />
          )}
          <text
            key={value}
            style={{
              fontSize: "10px",
              textAnchor: "middle",
              transform: vertical
                ? "translate(-15px, 0)"
                : "translate(0, 15px)",
            }}
          >
            {value / divisor}
          </text>
        </g>
      ))}
    </>
  );
};
