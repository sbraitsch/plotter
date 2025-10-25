import { useState, ChangeEvent } from "react";
import "@/styles/AdminModal.css";

interface InfoModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export default function InfoModal({ isOpen, onClose }: InfoModalProps) {
  if (!isOpen) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-container" onClick={(e) => e.stopPropagation()}>
        <h2 className="modal-title">FAQ</h2>

        <div className="info-content">
          <ul>
            <li>
              Click a house on the map to assign the lowest free priority to it
            </li>
            <li>Click again to remove it</li>
            <li>
              When using Target Mode, clicking on a house opens a modal that
              lets you prioritize manually
            </li>
          </ul>
        </div>
        <button className="btn submit-btn" onClick={onClose}>
          Got it!
        </button>
      </div>
    </div>
  );
}
