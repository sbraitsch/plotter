import React, { useState, useEffect } from "react";
import "@/styles/ControlPanel.css";
import { PlayerData, PlayerUpdate, updatePlayerData } from "../api/player";
import {
  CloudUpload,
  Lock,
  Unlock,
  Cog,
  Info,
  TestTubeDiagonal,
  TestTube,
  Trash2,
  PowerOff,
  Power,
  Download,
  Upload,
  NotebookPen,
} from "lucide-react";
import PlotGrid from "./PlotGrid";
import { useAuth, User } from "../context/AuthContext";
import AdminModal from "./AdminModal";
import ClearModal from "./ClearModal";
import { BASE_URL, fetchWithAuth } from "../api";
import ReactDOM from "react-dom";
import {
  Assignment,
  downloadAssignmentData,
  getOptimizedAssignments,
  optimizeAndLock,
  overwriteAssignments,
} from "../api/optimizer";
import InfoModal from "./InfoModal";

interface ControlPanelProps {
  user: User | undefined;
  playerData?: PlayerData;
  updatePlayerPlot: (plotId: number, value: number) => void;
  clearPlayerMappings: () => void;
  updatePlotAssignments: React.Dispatch<React.SetStateAction<Assignment[]>>;
  targetedMode: boolean;
  setTargetedMode: React.Dispatch<React.SetStateAction<boolean>>;
  contextDirty: boolean;
}

export default function ControlPanel({
  user,
  playerData,
  updatePlayerPlot,
  clearPlayerMappings,
  updatePlotAssignments,
  targetedMode,
  setTargetedMode,
  contextDirty,
}: ControlPanelProps) {
  const { setUser } = useAuth();
  const [showNotification, setShowNotification] = useState(false);
  const [isAdminModalOpen, setIsAdminModalOpen] = useState(false);
  const [isClearModalOpen, setIsClearModalOpen] = useState(false);
  const [isInfoModalOpen, setIsInfoModalOpen] = useState(false);
  const [notificationContent, setNotificationContent] = useState("");
  const [isPreviewing, setIsPreviewing] = useState(user?.community.locked);
  const [note, setNote] = useState(user?.note || "");
  const [noteEdited, setNoteEdited] = useState(false);

  const handleAdminModalSubmit = async (admin: number, member: number) => {
    try {
      await fetchWithAuth(`${BASE_URL}/community/config`, {
        method: "POST",
        body: JSON.stringify({ adminRank: admin, memberRank: member }),
      });
      setNotificationContent("Community guidelines updated.");
    } catch (err) {
      setNotificationContent("Error updating community guidelines.");
    }
    setShowNotification(true);
    setTimeout(() => setShowNotification(false), 5000);
    setIsAdminModalOpen(false);
  };

  const handleClearModalSubmit = () => {
    clearPlayerMappings();
    setNotificationContent("Plot mappings cleared.");
    setShowNotification(true);
    setTimeout(() => setShowNotification(false), 5000);
    setIsClearModalOpen(false);
  };

  const showAdminPanel = user?.isAdmin;

  const toggleManualAssign = () => {
    setTargetedMode((prev) => !prev);
  };

  const togglePreview = () => {
    setIsPreviewing((prev) => !prev);
  };

  const handleNoteChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setNote(e.target.value);
    setNoteEdited(true);
  };

  useEffect(() => {
    if (localStorage.getItem("showInfoModal")) {
      setIsInfoModalOpen(true);
    }
  });

  useEffect(() => {
    if (!isPreviewing && !user?.community.locked) {
      updatePlotAssignments([]);
      setIsPreviewing(false);
    } else if (isPreviewing) {
      async function getAssigments() {
        const results = await getOptimizedAssignments();
        updatePlotAssignments(results);
        setNotificationContent("Previewing optimized assignments.");
        setShowNotification(true);
        setTimeout(() => setShowNotification(false), 5000);
      }
      getAssigments();
    }
  }, [isPreviewing]);

  const triggerDownload = async () => {
    downloadAssignmentData();
  };

  const triggerUpload = async () => {
    const input = document.createElement("input");
    input.type = "file";
    input.accept = "application/json";
    input.click();

    input.onchange = async () => {
      const file = input.files?.[0];
      if (!file) return;

      try {
        const text = await file.text();
        const jsonData = JSON.parse(text);

        const data = await overwriteAssignments(jsonData);
        updatePlotAssignments(data);
        setNotificationContent("Upload successful.");
      } catch (err) {
        console.log(err);
        setNotificationContent("Upload failed.");
      } finally {
        setShowNotification(true);
        setTimeout(() => setShowNotification(false), 5000);
      }
    };
  };

  const lockCommunity = async () => {
    const results = await optimizeAndLock();
    if (!user?.community.locked) {
      updatePlotAssignments(results);
      setNotificationContent(
        "Displaying optimized assignments. Community locked.",
      );
    } else {
      updatePlotAssignments([]);
      setIsPreviewing(false);
      setNotificationContent("Community unlocked.");
    }
    setUser((prev) =>
      prev
        ? {
            ...prev,
            community: {
              ...prev.community,
              locked: !prev.community.locked,
            },
          }
        : prev,
    );
    setShowNotification(true);
    setTimeout(() => setShowNotification(false), 5000);
  };

  const handleSync = async () => {
    if (!playerData) return;
    const update: PlayerUpdate = {
      note: note,
      plotData: playerData.plotData,
    };
    try {
      await updatePlayerData(update);
      setNotificationContent("Plot mapping updated!");
      setNoteEdited(false);
    } catch (err) {
      setNotificationContent("Error updating plot mapping.");
    }

    setShowNotification(true);
    setTimeout(() => setShowNotification(false), 5000);
  };

  return (
    <>
      <div className="info-panel">
        <div className="btag-tile">
          <div className="btag-label">{user?.battletag}</div>
          <div className="btag-value">{user?.char}</div>
          <div className="btag-community">&lt;{user?.community.name}&gt;</div>
        </div>
        <div className="btn-group">
          {isPreviewing && user?.community.locked && (
            <>
              <button className="admin-btn" onClick={triggerDownload}>
                <Download />
              </button>
              <button className="admin-btn" onClick={triggerUpload}>
                <Upload />
              </button>
            </>
          )}
          {isPreviewing && (
            <button className="admin-btn" onClick={lockCommunity}>
              {user?.community.locked ? <Unlock /> : <Lock />}
            </button>
          )}
          {!user?.community.locked && !isPreviewing && (
            <button
              className="admin-btn"
              onClick={() => setIsClearModalOpen(true)}
            >
              <Trash2 />
            </button>
          )}
          <button
            className="admin-btn"
            onClick={() => setIsInfoModalOpen(true)}
          >
            <Info />
          </button>
          {showAdminPanel && (
            <button
              className="admin-btn"
              onClick={() => setIsAdminModalOpen(true)}
            >
              <Cog />
            </button>
          )}
        </div>
        {user?.community.locked ? (
          <div className="lock-notice">
            An admin has locked this community. Plot selection has been
            disabled.
          </div>
        ) : (
          <PlotGrid player={playerData} updatePlayerPlot={updatePlayerPlot} />
        )}

        <div className="toggle-group">
          {!user?.community.locked && (
            <>
              <button
                className={`toggle-wrapper ${!isPreviewing ? (targetedMode ? "active-btn" : "") : "locked-btn"}`}
                onClick={toggleManualAssign}
                disabled={isPreviewing}
              >
                <span className="toggle-label">Target Mode:</span>
                <label className="toggle">
                  <input
                    className="toggle-checkbox"
                    type="checkbox"
                    checked={targetedMode}
                    onChange={toggleManualAssign}
                    disabled={isPreviewing}
                  />
                  <div className="toggle-switch">
                    <div className="toggle-thumb">
                      {targetedMode ? (
                        <Power size={16} />
                      ) : (
                        <PowerOff size={16} />
                      )}
                    </div>
                  </div>
                </label>
              </button>
              {showAdminPanel && (
                <button
                  className={`toggle-wrapper ${isPreviewing ? "active-btn" : ""}`}
                  onClick={togglePreview}
                  disabled={user?.community.locked}
                >
                  <span className="toggle-label">Preview:</span>
                  <label className="toggle">
                    <input
                      className="toggle-checkbox"
                      type="checkbox"
                      checked={isPreviewing}
                      onChange={togglePreview}
                      disabled={user?.community.locked}
                    />
                    <div className="toggle-switch">
                      <div className="toggle-thumb">
                        {isPreviewing ? (
                          <TestTubeDiagonal size={16} />
                        ) : (
                          <TestTube size={16} />
                        )}
                      </div>
                    </div>
                  </label>
                </button>
              )}
            </>
          )}
        </div>
        <div className="textarea-wrapper">
          <textarea
            id="info"
            className="info-textarea"
            placeholder="Add a note"
            value={note}
            onChange={handleNoteChange}
            spellCheck={false}
          ></textarea>
          <label htmlFor="info" className="info-label">
            <NotebookPen />
          </label>
        </div>
        <button
          className="btn"
          disabled={!contextDirty && !noteEdited}
          onClick={handleSync}
        >
          <CloudUpload />
          Sync
        </button>
      </div>
      {showNotification &&
        ReactDOM.createPortal(
          <div
            style={{
              position: "fixed",
              inset: 0,
              display: "flex",
              justifyContent: "center",
              alignItems: "flex-end",
              pointerEvents: "none",
              paddingBottom: "16px",
            }}
          >
            <div className="notification-popup">{notificationContent}</div>
          </div>,
          document.body,
        )}
      <AdminModal
        isOpen={isAdminModalOpen}
        onClose={() => setIsAdminModalOpen(false)}
        onSubmit={handleAdminModalSubmit}
      />
      <ClearModal
        isOpen={isClearModalOpen}
        onClose={() => setIsClearModalOpen(false)}
        onSubmit={handleClearModalSubmit}
      />
      <InfoModal
        isOpen={isInfoModalOpen}
        onClose={() => setIsInfoModalOpen(false)}
      />
    </>
  );
}
