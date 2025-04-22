import * as d3 from "d3";
import { IPieChart } from "@/context/interfaces/chart";
import { useEffect, useRef } from "react";

export const Pie = ({data}: Record<string, IPieChart[]>) => {
  const pieChartRef = useRef(null)
  const renderPieChart = () => {
    const width = 180;
    const height = Math.min(width, 180);

    const color = d3
      .scaleOrdinal()
      .domain(data.map((d: IPieChart) => d.name))
      .range(["#EF476F", "#FFD166", "#06D6A0", "#118AB2", "#073B4C"]);

    const pie = d3
      .pie<IPieChart>()
      .sort(null)
      .value((d) => d.value);
    const radius = Math.min(width, height) / 2 - 1
    const arc = d3
      .arc<d3.PieArcDatum<IPieChart>>()
      .innerRadius(0)
      .outerRadius(radius);
    
    const labelRadius = radius * 0.8;

    const arcLabel = d3.arc<d3.PieArcDatum<IPieChart>>().innerRadius(labelRadius).outerRadius(labelRadius);
    const arcs = pie(data);

    const svg = d3
      .select(pieChartRef.current)
      .attr("width", width)
      .attr("height", height)
      .attr("viewBox", [-width / 2, -height / 2, width, height])
      .attr("style", "max-width: 100%; height: auto; font: 10px sans-serif;");

    svg
      .append("g")
      .selectAll()
      .data(arcs)
      .join("path")
      .attr("fill", (d) => color(d.data.name) as string)
      .attr("d", arc)
      .append("title")
      .text((d) => `${d.data.name}: ${d.data.value.toLocaleString("en-US")}`);

    svg
      .append("g")
      .attr("text-anchor", "middle")
      .selectAll()
      .data(arcs)
      .join("text")
      .attr("transform", (d) => `translate(${arcLabel.centroid(d)})`)
      .call((text) =>
        text
          .append("tspan")
          .attr("y", "-0.4em")
          .attr("font-weight", "bold")
          .text((d) => d.data.name)
      )
      .call((text) =>
        text
          .filter((d) => d.endAngle - d.startAngle > 0.25)
          .append("tspan")
          .attr("x", 0)
          .attr("y", "0.7em")
          .attr("fill-opacity", 0.7)
          .text((d) => d.data.value.toLocaleString("en-US"))
      );
  }

  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(() => {renderPieChart()}, [data]);

  return <svg ref={pieChartRef}></svg>;
}
