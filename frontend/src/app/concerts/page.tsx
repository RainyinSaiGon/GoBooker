"use client";

import { useState, useEffect } from "react";
import { useAuthStore } from "@/lib/store/authStore";
import LoginForm from "@/components/LoginForm";

export default function ConcertsPage() {
    const[mounted, setMounted] = useState(false);
    const { token } = useAuthStore();

    useEffect(() => {
        setMounted(true);
        }, []);

        // Prevent SSR flickering 
        if (!mounted) {
        return (
            <div className="flex-center animate-fade-in" style={{ minHeight: "100vh" }}>
            <div className="text-muted">Loading GoBooker...</div>
            </div>
        );
        }

        // Login View
        if (!token) {
        return <LoginForm />;
        }


    return (
        <div>
            <h1 className="text-3xl font-bold mb-4">Concerts</h1>
            <p className="text-lg text-gray-600">This is the concerts page.</p>
        </div>
    );    
}