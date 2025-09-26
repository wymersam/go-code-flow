import React from "react";
import GoFileUploader from "./components/GoFileUploader";
import Diagram from "./components/Diagram";
import "./app.css";

const App: React.FC = () => {
  const [graphData, setGraphData] = React.useState<{
    nodes: { id: string }[];
    links: any[];
  } | null>(null);
  const [summaries, setSummaries] = React.useState<Record<string, string>>({});
  const [showSummaries, setShowSummaries] = React.useState(false);
  const [highlightedFunc, setHighlightedFunc] = React.useState<string | null>(
    null
  );

  const handleSummariesUpdate = (
    newSummaries: Record<string, string>,
    show: boolean
  ) => {
    setSummaries(newSummaries);
    setShowSummaries(show && Object.keys(newSummaries).length > 0);
  };

  return (
    <div className="app-container">
      <header className="app-header">
        <h1>Go Function Call Graph</h1>
        <GoFileUploader
          onGraphData={setGraphData}
          onSummariesUpdate={handleSummariesUpdate}
        />
      </header>
      <main className="app-main">
        {graphData ? (
          <Diagram nodes={graphData.nodes} links={graphData.links} />
        ) : (
          <div className="no-graph-message">
            <p>No graph loaded.</p>
          </div>
        )}

        {showSummaries && (
          <aside className="summaries-sidebar">
            <div className="summaries-content">
              <h3>Function Summaries</h3>
              <button
                className="close-summaries"
                onClick={() => setShowSummaries(false)}
              >
                ×
              </button>
              <ul className="summaries-list">
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
          </aside>
        )}
      </main>
    </div>
  );
};

export default App;
