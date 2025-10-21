import { BASE_URL, fetchWithAuth } from "./index";

export interface ValidateResponse {
  battletag: string;
  isAdmin: boolean;
  community: Community;
}

export type Community = {
  id: string;
  name: string;
  locked: boolean;
};

export async function validateSession(): Promise<ValidateResponse> {
  const url = `${BASE_URL}/user/validate`;
  const data = await fetchWithAuth<ValidateResponse>(url);
  return data;
}
