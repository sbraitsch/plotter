import { useState, ChangeEvent } from "react";
import "@/styles/AdminModal.css";

interface AdminModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (value: number) => void;
  initialValue?: number;
}

export default function AdminModal({
  isOpen,
  onClose,
  onSubmit,
  initialValue = 0,
}: AdminModalProps) {
  const [value, setValue] = useState<string>(initialValue?.toString() || "");

  if (!isOpen) return null;

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
  };

  const handleSubmit = () => {
    const num = Number(value);
    if (!isNaN(num)) {
      onSubmit(num);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-container" onClick={(e) => e.stopPropagation()}>
        <h2 className="modal-title">Configure</h2>

        <div className="modal-field">
          <label htmlFor="minRank">Min. Rank for Admin:</label>
          <input
            id="minRank"
            type="number"
            value={value}
            onChange={handleChange}
            min={0}
          />
        </div>

        <button className="btn submit-btn" onClick={handleSubmit}>
          Submit
        </button>
      </div>
    </div>
  );
}
