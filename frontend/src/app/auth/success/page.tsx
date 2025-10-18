"use client";
import { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";

const AuthSuccess = () => {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const token = searchParams.get("token");

    if (!token || typeof token !== "string") {
      router.replace("/");
      return;
    }

    localStorage.setItem("session_token", token);

    router.replace("/");
  }, [router]);

  return <div>Logging in...</div>;
};

export default AuthSuccess;
