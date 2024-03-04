import * as d3 from "d3";
import { ScaleLinear } from "d3";
import { useMemo, useRef } from "react";
import { useDims } from "../Utils/Hooks";
import { Box } from "@chakra-ui/react";

const AXIS_MARGIN = { x: 20, y: 20 };

export const Graph = ({ data }: { data: { x: number; y: number }[][] }) => {
  const minX = Math.min(...data.flatMap((d) => d.map(({ x }) => x)));
  const maxX = Math.max(...data.flatMap((d) => d.map(({ x }) => x)));
  const minY = Math.min(...data.flatMap((d) => d.map(({ y }) => y)));
  const maxY = Math.max(...data.flatMap((d) => d.map(({ y }) => y)));

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

  const linePaths = data.map((d) => lineBuilder(d.map(({ x, y }) => [x, y])));

  return (
    <Box minWidth={200} minHeight={400} ref={ref}>
      <svg width={width} height={height}>
        <g
          width={innerWidth}
          height={innerHeight}
          transform={`translate(${AXIS_MARGIN.x}, ${AXIS_MARGIN.y})`}
        >
          <g
            transform={`translate(0, 0)`}
            shapeRendering={"geometricPrecision"}
          >
            <Axis scale={yScale} pixelsPerTick={40} vertical={true} />
          </g>

          <g
            transform={`translate(0, ${innerHeight})`}
            shapeRendering={"geometricPrecision"}
          >
            <Axis scale={xScale} pixelsPerTick={40} vertical={false} />
          </g>

          {linePaths.map((p, i) => (
            <path
              d={p ?? undefined}
              stroke={["blue", "red", "green", "orange", "purple"][i % 5]}
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
}: {
  scale: ScaleLinear<number, number>;
  pixelsPerTick: number;
  vertical: boolean;
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
            {value}
          </text>
        </g>
      ))}
    </>
  );
};
