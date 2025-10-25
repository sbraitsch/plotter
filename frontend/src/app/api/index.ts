export async function fetchWithAuth<T = void>(
  url: string,
  options: RequestInit = {},
): Promise<T> {
  const userId = localStorage.getItem("session_token");
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
    ...(userId ? { "X-Token": userId } : {}),
  };

  const res = await fetch(url, { ...options, headers });

  if (!res.ok) {
    const errorText = await res.text();
    throw new Error(errorText || res.statusText);
  }

  if (res.status === 204 || res.headers.get("Content-Length") === "0") {
    return undefined as T;
  }

  return (await res.json()) as T;
}

export const BASE_URL =
  typeof window !== "undefined" && window.location.hostname === "localhost"
    ? "http://localhost:8080"
    : "https://plotter.sbraitsch.dev/api";
