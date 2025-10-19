import { PlayerData } from "../api/player";
import "@/styles/PlotGrid.css";
import { getGradientColor, getLowestFreePriority, TOTAL_PLOTS } from "../utils";
import React, { useState } from "react";

interface PlotGridProps {
  player?: PlayerData;
  updatePlayerPlot: (plotId: number, value: number) => void;
}

export default function PlotGrid({ player, updatePlayerPlot }: PlotGridProps) {
  const [deleteMode, setDeleteMode] = useState(false);
  const plotIdToPriority = player
    ? Object.entries(player.plotData).reduce<Record<number, number>>(
        (acc, [plotId, plotPriority]) => {
          acc[Number(plotId)] = plotPriority;
          return acc;
        },
        {},
      )
    : [];

  const handleModeToggle = () => {
    setDeleteMode((prev) => !prev);
  };

  const handlePriorityUpdate = (plotId: number) => {
    // updatePlayerPlot(plotId, getLowestFreePriority(player!));
  };

  return (
    <>
      {/*<label className="switch">
        <input
          type="checkbox"
          checked={deleteMode}
          onChange={handleModeToggle}
        />
        <span className="slider round"></span>
      </label>*/}
      <div className="plot-grid">
        {Array.from({ length: TOTAL_PLOTS }, (_, i) => {
          const plotId = i + 1;
          const priority = plotIdToPriority[plotId];
          const bgColor = plotId ? getGradientColor(priority) : undefined;

          return (
            <div
              key={plotId}
              className={`plot-node ${priority ? "plot-node--active" : ""}`}
              style={priority ? { backgroundColor: bgColor } : undefined}
              onClick={() => handlePriorityUpdate(plotId)}
            >
              {priority ?? plotId}
            </div>
          );
        })}
      </div>
    </>
  );
}
