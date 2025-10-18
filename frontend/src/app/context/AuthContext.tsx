"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { Community, validateSession } from "../api/validate";

export type User = {
  battletag: string;
  community: Community;
  isAdmin: boolean;
};

interface AuthContextType {
  isKnown: boolean;
  user: User | undefined;
  setUser: React.Dispatch<React.SetStateAction<User | undefined>>;
  validateKnownUser: () => Promise<void>;
  loading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [isKnown, setIsKnown] = useState(false);
  const [user, setUser] = useState<User | undefined>(undefined);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem("session_token");
    const validate = async () => {
      if (token) {
        setIsKnown(true);
        await validateKnownUser();
      }
      setLoading(false);
    };
    validate();
  }, []);

  const validateKnownUser = async () => {
    try {
      const validationResponse = await validateSession();
      const u: User = { ...validationResponse };
      setUser(u);
    } catch (err) {
      localStorage.clear();
      setIsKnown(false);
      console.log(err);
    }
  };

  return (
    <AuthContext.Provider
      value={{ isKnown, user, setUser, validateKnownUser, loading }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within an AuthProvider");
  return ctx;
};
