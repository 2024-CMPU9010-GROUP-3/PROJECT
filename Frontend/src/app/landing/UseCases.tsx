"use client";

import { Check } from "lucide-react";
import { motion } from "framer-motion";

const container = {
  hidden: { opacity: 0 },
  show: {
    opacity: 1,
    transition: { staggerChildren: 0.2 },
  },
};

const item = {
  hidden: { opacity: 0, y: 20 },
  show: { opacity: 1, y: 0 },
};

export const UseCases = () => (
  <div
    className="py-24 bg-gradient-to-b from-gray-50 to-white relative overflow-hidden"
    id="use-cases"
  >
    <div className="absolute inset-0 bg-grid-pattern opacity-[0.02]" />

    <div className="relative max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
      <motion.h2
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        className="text-3xl sm:text-4xl font-extrabold text-center mb-16"
      >
        <span className="bg-clip-text text-transparent bg-gradient-to-r from-black to-gray-800">
          Magpie Works for Everyone
        </span>
      </motion.h2>

      <motion.div
        variants={container}
        initial="hidden"
        whileInView="show"
        viewport={{ once: true }}
        className="grid grid-cols-1 gap-8 lg:grid-cols-2"
      >
        {[
          {
            title: "Urban Planners & City Authorities",
            items: [
              "Optimize infrastructure planning",
              "Analyze public service distribution",
              "Generate stakeholder-ready reports",
            ],
          },
          {
            title: "General Users & Residents",
            items: [
              "Find nearby services and amenities",
              "Access transport options",
              "Plan efficient daily routes",
            ],
          },
        ].map((section) => (
          <motion.div
            key={section.title}
            variants={item}
            className="backdrop-blur-sm bg-white/50 p-8 rounded-2xl border border-gray-100 hover:border-gray-200 shadow-sm hover:shadow-md transition-all duration-300"
          >
            <h3 className="text-xl font-bold mb-6 text-gray-900">
              {section.title}
            </h3>
            <ul className="space-y-4">
              {section.items.map((text) => (
                <li key={text} className="flex items-center">
                  <div className="p-1 bg-green-50 rounded-full mr-3">
                    <Check className="h-4 w-4 text-green-500" />
                  </div>
                  <span className="text-gray-600">{text}</span>
                </li>
              ))}
            </ul>
          </motion.div>
        ))}
      </motion.div>
    </div>
  </div>
);
