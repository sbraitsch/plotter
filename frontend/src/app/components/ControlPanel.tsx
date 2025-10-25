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
} from "lucide-react";
import PlotGrid from "./PlotGrid";
import { useAuth, User } from "../context/AuthContext";
import AdminModal from "./AdminModal";
import ClearModal from "./ClearModal";
import { BASE_URL, fetchWithAuth } from "../api";
import ReactDOM from "react-dom";
import {
  Assignment,
  getOptimizedAssignments,
  optimizeAndLock,
} from "../api/optimizer";
import { getAppPageStaticInfo } from "next/dist/build/analysis/get-page-static-info";
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

  const handleAdminModalSubmit = async (minRank: number) => {
    try {
      await fetchWithAuth(`${BASE_URL}/community/config`, {
        method: "POST",
        body: JSON.stringify({ minRank: minRank }),
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

  useEffect(() => {
    if (!isPreviewing && !user?.community.locked) {
      updatePlotAssignments([]);
      setIsPreviewing(false);
      setNotificationContent("Preview stopped.");
    } else if (isPreviewing) {
      async function getAssigments() {
        const results = await getOptimizedAssignments();
        updatePlotAssignments(results);
        setNotificationContent("Previewing optimized assignments.");
      }
      getAssigments();
    }
    setShowNotification(true);
    setTimeout(() => setShowNotification(false), 5000);
  }, [isPreviewing]);

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
      battletag: playerData.battletag,
      plotData: playerData.plotData,
    };
    try {
      await updatePlayerData(update);
      setNotificationContent("Plot mapping updated!");
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
          {isPreviewing && (
            <button className="admin-btn" onClick={lockCommunity}>
              {user?.community.locked ? <Unlock /> : <Lock />}
            </button>
          )}
          {!user?.community.locked && (
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
            </>
          )}

          {/*{(isPreviewing || user?.community.locked) && (
            <button
              className={`toggle-wrapper
                               ${user?.community.locked ? "active-btn" : ""}`}
              onClick={lockCommunity}
            >
              <span className="toggle-label">Lock Community:</span>
              <label className="toggle">
                <input
                  className="toggle-checkbox"
                  type="checkbox"
                  checked={user?.community.locked}
                  onChange={lockCommunity}
                />
                <div className="toggle-switch">
                  <div className="toggle-thumb">
                    {user?.community.locked ? (
                      <Unlock size={16} />
                    ) : (
                      <Lock size={16} />
                    )}
                  </div>
                </div>
              </label>
            </button>
          )}*/}
        </div>
        <div className="spacer"></div>
        <button className="btn" disabled={!contextDirty} onClick={handleSync}>
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
