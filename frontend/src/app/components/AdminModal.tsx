import { useState, ChangeEvent, useEffect } from "react";
import "@/styles/AdminModal.css";
import { getCommunitySettings } from "../api/player";

interface AdminModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (admin: number, member: number) => void;
}

export default function AdminModal({
  isOpen,
  onClose,
  onSubmit,
}: AdminModalProps) {
  const [adminValue, setAdminValue] = useState("");
  const [memberValue, setMemberValue] = useState("");

  useEffect(() => {
    if (!isOpen) return;
    async function fetchData() {
      try {
        const { officerRank, memberRank } = await getCommunitySettings();
        setAdminValue(officerRank.toString());
        setMemberValue(memberRank.toString());
      } catch (err: any) {
        console.error(err);
      }
    }
    fetchData();
  }, [isOpen]);

  if (!isOpen) return null;

  const handleAdminChange = (e: ChangeEvent<HTMLInputElement>) => {
    setAdminValue(e.target.value);
  };

  const handleMemberChange = (e: ChangeEvent<HTMLInputElement>) => {
    setMemberValue(e.target.value);
  };

  const handleSubmit = () => {
    const admin = Number(adminValue);
    const member = Number(memberValue);
    if (!isNaN(admin) && !isNaN(member)) {
      onSubmit(admin, member);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-container" onClick={(e) => e.stopPropagation()}>
        <h2 className="modal-title">Community Guidelines</h2>

        <div className="info-content">
          Set the guild rank requirements to...
        </div>
        <div className="modal-field">
          <label htmlFor="admin">...be admin: </label>
          <input
            id="admin"
            type="number"
            value={adminValue}
            onChange={handleAdminChange}
            min={0}
          />
        </div>
        <div className="modal-field">
          <label htmlFor="join">...join the community: </label>
          <input
            id="join"
            type="number"
            value={memberValue}
            onChange={handleMemberChange}
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
