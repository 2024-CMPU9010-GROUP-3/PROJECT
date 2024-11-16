import { Building, Calendar, MapPin } from "lucide-react";
import { FeatureCard } from "./FeatureCard";

export const Features = () => (
  <div className="py-12 bg-white">
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div className="text-center">
        <h2 className="text-3xl font-extrabold text-gray-900">
          Features That Make Magpie Unmatched
        </h2>
      </div>
      <div className="mt-10 grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3">
        <FeatureCard
          icon={MapPin}
          title="Interactive Map Visualization"
          description="Search, filter, and visualize essential services with ease."
        />
        <FeatureCard
          icon={Building}
          title="Smart Search & Analysis"
          description="Get detailed insights and multi-layer data visualizations."
        />
        <FeatureCard
          icon={Calendar}
          title="Professional Tools"
          description="Export reports and analyze historical trends for smarter decisions."
        />
      </div>
    </div>
  </div>
);
