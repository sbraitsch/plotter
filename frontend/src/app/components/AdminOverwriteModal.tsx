import { useState, ChangeEvent } from "react";
import "@/styles/AdminModal.css";
import "@/styles/AdminOverwriteModal.css";
import { PlayerData } from "../api/player";
import React, { useEffect, useRef, ReactNode } from "react";
import { Info, ChevronDown, ChevronUp } from "lucide-react";
import {
  Assignment,
  OverwriteAssignment,
  overwriteSingleAssignment,
} from "../api/optimizer";

interface AdminOverwriteModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: () => void;
  plot: number;
  player: PlayerData;
  assignments: Assignment[];
  communityMembers: PlayerData[];
}

export default function AdminOverwriteModal({
  isOpen,
  onClose,
  onSubmit,
  plot,
  player,
  assignments,
  communityMembers,
}: AdminOverwriteModalProps) {
  const [value, setValue] = useState<string>(player.battletag);
  const [error, setError] = useState<ReactNode | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const [infoOpen, setInfoOpen] = useState(false);

  useEffect(() => {
    if (isOpen) {
      inputRef.current?.focus();
      setValue(assignments.find((ass) => ass.plot === plot)?.btag || "");
    }
  }, [isOpen]);

  const handleFocus = () => {
    inputRef.current?.select();
  };

  if (!isOpen) return null;

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
  };

  const handleSelect = (member: string) => {
    setValue(member);
  };

  const handleClose = () => {
    setError(undefined);
    setValue("");
    onClose();
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const parts = value.trim().split(" ");
    const b = parts[0];
    const c = parts[1] ?? parts[0];

    const body: OverwriteAssignment = {
      btag: b,
      char: c,
      plot: plot,
    };
    await overwriteSingleAssignment(body);
    onSubmit();
  };

  const getAssignment = (btag: string): number | undefined => {
    return assignments.find((ass) => ass.btag === btag)?.plot;
  };

  return (
    <div className="modal-overlay" onClick={handleClose}>
      <div className="modal-container" onClick={(e) => e.stopPropagation()}>
        <form onSubmit={handleSubmit}>
          <h2 className="modal-title">Overwrite Assignment of Plot #{plot}</h2>
          <div className="modal-field">
            <input
              id="btag"
              ref={inputRef}
              type="text"
              value={value}
              onChange={handleChange}
              onFocus={handleFocus}
              min={1}
              max={53}
            />
          </div>
          <button
            type="button"
            onClick={() => setInfoOpen(!infoOpen)}
            className="hint"
          >
            <div>
              {infoOpen ? (
                <div>
                  <ul>
                    <li>
                      Select one of the community members below, or type the
                      BattleTag by hand.
                    </li>
                    <li>
                      If the BattleTag is not yet registered, a new community
                      member will be created if possible.
                    </li>
                    <li>
                      If you overwrite an assignment, the old asignee will be
                      homeless. Do not forget to assign them elsewhere.
                    </li>
                    <li>
                      You can add a name by putting it after the BattleTag,
                      separated by whitespace (e.g. [btag] [name]).
                    </li>
                  </ul>
                </div>
              ) : (
                <span
                  style={{
                    display: "inline-flex",
                    alignItems: "center",
                    gap: "0.5rem",
                  }}
                >
                  <Info />
                  Click to expand info
                </span>
              )}
            </div>
          </button>

          <div className="player-list-wrapper">
            <ul className="player-list">
              {communityMembers.map((member, idx) => (
                <li
                  key={idx}
                  className={`player-list-item ${value === `${member.battletag} ${member.char}` ? "selected" : ""}`}
                  onClick={() =>
                    handleSelect(`${member.battletag} ${member.char}`)
                  }
                >
                  <span>
                    {member.battletag}{" "}
                    <span style={{ color: "gray" }}>{member.char}</span>
                  </span>
                  <span>{getAssignment(member.battletag)}</span>
                </li>
              ))}
            </ul>
          </div>
          <button type="submit" className="submit-btn" disabled={!!error}>
            Confirm
          </button>
        </form>
      </div>
    </div>
  );
}
