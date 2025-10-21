import { BASE_URL, fetchWithAuth } from ".";

export interface Assignment {
  player: string;
  plot: number;
  score: number;
}

export async function getOptimizedAssignments(): Promise<Assignment[]> {
  try {
    const url = `${BASE_URL}/community/optimize`;
    const data = await fetchWithAuth<Assignment[]>(url);
    return data;
  } catch (err) {
    throw new Error("Optimizer failed");
  }
}

export async function optimizeAndLock(): Promise<Assignment[]> {
  try {
    const url = `${BASE_URL}/community/lock`;
    const data = await fetchWithAuth<Assignment[]>(url, { method: "POST" });
    return data;
  } catch (err) {
    throw new Error("Optimizer failed");
  }
}

export async function getAssignedPlots(): Promise<Assignment[]> {
  try {
    const url = `${BASE_URL}/community/assignments`;
    const data = await fetchWithAuth<Assignment[]>(url);
    return data;
  } catch (err) {
    throw new Error("Failed to fetch assignment data.");
  }
}
