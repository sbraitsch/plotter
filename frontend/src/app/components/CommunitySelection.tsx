"use client";

import React, { useState, useEffect } from "react";
import "@/styles/CommunitySelection.css";
import { BASE_URL, fetchWithAuth } from "../api";
import { useAuth, User } from "../context/AuthContext";
import { Community } from "../api/validate";

interface CommunityResponse {
  id: string;
  name: string;
  realm: string;
  locked: boolean;
}
const CommunitySelection: React.FC = () => {
  const { user, setUser } = useAuth();
  const [options, setOptions] = useState<CommunityResponse[]>([]);
  const [selected, setSelected] = useState<CommunityResponse | undefined>(
    undefined,
  );

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchOptions = async () => {
      setLoading(true);
      try {
        const communities = await fetchWithAuth<CommunityResponse[]>(
          `${BASE_URL}/auth/bnet/guilds`,
        );

        setOptions(communities);
      } catch (err) {
        console.error(err);
        setError("Failed to load options");
      } finally {
        setLoading(false);
      }
    };

    fetchOptions();
  }, []);

  const handleSelect = (com: CommunityResponse) => {
    setSelected(com);
  };

  const handleSubmit = async (com: CommunityResponse) => {
    try {
      await fetchWithAuth<Community>(`${BASE_URL}/community/join/${com.id}`, {
        method: "POST",
      });

      setUser((prev) => {
        if (!prev) return prev;
        return {
          ...prev,
          community: {
            id: com.id,
            name: com.name,
            realm: com.realm,
            locked: com.locked,
          },
        };
      });
    } catch (err) {
      setError("Failed to join community");
    }
  };

  return (
    <div className="bnet-list-wrapper">
      <h2 className="bnet-list-title">Choose your Community</h2>
      <p className="bnet-list-subtitle">
        Don't fuck this up. You won't be able to change it later.
      </p>

      {loading && <p>Loading...</p>}
      {error && <p className="error">{error}</p>}

      {!loading && !error && (
        <ul className="bnet-list">
          {options.map((opt) => (
            <li
              key={opt.id}
              className={`bnet-list-item ${selected?.id === opt.id ? "selected" : ""}`}
              onClick={() => handleSelect(opt)}
            >
              {opt.name}
            </li>
          ))}
        </ul>
      )}
      <div className="spacer"></div>
      <button
        className="bnet-submit-btn"
        disabled={!selected}
        onClick={() => {
          handleSubmit(selected!!);
        }}
      >
        Continue
      </button>
    </div>
  );
};

export default CommunitySelection;
