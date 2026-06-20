import type { Plan } from "@/domain/plans/types";

export async function getPlanByShareToken(shareToken: string): Promise<Plan | null> {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/plans/${shareToken}`);

  if (res.status === 404) {
    return null;
  }

  if (!res.ok) {
    throw new Error(`Failed to fetch plan: ${res.status}`);
  }

  return res.json();
}
