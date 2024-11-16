import { FC } from "react";

interface FeatureCardProps {
  icon: FC<{ className?: string }>;
  title: string;
  description: string;
}

export const FeatureCard: FC<FeatureCardProps> = ({
  icon: Icon,
  title,
  description,
}) => (
  <div className="p-6 bg-white rounded-lg shadow-md">
    <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mb-4">
      <Icon className="h-6 w-6 text-blue-600" />
    </div>
    <h3 className="text-lg font-medium text-gray-900">{title}</h3>
    <p className="mt-2 text-gray-500">{description}</p>
  </div>
);
