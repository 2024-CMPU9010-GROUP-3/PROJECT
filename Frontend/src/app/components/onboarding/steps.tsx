import { Tour } from "onborda/dist/types";

export const steps: Tour[] = [
    {
        tour: "firsttour",
        steps: [
            {
                icon: <>👋</>,
                title: "Tour 1, Step 1",
                content: <>First tour, first step</>,
                selector: "#tour1-step1",
                side: "top",
                showControls: true,
                pointerPadding: 10,
                pointerRadius: 10,
                nextRoute: "/foo",
                prevRoute: "/bar"
            },
        ],
    }
];