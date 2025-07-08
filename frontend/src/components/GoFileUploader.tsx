// src/components/GoFileUploader.tsx
import React, { useState } from "react";

interface GoFileUploaderProps {
  onGraphData: (data: { nodes: { id: string }[]; links: any[] }) => void;
}

const GoFileUploader: React.FC<GoFileUploaderProps> = ({ onGraphData }) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!file.name.endsWith(".go")) {
      alert("Please select a Go source file (.go)");
      return;
    }

    setError(null);
    setLoading(true);

    const formData = new FormData();
    formData.append("file", file);

    fetch("http://localhost:8080/parse", {
      method: "POST",
      body: formData,
    })
      .then(async (response) => {
        if (!response.ok) {
          throw new Error(await response.text());
        }
        return response.json();
      })
      .then((data) => {
        // Convert nodes from string[] to { id: string }[]
        const preparedData = {
          nodes: data.nodes.map((id: string) => ({ id })),
          links: data.links,
        };
        onGraphData(preparedData);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  };

  return (
    <div>
      <h2>Upload Go source file</h2>
      <input type="file" accept=".go" onChange={handleFileChange} />
      {loading && <p>Processing file...</p>}
      {error && <p style={{ color: "red" }}>{error}</p>}
    </div>
  );
};

export default GoFileUploader;
