"use client";
import React, { useState } from "react";
import type { CardComponentProps } from "onborda";
import { useOnborda } from "onborda";
import confetti from "canvas-confetti";

// Shadcn
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../../../components/ui/card";

// Icons
import { X } from "lucide-react";

export const TourCard: React.FC<CardComponentProps> = ({
  step,
  currentStep: initialStep,
  totalSteps,
  nextStep,
  prevStep,
  arrow,
}) => {
  // Onborda hooks
  const { closeOnborda } = useOnborda();

  // State for current step
  const [currentStep, setCurrentStep] = useState(initialStep);

  function handleConfetti() {
    closeOnborda();
    
    // Run confetti once immediately
    confetti({
        particleCount: 100,
        startVelocity: 30,
        spread: 360,
        origin: {
            x: Math.random(),
            y: Math.random() - 0.2
        }
    });

    // Set up interval to run confetti 5 times, once per second
    let count = 0;
    const interval = setInterval(() => {
        confetti({
            particleCount: 100,
            startVelocity: 30,
            spread: 360,
            origin: {
                x: Math.random(),
                y: Math.random() - 0.2
            }
        });
        count++;
        if (count === 4) {
            clearInterval(interval);
        }
    }, 500);
}

  function handleNextStep() {
    setCurrentStep((prev) => Math.min(prev + 1, totalSteps - 1));
    nextStep();
  }

  function handlePrevStep() {
    setCurrentStep((prev) => Math.max(prev - 1, 0));
    prevStep();
  }

  return (
    <Card className="relative border-0 min-w-[300px] w-full max-w-[90%] md:max-w-[70%] lg:max-w-[50%] z-[999] bg-white border-none mx-auto">
      <CardHeader>
        <div className="flex items-start justify-between w-full space-x-4">
          <div className="flex flex-col space-y-2">
            <CardDescription className="text-black/50">
              {currentStep + 1} of {totalSteps}
            </CardDescription>
            <CardTitle className="mb-2 text-lg text-black">
              {step.icon} {step.title}
            </CardTitle>
          </div>
          <Button
            variant="ghost"
            className="text-black/50 absolute top-4 right-2 hover:bg-transparent hover:text-black/80"
            size="icon"
            onClick={() => closeOnborda()}
          >
            <X size={16} />
          </Button>
        </div>
      </CardHeader>
      <CardContent className="text-black">{step.content}</CardContent>
      <CardFooter className="text-black">
        <div className="flex justify-between w-full gap-4">
          {currentStep !== 0 && (
            <Button
              onClick={handlePrevStep}
              className="bg-zinc-900 hover:bg-zinc-800 text-white hover:text-white"
            >
              Previous
            </Button>
          )}
          {currentStep + 1 !== totalSteps && (
            <Button
              onClick={handleNextStep}
              className="bg-zinc-900 hover:bg-zinc-800 text-white hover:text-white ml-auto"
            >
              Next
            </Button>
          )}
          {currentStep + 1 === totalSteps && (
            <Button
              className="bg-zinc-900 hover:bg-zinc-800 text-white hover:text-white ml-auto"
              onClick={handleConfetti}
            >
              🎉 Finish!
            </Button>
          )}
        </div>
      </CardFooter>
      <span className="text-white">{arrow}</span>
    </Card>
  );
};