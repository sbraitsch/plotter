"use client";

import React, { useState, FormEvent } from "react";
import "@/styles/AuthModal.css";
import { BASE_URL } from "../api";

const AuthModal: React.FC = () => {
  const [loading, setLoading] = useState(false);

  const handleLogin = () => {
    const url = `${BASE_URL}/auth/bnet/login`;
    setLoading(true);
    window.location.href = url;
  };

  return (
    <div className="bnet-login-container">
      <button onClick={handleLogin} className="bnet-login-btn">
        <img src="/bnet.svg" alt="Battle.net" className="bnet-icon" />
        Log in with Battle.net
      </button>
      {loading && (
        <div className="spinner-container">
          <div className="spinner" />
        </div>
      )}
    </div>
  );
};

export default AuthModal;
