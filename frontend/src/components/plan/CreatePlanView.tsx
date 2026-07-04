"use client";

import { useRouter } from "next/navigation";

import { useCreatePlanContainer } from "@/containers/CreatePlanContainer";

export function CreatePlanView() {
  const router = useRouter();
  const { title, loading, error, onChangeTitle, onSubmit } = useCreatePlanContainer();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const plan = await onSubmit();
    if (plan) router.push(`/plans/${plan.share_token}`);
  };

  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-6">
      <h1 className="text-2xl font-semibold">New Plan</h1>
      <form onSubmit={handleSubmit} className="flex w-full max-w-sm flex-col gap-4">
        <input
          type="text"
          value={title}
          onChange={(e) => onChangeTitle(e.target.value)}
          placeholder="Plan title"
          disabled={loading}
          className="rounded border px-3 py-2 outline-none focus:ring-2"
        />
        {error && <p className="text-sm text-red-600">{error}</p>}
        <button
          type="submit"
          disabled={loading || !title.trim()}
          className="rounded bg-black px-4 py-2 text-white disabled:opacity-50"
        >
          {loading ? "Creating…" : "Create Plan"}
        </button>
      </form>
    </div>
  );
}
