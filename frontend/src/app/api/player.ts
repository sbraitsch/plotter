import { BASE_URL, fetchWithAuth } from "./index";

interface CommunityData {
  members: PlayerData[];
}

export interface PlayerData {
  id: string;
  battletag: string;
  plotData: Record<number, number>;
}

export interface PlayerUpdate {
  battletag: string;
  plotData: Record<number, number>;
}

export async function getCommunityData(): Promise<PlayerData[]> {
  const url = `${BASE_URL}/community`;
  const data = await fetchWithAuth<CommunityData>(url);
  return data.members.map((p) => ({
    ...p,
    plotData: p.plotData || {},
  }));
}

export async function updatePlayerData(
  update: PlayerUpdate,
): Promise<PlayerData[]> {
  const url = `${BASE_URL}/update`;

  const data = await fetchWithAuth<CommunityData>(url, {
    method: "POST",
    body: JSON.stringify(update),
  });

  return data.members.map((p) => ({
    ...p,
    plotData: p.plotData || {},
  }));
}
