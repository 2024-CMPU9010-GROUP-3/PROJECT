// FeatureCard.tsx
import { LucideIcon } from "lucide-react";
import { motion } from "framer-motion";

interface FeatureCardProps {
  icon: LucideIcon;
  title: string;
  description: string;
  className?: string;
}

const item = {
  hidden: { opacity: 0, y: 20 },
  show: { opacity: 1, y: 0 },
};

export const FeatureCard = ({
  icon: Icon,
  title,
  description,
  className,
}: FeatureCardProps) => {
  return (
    <motion.div variants={item} className={`rounded-2xl p-6 ${className}`}>
      <div className="flex flex-col items-center text-center space-y-4">
        <div className="p-3 bg-gradient-to-br from-gray-50 to-gray-200 rounded-xl">
          <Icon className="w-6 h-6 text-gray-800" />
        </div>

        <h3 className="text-xl font-bold text-gray-900">{title}</h3>

        <p className="text-gray-600 leading-relaxed">{description}</p>
      </div>
    </motion.div>
  );
};
