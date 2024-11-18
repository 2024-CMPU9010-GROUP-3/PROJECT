"use client";

import { ChevronRight } from "lucide-react";
import { motion } from "framer-motion";

export const CTA = () => (
  <div className="bg-gradient-to-r from-blue-600 to-blue-700 relative overflow-hidden">
    <div className="absolute inset-0 bg-grid-pattern opacity-[0.02]" />

    <div className="relative max-w-6xl mx-auto py-16 px-4 sm:px-6 lg:py-20 lg:px-8 lg:flex lg:items-center lg:justify-between">
      <motion.div
        initial={{ opacity: 0, x: -20 }}
        whileInView={{ opacity: 1, x: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.6 }}
      >
        <h2 className="text-3xl font-extrabold tracking-tight sm:text-4xl">
          <span className="block text-white">
            Ready to Transform How You Plan?
          </span>
          <span className="block text-blue-200 mt-2 text-2xl sm:text-3xl">
            Join countless professionals making smarter decisions.
          </span>
        </h2>
      </motion.div>

      <motion.div
        initial={{ opacity: 0, x: 20 }}
        whileInView={{ opacity: 1, x: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.6 }}
        className="mt-8 flex lg:mt-0 lg:flex-shrink-0"
      >
        <motion.button
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          className="group px-8 py-3 bg-white rounded-xl font-medium text-lg flex items-center justify-center hover:shadow-lg transition-all duration-200 text-blue-600"
        >
          Sign Up Now
          <ChevronRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
        </motion.button>
      </motion.div>
    </div>

    <div className="absolute bottom-0 left-0 right-0 h-px bg-gradient-to-r from-transparent via-blue-300/20 to-transparent" />
  </div>
);
