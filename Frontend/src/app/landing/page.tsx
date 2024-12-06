import { CTA } from "./CTA";
import { Features } from "./Features";
import HeroSection from "./HeroSection";
import { Navbar } from "./Navbar";
import { UseCases } from "./UseCases";
import WhatIsMagpie from "./WhatIsMagpie";
const page = () => {
  return (
    <div>
      <Navbar />
      <HeroSection />
      <WhatIsMagpie />
      <Features />
      <UseCases />
      <CTA />
    </div>
  );
};
export default page;
