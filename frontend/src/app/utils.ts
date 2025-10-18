import { PlayerData } from "./api/player";

export const TOTAL_PLOTS = 53;

export const getGradientColor = (index: number) => {
  const hue = 120 - ((index - 1) / (TOTAL_PLOTS - 1)) * 120;
  return `hsl(${hue}, 70%, 60%)`;
};

export const getLowestFreePriority = (player: PlayerData) => {
  const usedPriorities = new Set(Object.values(player.plotData));

  let priority = 1;
  while (usedPriorities.has(priority)) {
    priority += 1;
  }

  return priority;
};
