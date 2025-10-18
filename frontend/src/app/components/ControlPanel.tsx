import React, { useState } from "react";
import "@/styles/ControlPanel.css";
import { PlayerData, PlayerUpdate, updatePlayerData } from "../api/player";
import { CloudUpload, Lock, Unlock, Cog, TestTubeDiagonal } from "lucide-react";
import PlotGrid from "./PlotGrid";
import { useAuth, User } from "../context/AuthContext";
import AdminModal from "./AdminModal";
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
  updatePlotAssignments: React.Dispatch<React.SetStateAction<Assignment[]>>;
}

export default function ControlPanel({
  user,
  playerData,
  updatePlayerPlot,
  updatePlotAssignments,
}: ControlPanelProps) {
  const { setUser } = useAuth();
  const [showNotification, setShowNotification] = useState(false);
  const [isAdminModalOpen, setIsAdminModalOpen] = useState(false);
  const [notificationContent, setNotificationContent] = useState("");

  const handleOpenModal = () => setIsAdminModalOpen(true);
  const handleCloseModal = () => setIsAdminModalOpen(false);

  const handleModalSubmit = async (minRank: number) => {
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

  const showAdminPanel = user?.isAdmin;

  const runOptimizer = async () => {
    const results = await getOptimizedAssignments();
    updatePlotAssignments(results);
    setNotificationContent("Displaying optimized assignments.");
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
          <PlotGrid player={playerData} updatePlayerPlot={updatePlayerPlot} />
        )}
        <div className="spacer"></div>
        {showAdminPanel && (
          <div className="admin-controls">
            <button className="admin-btn" onClick={handleOpenModal}>
              <Cog size={16} />
              Configure
            </button>
            <button className="admin-btn" onClick={runOptimizer}>
              <TestTubeDiagonal size={16} />
              Preview
            </button>
            <button className="admin-btn" onClick={lockCommunity}>
              {user?.community.locked ? (
                <Unlock size={16} />
              ) : (
                <Lock size={16} />
              )}
              {user?.community.locked ? "Unlock Community" : "Lock Community"}
            </button>
          </div>
        )}
        <button className="btn" onClick={handleSync}>
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
        onClose={handleCloseModal}
        onSubmit={handleModalSubmit}
      />
    </>
  );
}
