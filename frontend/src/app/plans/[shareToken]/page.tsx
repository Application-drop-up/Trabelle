import { PlanView } from "@/components/plan/PlanView";

export default async function PlanPage({ params }: { params: Promise<{ shareToken: string }> }) {
  const { shareToken } = await params;
  return <PlanView shareToken={shareToken} />;
}
