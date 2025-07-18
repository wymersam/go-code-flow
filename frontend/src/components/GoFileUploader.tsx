import React, { useState } from "react";

interface ParseResponse {
  nodes: string[];
  links: { source: string; target: string }[];
  summaries: Record<string, string>;
}

interface GoFileUploaderProps {
  onGraphData: (data: { nodes: { id: string }[]; links: any[] }) => void;
}

const GoFileUploader: React.FC<GoFileUploaderProps> = ({ onGraphData }) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [summaries, setSummaries] = useState<Record<string, string>>({});
  const [highlightedFunc, setHighlightedFunc] = useState<string | null>(null);
  const [selectedFunction, setSelectedFunction] = useState<string>("");
  const [enableSummary, setEnableSummary] = useState(false);
  const [rawGraph, setRawGraph] = useState<ParseResponse | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!file.name.endsWith(".zip")) {
      alert("Please select a ZIP file (.zip)");
      return;
    }

    setError(null);
    setLoading(true);
    setHighlightedFunc(null);

    const formData = new FormData();
    formData.append("repo", file);
    formData.append("enableSummary", enableSummary ? "true" : "false");

    fetch("http://localhost:8080/parse", {
      method: "POST",
      body: formData,
    })
      .then(async (response) => {
        if (!response.ok) {
          throw new Error(await response.text());
        }
        return response.json() as Promise<ParseResponse>;
      })
      .then((data) => {
        setRawGraph(data);
        setSummaries(data.summaries);
        setLoading(false);

        // Optional: Show nothing until a function is selected
        onGraphData({
          nodes: [],
          links: [],
        });

        // Or show full graph by default:
        // onGraphData({
        //   nodes: data.nodes.map((id) => ({ id })),
        //   links: data.links,
        // });
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  };

  const getSubgraph = (
    startNode: string,
    links: { source: string; target: string }[]
  ): { nodes: string[]; links: { source: string; target: string }[] } => {
    const visited = new Set<string>();
    const queue = [startNode];
    const subLinks: { source: string; target: string }[] = [];

    while (queue.length > 0) {
      const current = queue.shift();
      if (!current || visited.has(current)) continue;
      visited.add(current);

      const outgoing = links.filter((link) => link.source === current);
      subLinks.push(...outgoing);

      for (const link of outgoing) {
        if (!visited.has(link.target)) {
          queue.push(link.target);
        }
      }
    }

    return {
      nodes: Array.from(visited),
      links: subLinks,
    };
  };

  const updateGraphFromFunction = (data: ParseResponse, func: string) => {
    let filteredNodes = data.nodes.map((id) => ({ id }));
    let filteredLinks = data.links;

    if (func) {
      const { nodes: ids, links } = getSubgraph(func, data.links);
      filteredNodes = ids.map((id) => ({ id }));
      filteredLinks = links;
    }

    onGraphData({
      nodes: filteredNodes,
      links: filteredLinks,
    });
  };

  const handleFunctionSelect = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const func = e.target.value;
    setSelectedFunction(func);
    if (rawGraph) {
      updateGraphFromFunction(rawGraph, func);
    }
  };

  return (
    <div className="container">
      <h2>Upload Go repository (.zip)</h2>

      <label>
        <input
          type="checkbox"
          checked={enableSummary}
          onChange={() => setEnableSummary(!enableSummary)}
        />
        Generate function summaries <span>(uses OpenAI API)</span>
      </label>

      <input type="file" accept=".zip" onChange={handleFileChange} />

      {loading && <p className="loading-text">Processing file...</p>}
      {error && <p className="error-text">{error}</p>}

      {rawGraph && rawGraph.nodes.length > 0 && (
        <div className="function-selector">
          <label>Select function to visualize:</label>
          <select value={selectedFunction} onChange={handleFunctionSelect}>
            <option value="">(All functions)</option>
            {rawGraph.nodes.map((name) => (
              <option key={name} value={name}>
                {name}
              </option>
            ))}
          </select>
        </div>
      )}

      {enableSummary && Object.keys(summaries).length > 0 && (
        <div className="function-summary-container">
          <h3>Functions</h3>
          <ul>
            {Object.keys(summaries).map((funcName) => (
              <li
                key={funcName}
                onClick={() => setHighlightedFunc(funcName)}
                className={`function-list-item ${
                  highlightedFunc === funcName ? "selected" : ""
                }`}
              >
                {funcName}
              </li>
            ))}
          </ul>

          {highlightedFunc && (
            <div className="function-summary-box">
              <h4>Summary for {highlightedFunc}</h4>
              <p>{summaries[highlightedFunc] || "No summary available."}</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default GoFileUploader;
