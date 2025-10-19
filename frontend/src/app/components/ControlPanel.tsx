import React, { useState } from "react";
import "@/styles/ControlPanel.css";
import { PlayerData, PlayerUpdate, updatePlayerData } from "../api/player";
import {
  CloudUpload,
  Lock,
  Unlock,
  Cog,
  TestTubeDiagonal,
  TestTube,
  Trash2,
  Hand,
  Target,
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
  const [notificationContent, setNotificationContent] = useState("");
  const [isPreviewing, setIsPreviewing] = useState(user?.community.locked);

  const handleAdminModalSubmit = async (minRank: number) => {
    try {
      await fetchWithAuth(`${BASE_URL}/config`, {
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

  const runOptimizer = async () => {
    if (isPreviewing && !user?.community.locked) {
      updatePlotAssignments([]);
      setIsPreviewing(false);
      setNotificationContent("Preview stopped.");
    } else if (!isPreviewing) {
      const results = await getOptimizedAssignments();
      updatePlotAssignments(results);
      setIsPreviewing(true);
      setNotificationContent("Previewing optimized assignments.");
    } else {
      setNotificationContent("Community is locked. Can't leave Preview Mode.");
    }
    setShowNotification(true);
    setTimeout(() => setShowNotification(false), 5000);
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
          <div className="btag-value">{user?.battletag}</div>
          <div className="btag-community">&lt;{user?.community.name}&gt;</div>
        </div>
        {user?.community.locked ? (
          <div className="lock-notice">
            An admin has locked this community. Plot selection has been
            disabled.
          </div>
        ) : (
          <>
            <PlotGrid player={playerData} updatePlayerPlot={updatePlayerPlot} />
            <div className="btn-group">
              <button
                className="admin-btn"
                onClick={() => setIsClearModalOpen(true)}
              >
                <Trash2 />
                Clear
              </button>
              <button
                className={`admin-btn ${targetedMode ? "active-btn" : ""}`}
                onClick={toggleManualAssign}
              >
                <Target />
                Targeted
              </button>
            </div>
            <div className="lock-notice">
              Click a house on the map to assign the lowest free priority to it.
              <br />
              <br />
              Click it again to remove it.
              <br />
              <br />
              When using Targeted mode, clicking on a house opens a modal that
              lets you prioritize manually.
            </div>
          </>
        )}
        <div className="spacer"></div>
        {showAdminPanel && (
          <div className="btn-group">
            {(isPreviewing || user.community.locked) && (
              <button
                className={`admin-btn
                 ${user.community.locked ? "active-btn" : ""}`}
                onClick={lockCommunity}
              >
                {user.community.locked ? (
                  <Unlock size={16} />
                ) : (
                  <Lock size={16} />
                )}
                {user.community.locked ? "Unlock Community" : "Lock Community"}
              </button>
            )}
            <button
              className="admin-btn"
              onClick={() => setIsAdminModalOpen(true)}
            >
              <Cog color="grey" />
              Configure
            </button>
            <button
              className={`admin-btn ${isPreviewing ? "active-btn" : ""}`}
              onClick={runOptimizer}
            >
              {isPreviewing ? <TestTubeDiagonal /> : <TestTube />}
              Preview
            </button>
          </div>
        )}

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
    </>
  );
}
