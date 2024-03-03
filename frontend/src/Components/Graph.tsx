import * as d3 from "d3";
import { ScaleLinear } from "d3";

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
    .domain([minY, maxY])
    .range([0, height]);

  const lineBuilder = d3
    .line()
    .x(([x]) => xScale(x))
    .y(([, y]) => yScale(y));

  const linePath = lineBuilder(data.map(({ x, y }) => [x, y]));

  console.log(`Attempted to create a graph. ${linePath}`);

  return (
    <svg width={width} height={height}>
      <g>
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
