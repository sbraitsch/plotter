import { BASE_URL, fetchWithAuth } from ".";

export interface Assignment {
  btag: string;
  char: string;
  plot: number;
  score: number;
}

export interface OverwriteAssignment {
  btag: string;
  char: string;
  plot: number;
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

export async function finalizeAssignments(): Promise<void> {
  try {
    const url = `${BASE_URL}/community/finalize`;
    await fetchWithAuth(url, {
      method: "POST",
    });
    return;
  } catch (err) {
    throw new Error("Finalizing failed");
  }
}

export async function downloadAssignmentData(): Promise<void> {
  try {
    const userId = localStorage.getItem("session_token");
    const headers: HeadersInit = {
      "Content-Type": "application/json",
      ...(userId ? { "X-Token": userId } : {}),
    };
    const url = `${BASE_URL}/community/download`;
    await fetch(url, { headers })
      .then(async (res) => {
        if (!res.ok) throw new Error("Download failed");
        return res.blob().then((blob) => ({ blob, res }));
      })
      .then(({ blob, res }) => {
        const disposition = res.headers.get("Content-Disposition");
        let filename = "community_data.json";
        if (disposition && disposition.includes("filename=")) {
          const match = disposition.match(/filename="(.+)"/);
          if (match) filename = match[1];
        }

        const url = window.URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        a.remove();
        window.URL.revokeObjectURL(url);
      })
      .catch(console.error);
  } catch (err) {
    throw new Error("Download failed");
  }
}

export async function overwriteAssignments(json: any): Promise<Assignment[]> {
  try {
    const url = `${BASE_URL}/community/upload`;
    const data = await fetchWithAuth<Assignment[]>(url, {
      method: "POST",
      body: JSON.stringify(json),
    });
    return data;
  } catch (err) {
    throw new Error("Upload failed");
  }
}

export async function overwriteSingleAssignment(
  body: OverwriteAssignment,
): Promise<void> {
  try {
    const url = `${BASE_URL}/community/assignments`;
    await fetchWithAuth<Assignment[]>(url, {
      method: "POST",
      body: JSON.stringify(body),
    });
    return;
  } catch (err) {
    throw new Error("Overwrite failed");
  }
}

export async function optimizeAndLock(): Promise<Assignment[]> {
  try {
    const url = `${BASE_URL}/community/lock`;
    const data = await fetchWithAuth<Assignment[]>(url, {
      method: "POST",
    });
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
