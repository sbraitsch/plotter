"use client";
import { useEffect, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";

function AuthSuccessInner() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const token = searchParams.get("token");

    console.log(token);
    if (!token || typeof token !== "string") {
      router.replace("/");
      return;
    }

    localStorage.setItem("session_token", token);
    router.replace("/");
  }, [router, searchParams]);

  return <div>Logging in...</div>;
}

export default function AuthSuccess() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <AuthSuccessInner />
    </Suspense>
  );
}
