import { BASE_URL, fetchWithAuth } from "./index";

interface CommunityData {
  members: PlayerData[];
}

export interface PlayerData {
  id: string;
  battletag: string;
  char: string;
  plotData: Record<number, number>;
}

export interface PlayerUpdate {
  battletag: string;
  plotData: Record<number, number>;
}

export interface PlotEntry {
  char: string;
  prio: number;
}

export interface CommunitySettings {
  officerRank: number;
  memberRank: number;
}

export async function getCommunityData(): Promise<PlayerData[]> {
  const url = `${BASE_URL}/community`;
  const data = await fetchWithAuth<CommunityData>(url);
  return data.members.map((p) => ({
    ...p,
    plotData: p.plotData || {},
  }));
}

export async function getCommunitySettings(): Promise<CommunitySettings> {
  const url = `${BASE_URL}/community/config`;
  const data = await fetchWithAuth<CommunitySettings>(url);
  return data;
}

export async function updatePlayerData(
  update: PlayerUpdate,
): Promise<PlayerData[]> {
  const url = `${BASE_URL}/user/update`;

  const data = await fetchWithAuth<CommunityData>(url, {
    method: "POST",
    body: JSON.stringify(update),
  });

  return data.members.map((p) => ({
    ...p,
    plotData: p.plotData || {},
  }));
}

export function buildPlotMap(
  players: PlayerData[],
): Record<number, PlotEntry[]> {
  const plotMap: Record<number, PlotEntry[]> = {};

  players.forEach((player) => {
    for (const [plotIdStr, prio] of Object.entries(player.plotData)) {
      const plotId = Number(plotIdStr);
      if (!plotMap[plotId]) {
        plotMap[plotId] = [];
      }
      plotMap[plotId].push({
        char: player.char,
        prio,
      });
    }
  });

  return plotMap;
}
