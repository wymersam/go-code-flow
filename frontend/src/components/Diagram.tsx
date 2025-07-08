import React, { useEffect, useRef } from "react";
import * as d3 from "d3";

type NodeType = {
  id: string;
  x?: number;
  y?: number;
  fx?: number | null;
  fy?: number | null;
};
type LinkType = { source: string | NodeType; target: string | NodeType };

type DiagramProps = {
  nodes: NodeType[];
  links: LinkType[];
};

const Diagram: React.FC<DiagramProps> = ({ nodes, links }) => {
  const svgRef = useRef<SVGSVGElement | null>(null);
  const gRef = useRef<SVGGElement | null>(null);

  useEffect(() => {
    if (!nodes.length || !links.length || !svgRef.current || !gRef.current)
      return;

    const svg = d3.select(svgRef.current);
    const g = d3.select(gRef.current);
    svg.selectAll("g > *").remove(); // clear inside <g> only

    const width = svgRef.current.clientWidth;
    const height = svgRef.current.clientHeight;

    // Zoom/pan behavior
    svg.call(
      d3.zoom<SVGSVGElement, unknown>().on("zoom", (event) => {
        g.attr("transform", event.transform);
      })
    );

    // Create simulation
    const simulation = d3
      .forceSimulation(nodes as any)
      .force(
        "link",
        d3
          .forceLink(links as any)
          .id((d: any) => d.id)
          .distance(100)
      )
      .force("charge", d3.forceManyBody().strength(-300))
      .force("center", d3.forceCenter(width / 2, height / 2));

    // Draw links
    const link = g
      .append("g")
      .attr("stroke", "#999")
      .attr("stroke-opacity", 0.6)
      .selectAll("line")
      .data(links)
      .join("line")
      .attr("stroke-width", 1.5);

    // Draw nodes
    const node = g
      .append("g")
      .attr("stroke", "#fff")
      .attr("stroke-width", 1.5)
      .selectAll("circle")
      .data(nodes)
      .join("circle")
      .attr("r", 10)
      .attr("fill", "#69b3a2")
      .call(
        d3
          .drag()
          .on("start", (event, d: any) => {
            if (!event.active) simulation.alphaTarget(0.3).restart();
            d.fx = d.x;
            d.fy = d.y;
          })
          .on("drag", (event, d: any) => {
            d.fx = event.x;
            d.fy = event.y;
          })
          .on("end", (event, d: any) => {
            if (!event.active) simulation.alphaTarget(0);
            d.fx = null;
            d.fy = null;
          })
      );

    // Labels
    const label = g
      .append("g")
      .selectAll("text")
      .data(nodes)
      .join("text")
      .text((d) => d.id)
      .attr("font-size", 10)
      .attr("dx", 12)
      .attr("dy", ".35em");

    // Tick
    simulation.on("tick", () => {
      link
        .attr("x1", (d: any) => (d.source as any).x)
        .attr("y1", (d: any) => (d.source as any).y)
        .attr("x2", (d: any) => (d.target as any).x)
        .attr("y2", (d: any) => (d.target as any).y);

      node.attr("cx", (d: any) => d.x).attr("cy", (d: any) => d.y);
      label.attr("x", (d: any) => d.x).attr("y", (d: any) => d.y);
    });

    return () => simulation.stop();
  }, [nodes, links]);

  return (
    <svg ref={svgRef} style={{ width: "100%", height: "100vh" }}>
      <g ref={gRef}></g>
    </svg>
  );
};

export default Diagram;
