import { BASE_URL, fetchWithAuth } from "./index";

export interface ValidateResponse {
  battletag: string;
  char: string;
  note: string;
  isAdmin: boolean;
  community: Community;
}

export type Community = {
  id: string;
  name: string;
  realm: string;
  locked: boolean;
  finalized: boolean;
};

export async function validateSession(): Promise<ValidateResponse> {
  const url = `${BASE_URL}/user/validate`;
  const data = await fetchWithAuth<ValidateResponse>(url);
  return data;
}
