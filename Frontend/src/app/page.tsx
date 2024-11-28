import LandingPage from "./landing/page";
import { HydrationBoundary, dehydrate } from "@tanstack/react-query";
import { getQueryClient } from "@/app/get-query-client";

const Page = () => {
  const queryClient = getQueryClient();

  return (
    <>
      <HydrationBoundary state={dehydrate(queryClient)}>
        <LandingPage />
      </HydrationBoundary>
    </>
  );
};

export default Page;
