import { useState, ChangeEvent } from "react";
import "@/styles/AdminModal.css";

interface ClearModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: () => void;
}

export default function ClearModal({
  isOpen,
  onClose,
  onSubmit,
}: ClearModalProps) {
  if (!isOpen) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-container" onClick={(e) => e.stopPropagation()}>
        <h2 className="modal-title">Note</h2>

        <div className="text-content">
          Changes will only be persisted if you{" "}
          <div className="inline-btn">Sync</div> them after.
        </div>

        <button className="btn submit-btn" onClick={onSubmit}>
          Reset
        </button>
      </div>
    </div>
  );
}
