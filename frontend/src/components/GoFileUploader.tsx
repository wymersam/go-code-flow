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
  const [selectedFunc, setSelectedFunc] = useState<string | null>(null);
  const [enableSummary, setEnableSummary] = useState(false);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!file.name.endsWith(".go")) {
      alert("Please select a Go source file (.go)");
      return;
    }

    setError(null);
    setLoading(true);
    setSelectedFunc(null);

    const formData = new FormData();
    formData.append("file", file);
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
        const nodes = data.nodes.map((id) => ({ id }));
        onGraphData({ nodes, links: data.links });
        setSummaries(data.summaries);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  };

  return (
    <div className="container">
      <h2>Upload Go source file</h2>

      <label>
        <input
          type="checkbox"
          checked={enableSummary}
          onChange={() => setEnableSummary(!enableSummary)}
        />
        Generate function summaries <span>(uses OpenAI API)</span>
      </label>

      <input type="file" accept=".go" onChange={handleFileChange} />

      {loading && <p className="loading-text">Processing file...</p>}
      {error && <p className="error-text">{error}</p>}

      {enableSummary && Object.keys(summaries).length > 0 && (
        <div className="function-summary-container">
          <h3>Functions</h3>
          <ul>
            {Object.keys(summaries).map((funcName) => (
              <li
                key={funcName}
                onClick={() => setSelectedFunc(funcName)}
                className={`function-list-item ${
                  selectedFunc === funcName ? "selected" : ""
                }`}
              >
                {funcName}
              </li>
            ))}
          </ul>

          {selectedFunc && (
            <div className="function-summary-box">
              <h4>Summary for {selectedFunc}</h4>
              <p>{summaries[selectedFunc] || "No summary available."}</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default GoFileUploader;
