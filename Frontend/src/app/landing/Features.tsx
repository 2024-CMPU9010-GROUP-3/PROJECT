"use client";

import { Building, Calendar, MapPin } from "lucide-react";
import { motion } from "framer-motion";
import { FeatureCard } from "./FeatureCard";

const container = {
  hidden: { opacity: 0 },
  show: {
    opacity: 1,
    transition: {
      staggerChildren: 0.2,
    },
  },
};

export const Features = () => (
  <div className="py-24 bg-gradient-to-b from-white to-gray-50 relative overflow-hidden">
    <div className="absolute inset-0 bg-grid-pattern opacity-[0.02]" />

    <div className="relative max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.8 }}
        className="text-center"
      >
        <h2 className="text-3xl sm:text-4xl font-extrabold">
          <span className="bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-purple-600">
            Features That Make Magpie Unmatched
          </span>
        </h2>
      </motion.div>

      <motion.div
        variants={container}
        initial="hidden"
        whileInView="show"
        viewport={{ once: true }}
        className="mt-16 grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3"
      >
        <FeatureCard
          icon={MapPin}
          title="Interactive Map Visualization"
          description="Search, filter, and visualize essential services with ease."
          className="backdrop-blur-sm bg-white/50 border border-gray-100 hover:border-blue-100 shadow-sm hover:shadow-md transition-all duration-300"
        />
        <FeatureCard
          icon={Building}
          title="Smart Search & Analysis"
          description="Get detailed insights and multi-layer data visualizations."
          className="backdrop-blur-sm bg-white/50 border border-gray-100 hover:border-blue-100 shadow-sm hover:shadow-md transition-all duration-300"
        />
        <FeatureCard
          icon={Calendar}
          title="Professional Tools"
          description="Export reports and analyze historical trends for smarter decisions."
          className="backdrop-blur-sm bg-white/50 border border-gray-100 hover:border-blue-100 shadow-sm hover:shadow-md transition-all duration-300"
        />
      </motion.div>
    </div>
  </div>
);
