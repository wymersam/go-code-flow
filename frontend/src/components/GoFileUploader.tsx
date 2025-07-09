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
    <div
      style={{
        background: "#f9f9f9",
        padding: "1.5rem",
        borderRadius: "12px",
        boxShadow: "0 2px 10px rgba(0,0,0,0.1)",
        marginBottom: "2rem",
      }}
    >
      <h2 style={{ marginBottom: "1rem", color: "#333" }}>
        Upload Go source file
      </h2>

      <label
        style={{
          display: "flex",
          alignItems: "center",
          marginBottom: "1rem",
          fontSize: "0.95rem",
          color: "#444",
        }}
      >
        <input
          type="checkbox"
          checked={enableSummary}
          onChange={() => setEnableSummary(!enableSummary)}
          style={{ marginRight: "0.5rem" }}
        />
        Generate function summaries{" "}
        <span
          style={{
            fontStyle: "italic",
            fontSize: "0.85rem",
            marginLeft: "0.5rem",
            color: "#888",
          }}
        >
          (uses OpenAI API)
        </span>
      </label>

      <input
        type="file"
        accept=".go"
        onChange={handleFileChange}
        style={{
          padding: "0.4rem",
          border: "1px solid #ccc",
          borderRadius: "6px",
          cursor: "pointer",
          marginBottom: "1rem",
        }}
      />

      {loading && <p style={{ color: "#0078D4" }}>Processing file...</p>}
      {error && <p style={{ color: "red" }}>{error}</p>}

      {enableSummary && Object.keys(summaries).length > 0 && (
        <div style={{ marginTop: "2rem" }}>
          <h3 style={{ marginBottom: "0.5rem", color: "#222" }}>Functions</h3>
          <ul style={{ paddingLeft: 0 }}>
            {Object.keys(summaries).map((funcName) => (
              <li
                key={funcName}
                onClick={() => setSelectedFunc(funcName)}
                style={{
                  listStyle: "none",
                  padding: "8px 12px",
                  marginBottom: 6,
                  border:
                    selectedFunc === funcName
                      ? "2px solid #0078D4"
                      : "1px solid #ddd",
                  borderRadius: 6,
                  backgroundColor:
                    selectedFunc === funcName ? "#eaf4fe" : "#fff",
                  cursor: "pointer",
                  transition: "background-color 0.2s ease",
                }}
              >
                {funcName}
              </li>
            ))}
          </ul>

          {selectedFunc && (
            <div
              style={{
                marginTop: "1rem",
                padding: "1rem",
                background: "#fff",
                border: "1px solid #ddd",
                borderRadius: "8px",
              }}
            >
              <h4 style={{ marginBottom: "0.5rem", color: "#0078D4" }}>
                Summary for {selectedFunc}
              </h4>
              <p
                style={{
                  fontSize: "0.95rem",
                  lineHeight: "1.4",
                  color: "#333",
                }}
              >
                {summaries[selectedFunc] || "No summary available."}
              </p>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default GoFileUploader;
