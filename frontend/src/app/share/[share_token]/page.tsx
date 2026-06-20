import { notFound } from "next/navigation";

import { getPlanByShareToken } from "@/domain/plans/api";
import { PlanView } from "@/domain/plans/components/PlanView";

export default async function SharedPlanPage({
  params,
}: {
  params: Promise<{ share_token: string }>;
}) {
  const { share_token } = await params;
  const plan = await getPlanByShareToken(share_token);

  if (!plan) {
    notFound();
  }

  return <PlanView plan={plan} />;
}
