import * as d3 from "d3";
import { ScaleLinear } from "d3";
import { useMemo } from "react";
import { inlineLog } from "../Utils/Utils";

const AXIS_MARGIN = { x: 20, y: 20 };

export const Graph = ({
  width,
  height,
  data,
}: {
  width: number;
  height: number;
  data: { x: number; y: number }[];
}) => {
  const minX = Math.min(...data.map(({ x }) => x));
  const maxX = Math.max(...data.map(({ x }) => x));
  const minY = Math.min(...data.map(({ y }) => y));
  const maxY = Math.max(...data.map(({ y }) => y));

  const xScale: ScaleLinear<number, number> = d3
    .scaleLinear()
    .domain([minX, maxX])
    .range([0, width]);

  const yScale: ScaleLinear<number, number> = d3
    .scaleLinear()
    .domain([maxY, minY])
    .range([0, height]);

  const lineBuilder = d3
    .line()
    .x(([x]) => xScale(x))
    .y(([, y]) => yScale(y));

  const linePath = lineBuilder(data.map(({ x, y }) => [x, y]));

  console.log(`Attempted to create a graph. ${linePath}`);

  return (
    <svg width={width + AXIS_MARGIN.x * 2} height={height + AXIS_MARGIN.y * 2}>
      <g
        width={width}
        height={height}
        transform={`translate(${AXIS_MARGIN.x}, ${AXIS_MARGIN.y})`}
      >
        <g transform={`translate(0, 0)`} shapeRendering={"geometricPrecision"}>
          <Axis scale={yScale} pixelsPerTick={40} vertical={true} />
        </g>

        <g
          transform={`translate(0, ${height})`}
          shapeRendering={"geometricPrecision"}
        >
          <Axis scale={xScale} pixelsPerTick={40} vertical={false} />
        </g>

        <path
          d={linePath ?? undefined}
          stroke="blue"
          fill="none"
          strokeWidth={3}
        />
      </g>
    </svg>
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
    const diff = Math.abs(range[1] - range[0]);
    const numberOfTicksTarget = Math.round(diff / pixelsPerTick);

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
