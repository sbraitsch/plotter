"use client";

import AuthModal from "./components/AuthModal";
import CommunitySelection from "./components/CommunitySelection";
import MapComponent from "./components/MapComponent";
import { AuthProvider, useAuth } from "./context/AuthContext";

export default function Home() {
  return (
    <AuthProvider>
      <HomeContent />
    </AuthProvider>
  );
}

const HomeContent: React.FC = () => {
  const { isKnown, user, loading } = useAuth();

  if (loading) return null;

  if (!isKnown) return <AuthModal />;
  if (!user?.community.id) return <CommunitySelection />;

  const mainContainerStyle: React.CSSProperties = {
    minHeight: "100vh",
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
  };

  const contentWrapperStyle: React.CSSProperties = {
    width: "100%",
    maxWidth: "90%",
    margin: "0 auto",
  };

  return (
    <div style={mainContainerStyle}>
      <div style={contentWrapperStyle}>
        <MapComponent />
      </div>
    </div>
  );
};
