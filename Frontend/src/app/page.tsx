import LoginPage from "./login/page";
import { HydrationBoundary, dehydrate } from "@tanstack/react-query";
import { getQueryClient } from "@/app/get-query-client";

const Page = () => {
  const queryClient = getQueryClient();

  return (
    <>
      <HydrationBoundary state={dehydrate(queryClient)}>
        <LoginPage />
      </HydrationBoundary>
    </>
  );
};

export default Page;
