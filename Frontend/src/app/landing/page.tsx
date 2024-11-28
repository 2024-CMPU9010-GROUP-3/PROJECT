import { Navbar } from "../../components/global/Navbar";
import { CTA } from "./CTA";
import { Features } from "./Features";
import HeroSection from "./HeroSection";
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
