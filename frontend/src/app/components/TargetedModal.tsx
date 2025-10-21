import { useState, ChangeEvent } from "react";
import "@/styles/AdminModal.css";
import { PlayerData } from "../api/player";
import { getLowestFreePriority } from "../utils";
import React, { useEffect, useRef, ReactNode } from "react";

interface TargetedModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (plot: number, prio: number) => void;
  plot: number;
  player: PlayerData;
}

export default function TargetedModal({
  isOpen,
  onClose,
  onSubmit,
  plot,
  player,
}: TargetedModalProps) {
  const [value, setValue] = useState<string>(
    getLowestFreePriority(player).toString() || "",
  );
  const [error, setError] = useState<ReactNode | null>(null);
  const [showConfirm, setShowConfirm] = useState<boolean>(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const reassignError = (plot: number) => {
    return (
      <>
        <span
          style={{
            display: "inline-flex",
            alignItems: "center",
          }}
        >
          Priority already assigned to
          <img
            src={"./house_pop_40.png"}
            style={{ width: 28, height: 28, marginLeft: 4 }}
          />{" "}
          <span style={{ fontWeight: "bold" }}>#{plot}</span>
        </span>
        <br />
        <br />
        <span style={{ color: "white" }}>
          Click <span className="inline-btn">Confirm</span> to reassign.
        </span>
      </>
    );
  };

  useEffect(() => {
    setValue(getLowestFreePriority(player).toString() || "");
  }, [player]);

  useEffect(() => {
    if (isOpen) {
      inputRef.current?.focus();
    }
  }, [isOpen]);

  const handleFocus = () => {
    inputRef.current?.select();
  };

  if (!isOpen) return null;

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    checkValidity(val);
    setValue(e.target.value);
  };

  const handleClose = () => {
    setError(undefined);
    setShowConfirm(false);
    setValue(getLowestFreePriority(player).toString() || "");
    onClose();
  };

  const checkValidity = (value: string) => {
    const num = Number(value);
    const maybePlot = isUsed(player.plotData, num);
    if (!isNaN(num) && isInRange(num) && !maybePlot) {
      setError(undefined);
      setShowConfirm(false);
    } else if (!isInRange(num)) {
      setError(`The priority is out of the accepted bounds.`);
      setShowConfirm(false);
    } else {
      setError(reassignError(maybePlot!));
      setShowConfirm(true);
    }
    return false;
  };

  const isUsed = (
    record: Record<number, number>,
    num: number,
  ): number | undefined => {
    for (const [key, value] of Object.entries(record)) {
      if (value === num) return Number(key);
    }
    return undefined;
  };

  const isInRange = (num: number) => num > 0 && num <= 53;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const num = Number(value);

    onSubmit(plot, num);
  };

  return (
    <div className="modal-overlay" onClick={handleClose}>
      <div className="modal-container" onClick={(e) => e.stopPropagation()}>
        <form onSubmit={handleSubmit}>
          <h2 className="modal-title">Modify Priority</h2>

          <div className="modal-field">
            <label
              htmlFor="priority"
              style={{
                display: "inline-flex",
                alignItems: "center",
              }}
            >
              Set the priority for{" "}
              <img
                src={"./house_pop_40.png"}
                style={{ width: 28, height: 28, marginLeft: 4 }}
              />{" "}
              <span style={{ fontWeight: "bold" }}> #{plot}</span>
            </label>
            <input
              id="priority"
              ref={inputRef}
              type="number"
              value={value}
              onChange={handleChange}
              onFocus={handleFocus}
              min={1}
              max={53}
            />
          </div>

          <div className="input-error">{error}</div>
          {showConfirm ? (
            <>
              <div className="btn-group">
                <button
                  type="button"
                  className="btn cancel-btn"
                  onClick={handleClose}
                >
                  Cancel
                </button>
                <button type="submit" className="btn submit-btn">
                  Confirm
                </button>
              </div>
            </>
          ) : (
            <button type="submit" className="submit-btn" disabled={!!error}>
              Confirm
            </button>
          )}
        </form>
      </div>
    </div>
  );
}
