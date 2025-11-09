import "@/styles/AdminModal.css";

interface InfoModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export default function InfoModal({ isOpen, onClose }: InfoModalProps) {
  if (!isOpen) return null;

  const handleClose = () => {
    localStorage.removeItem("showInfoModal");
    onClose();
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-container" onClick={(e) => e.stopPropagation()}>
        <h2 className="modal-title">FAQ</h2>

        <div className="info-content">
          <ul>
            <li>1 = Highest Priority, 53 = Lowest Priority</li>
            <li>
              Clicking a house on the map will assign the highest free priority
              to it.
            </li>
            <li>
              Clicking on a house that is already prioritized will remove the
              priority.
            </li>
            <li>
              Any local changes are only persisted after clicking the Sync
              button.
            </li>
            <li>
              When using Target Mode, clicking on a house will open a window
              where you can manually input a priority.
            </li>
            <li>
              Hovering over a house will show you a list of interested people.
            </li>
          </ul>
        </div>
        <button className="btn submit-btn" onClick={handleClose}>
          Got it!
        </button>
      </div>
    </div>
  );
}
