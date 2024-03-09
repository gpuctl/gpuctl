import * as d3 from "d3";
import { ScaleLinear, ScaleTime } from "d3";
import { useRef } from "react";
import { useDims } from "../Utils/Hooks";
import { Box } from "@chakra-ui/react";
import { catNulls, chunks, mapNotNulls } from "../Utils/Utils";
import { GRAPH_COLS } from "./WorkstationView";

import {
  utcDay,
  utcFormat,
  utcHour,
  utcMinute,
  utcMonth,
  utcSecond,
  utcWeek,
  utcYear,
} from "d3";

const AXIS_MARGIN = { x: 20, y: 20 };

export const Graph = ({
  data,
  xlabel,
  maxPoints,
}: {
  data: ({ off: number; line: { x: number; y: number }[] } | null)[];
  xlabel: string;
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

  const xScale: ScaleTime<number, number> = d3
    .scaleTime()
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
            <Axis
              min={0}
              max={innerHeight}
              ticks={yScale
                .ticks(10)
                .map((value) => ({ value, offset: yScale(value) }))}
              fmt={(v) => v.toString()}
              pixelsPerTick={40}
              vertical={true}
            />
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
            <Axis
              min={0}
              max={innerWidth}
              ticks={xScale
                .ticks(10)
                .map((value) => ({ value, offset: xScale(value) }))}
              fmt={multiFormat}
              pixelsPerTick={40}
              vertical={false}
            />
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

const formatMillisecond = utcFormat("%S.%L"),
  formatSecond = utcFormat("%H:%M:%S"),
  formatMinute = utcFormat("%H:%M"),
  formatHour = utcFormat("%H:%M"),
  formatDay = utcFormat("%a %d"),
  formatWeek = utcFormat("%b %d"),
  formatMonth = utcFormat("%B"),
  formatYear = utcFormat("%Y");

const multiFormat = (date: Date) =>
  (utcSecond(date) < date
    ? formatMillisecond
    : utcMinute(date) < date
      ? formatSecond
      : utcHour(date) < date
        ? formatMinute
        : utcDay(date) < date
          ? formatHour
          : utcMonth(date) < date
            ? utcWeek(date) < date
              ? formatDay
              : formatWeek
            : utcYear(date) < date
              ? formatMonth
              : formatYear)(date);

export const Axis = <T,>({
  min,
  max,
  pixelsPerTick,
  vertical,
  ticks,
  fmt,
}: {
  min: number;
  max: number;
  pixelsPerTick: number;
  vertical: boolean;
  ticks: { value: T; offset: number }[];
  fmt: (t: T) => string;
}) => {
  return (
    <>
      {/* Main axis line */}
      <path
        d={(vertical
          ? ["M", 0, min, "L", 0, max]
          : ["M", min, 0, "L", max, 0]
        ).join(" ")}
        fill="none"
        stroke="currentColor"
      />

      {/* Ticks and labels */}
      {ticks.map(({ value, offset }, i) => (
        <g
          key={i}
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
            key={i}
            style={{
              fontSize: "10px",
              textAnchor: "middle",
              transform: vertical
                ? "translate(-15px, 0)"
                : "translate(0, 15px)",
            }}
          >
            {typeof value === "number" ? value : fmt(value)}
          </text>
        </g>
      ))}
    </>
  );
};
