"use client";

import { ArrowRight, MapPin } from "lucide-react";
import { motion } from "framer-motion";

const HeroSection = () => {
  return (
    <div className="relative min-h-[80vh] bg-gradient-to-br from-gray-50 via-white to-gray-50 flex items-center">
      <div className="absolute inset-0 bg-grid-pattern opacity-[0.02]" />

      <div className="relative max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          className="text-center space-y-8"
        >
          <h1 className="text-4xl sm:text-5xl md:text-6xl font-bold text-gray-900 tracking-tight">
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-black to-gray-800">
              Magpie
            </span>
            <span className="block mt-2 text-2xl sm:text-3xl md:text-4xl text-gray-600 font-medium">
              Services at a glance
            </span>
          </h1>

          <p className="max-w-xl mx-auto text-lg sm:text-xl text-gray-600 leading-relaxed">
            Access Dublin&apos;s most comprehensive database of public
            services—quickly, easily, and reliably.
          </p>

          <div className="flex flex-col sm:flex-row gap-4 justify-center mt-8">
            <motion.button
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              className="group px-8 py-3 bg-gradient-to-r from-black to-gray-800 text-white rounded-xl font-medium text-lg flex items-center justify-center hover:shadow-lg transition-all duration-200"
            >
              Explore the Map
              <MapPin className="ml-2 h-5 w-5 group-hover:rotate-12 transition-transform" />
            </motion.button>

            <motion.button
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              className="group px-8 py-3 bg-white text-black rounded-xl font-medium text-lg flex items-center justify-center border border-gray-100 hover:border-gray-200 hover:shadow-md transition-all duration-200"
            >
              Learn More
              <ArrowRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
            </motion.button>
          </div>
        </motion.div>
      </div>

      <div className="absolute bottom-0 left-0 right-0 h-32 bg-gradient-to-t from-white to-transparent" />
    </div>
  );
};

export default HeroSection;
