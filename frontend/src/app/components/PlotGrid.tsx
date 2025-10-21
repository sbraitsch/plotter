import { PlayerData } from "../api/player";
import "@/styles/PlotGrid.css";
import { getGradientColor, getLowestFreePriority, TOTAL_PLOTS } from "../utils";
import React, { useState } from "react";
import { fetchWithAuth, BASE_URL } from "../api";

interface PlotGridProps {
  player?: PlayerData;
  updatePlayerPlot: (plotId: number, value: number) => void;
}

export default function PlotGrid({ player, updatePlayerPlot }: PlotGridProps) {
  const plotIdToPriority = player
    ? Object.entries(player.plotData).reduce<Record<number, number>>(
        (acc, [plotId, plotPriority]) => {
          acc[Number(plotId)] = plotPriority;
          return acc;
        },
        {},
      )
    : [];

  return (
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
          >
            {priority ?? plotId}
          </div>
        );
      })}
    </div>
  );
}
