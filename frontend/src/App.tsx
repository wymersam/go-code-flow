import React from "react";
import GoFileUploader from "./components/GoFileUploader";
import Diagram from "./components/Diagram";
import "./app.css";

const App: React.FC = () => {
  const [graphData, setGraphData] = React.useState<{
    nodes: { id: string }[];
    links: any[];
  } | null>(null);

  return (
    <div>
      <h1>Go Function Call Graph</h1>
      <GoFileUploader onGraphData={setGraphData} />
      {graphData ? (
        <Diagram nodes={graphData.nodes} links={graphData.links} />
      ) : (
        <p>No graph loaded.</p>
      )}
    </div>
  );
};

export default App;
