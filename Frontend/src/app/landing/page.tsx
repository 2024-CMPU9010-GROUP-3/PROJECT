import { CTA } from "./CTA";
import { Features } from "./Features";
import HeroSection from "./HeroSection";
import { Navbar } from "./Navbar";
import { UseCases } from "./UseCases";
const page = () => {
  return (
    <div>
      <Navbar />
      <HeroSection />
      <Features />
      <UseCases />
      <CTA />
    </div>
  );
};
export default page;
