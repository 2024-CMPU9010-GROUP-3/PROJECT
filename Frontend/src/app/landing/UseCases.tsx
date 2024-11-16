import { Check } from "lucide-react";

export const UseCases = () => (
  <div className="py-12 bg-gray-50">
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <h2 className="text-3xl font-extrabold text-center text-gray-900 mb-12">
        Magpie Works for Everyone
      </h2>
      <div className="grid grid-cols-1 gap-8 lg:grid-cols-2">
        <div className="bg-white p-8 rounded-lg shadow-md">
          <h3 className="text-xl font-bold mb-4">
            Urban Planners & City Authorities
          </h3>
          <ul className="space-y-4">
            {[
              "Optimize infrastructure planning",
              "Analyze public service distribution",
              "Generate stakeholder-ready reports",
            ].map((item) => (
              <li key={item} className="flex items-center">
                <Check className="h-5 w-5 text-green-500 mr-2" />
                <span>{item}</span>
              </li>
            ))}
          </ul>
        </div>
        <div className="bg-white p-8 rounded-lg shadow-md">
          <h3 className="text-xl font-bold mb-4">
            Event Planners & Organizers
          </h3>
          <ul className="space-y-4">
            {[
              "Find the perfect venue",
              "Assess transport accessibility",
              "Plan logistics with confidence",
            ].map((item) => (
              <li key={item} className="flex items-center">
                <Check className="h-5 w-5 text-green-500 mr-2" />
                <span>{item}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  </div>
);
