"use client";

import { useState, useEffect } from "react";
import { Card, CardHeader, CardContent, CardFooter } from "@/components/ui/registry/card";
import { Button } from "@/components/ui/registry/button";
import { cn } from "@/lib/utils";

const COOKIE_NAME = "hasConsent";
const COOKIE_VALUE_ACCEPT = "accepted";
const COOKIE_VALUE_REJECT = "rejected";
const COOKIE_EXPIRATION_TIME = 31104000000; // 12 months in milliseconds
const trackingCookiesNames = ["_ga", "_gid", "_gat"]; // List of tracking cookies to delete if rejected

export default function CookieBanner() {
  const [isBannerVisible, setIsBannerVisible] = useState(false);

  useEffect(() => {
    // Check if user has already accepted or rejected cookies
    const consent = localStorage.getItem(COOKIE_NAME) || getCookie(COOKIE_NAME);
    const isBot = /bot|crawler|spider|crawling/i.test(navigator.userAgent);
    // Check if Do Not Track is enabled
    const dnt = navigator.doNotTrack;
    const isToTrack = dnt === "1" || dnt === "yes" ? false : true;

    if (!consent && !isBot && isToTrack) {
      setIsBannerVisible(true);
    }
  }, []);

  const handleAccept = () => {
    setConsent(true);
    setIsBannerVisible(false);
    // Add logic for enabling tracking or analytics
  };

  const handleReject = () => {
    setConsent(false);
    setIsBannerVisible(false);
    // Delete tracking cookies
    trackingCookiesNames.forEach(deleteCookie);
    // Disable analytics
  };

  const setConsent = (consent: boolean) => {
    const value = consent ? COOKIE_VALUE_ACCEPT : COOKIE_VALUE_REJECT;
    localStorage.setItem(COOKIE_NAME, value);
    setCookie(COOKIE_NAME, value, COOKIE_EXPIRATION_TIME);
  };

  const getCookie = (name: string) => {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop()?.split(";").shift();
  };

  const setCookie = (name: string, value: string, expirationMs: number) => {
    const date = new Date();
    date.setTime(date.getTime() + expirationMs);
    document.cookie = `${name}=${value}; expires=${date.toUTCString()}; path=/; SameSite=Lax; Secure`;
  };

  const deleteCookie = (name: string) => {
    document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; SameSite=Lax; Secure`;
  };

  if (!isBannerVisible) return null; // Hide banner if user has already consented or rejected

  return (
    <div className="fixed bottom-0 left-0 right-0 p-4 bg-gray-900 text-white z-50">
      <Card className="max-w-xl mx-auto">
        <CardHeader>
          <h2 className="text-lg font-semibold">Cookies & Privacy</h2>
        </CardHeader>
        <CardContent>
          <p>We use cookies to enhance your experience. By clicking "Accept", you consent to the use of cookies for analytics and personalized content.</p>
        </CardContent>
        <CardFooter className="flex justify-end space-x-4">
          <Button variant="outline" onClick={handleReject}>Reject</Button>
          <Button onClick={handleAccept}>Accept</Button>
        </CardFooter>
      </Card>
    </div>
  );
}
